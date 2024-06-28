//go:build context_cache

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	knowledgeKey = "moonshot_knowledge_0002"
)

const (
	cacheStatusPending  = "pending"
	cacheStatusReady    = "ready"
	cacheStatusExpired  = "expired"
	cacheStatusError    = "error"
	cacheStatusInactive = "inactive"
)

var documents = []string{
	"moonshot.pdf",
	"moonshot-context-cache.pdf",
}

func initKnowledge(
	ctx context.Context,
	client Client[moonshot],
	manager ContextCacheManager,
) (messages []*Message, cacheID string, err error) {
	createTableOnce.Do(func() {
		err = manager.createTable(ctx)
	})
	if err != nil {
		return nil, "", err
	}
	var cache *ContextCache
	cache, err = manager.Get(ctx, knowledgeKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, "", err
	}
	if cache.ID != "" {
		cache, err = client.RetrieveContextCache(ctx, cache.ID)
		if err != nil {
			return nil, "", err
		}
		if err = manager.Set(ctx, knowledgeKey, cache); err != nil {
			return nil, "", err
		}
		if cache.Status == cacheStatusReady {
			return cache.Messages, cache.ID, nil
		}
	}
	if cache.ID == "" || cache.Status == cacheStatusError {
		for _, document := range documents {
			file, err := os.Open(document)
			if err != nil {
				return nil, "", err
			}
			uploadedFile, err := client.UploadFile(ctx, &UploadFileRequest{
				File:    file,
				Purpose: "file-extract",
			})
			file.Close()
			if err != nil {
				return nil, "", err
			}
			fileContent, err := client.RetrieveFileContent(ctx, uploadedFile.ID)
			if err != nil {
				return nil, "", err
			}
			messages = append(messages, &Message{
				Role:    RoleSystem,
				Content: &Content{Text: string(fileContent)},
			})
		}
		cache, err = client.CreateContextCache(ctx, &CreateContextCacheRequest{
			Messages: messages,
			Model:    "moonshot-v1",
			TTL:      3600,
		})
		if err != nil {
			return nil, "", err
		}
		if err = manager.Set(ctx, knowledgeKey, cache); err != nil {
			return nil, "", err
		}
	} else {
		messages = cache.Messages
	}
	return messages, "", nil
}

var createTableOnce sync.Once

//go:generate go run -mod=mod github.com/x5iu/defc generate --output=context_cache_manager.gen.go --features=sqlx/nort
type ContextCacheManager interface {
	// createTable Exec
	/*
		create table if not exists context_cache
		(
		    id    integer not null
					constraint context_cache_pk
						primary key autoincrement,
		    key            text    not null,
		    cache_id       text    not null,
			cache_status   text    not null,
			cache_messages text    not null
		);
	*/
	createTable(ctx context.Context) error

	// Set Exec arguments=args
	/*
		delete from context_cache where key = {{ $.args.Add $.key }};
		insert into context_cache (
			key,
			cache_id,
			cache_status,
			cache_messages
		) values (
			{{ $.args.Add $.key }},
			{{ $.args.Add $.cache.ID }},
			{{ $.args.Add $.cache.Status }},
			{{ $.args.Add $.cache.Messages }}
		);
	*/
	Set(ctx context.Context, key string, cache *ContextCache) error

	// Get Query One Const
	// select cache_id, cache_status, cache_messages from context_cache where key = ?;
	Get(ctx context.Context, key string) (*ContextCache, error)
}

type moonshot struct {
	baseUrl string
	key     string
	client  *http.Client
	log     func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration)
}

func (m moonshot) BaseUrl() string      { return m.baseUrl }
func (m moonshot) Key() string          { return m.key }
func (m moonshot) Client() *http.Client { return m.client }

func (m moonshot) Log(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
	m.log(ctx, caller, request, response, elapse)
}

func main() {
	ctx := context.Background()

	client := NewClient[moonshot](moonshot{
		baseUrl: "https://api.moonshot.cn/v1",
		key:     os.Getenv("MOONSHOT_API_KEY"),
		client:  http.DefaultClient,
		log: func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
			var usingCache string
			if response != nil {
				cacheID := response.Header.Get("Msh-Context-Cache-Id")
				if cacheID != "" {
					tokenTTL := response.Header.Get("Msh-Context-Cache-Token-Ttl")
					tokenSaved := response.Header.Get("Msh-Context-Cache-Token-Saved")
					usingCache = fmt.Sprintf("; using cache: %s, ttl=%s, tokens=%s",
						cacheID,
						tokenTTL,
						tokenSaved,
					)
				}
			}
			log.Printf("[%s] %s %s%s", caller, request.URL, elapse, usingCache)
		},
	})

	cacheManager := NewContextCacheManager("sqlite3", os.Getenv("MOONSHOT_CACHE_DB"))

	messages, cacheID, err := initKnowledge(ctx, client, cacheManager)
	if err != nil {
		if parsed := ParseError(err); parsed != nil {
			log.Fatalln("("+parsed.Type+")", parsed.Message)
		}
		log.Fatalln(err)
	}

	chatWithCacheByContext(ctx, client, messages, cacheID)
	chatWithCacheByMessage(ctx, client, messages, cacheID)

	if err = RemoveAllFiles(ctx, client); err != nil {
		if parsed := ParseError(err); parsed != nil {
			log.Fatalln("("+parsed.Type+")", parsed.Message)
		}
		log.Fatalln(err)
	}
}

const input = "请列出 chat 所有支持的参数；并告诉我如何使用 Cache。"

func chatWithCacheByContext(
	ctx context.Context,
	client Client[moonshot],
	messages []*Message,
	cacheID string,
) {
	if cacheID != "" {
		ctx = withCacheOptions(ctx, &ContextCacheOptions{
			CacheID:  cacheID,
			ResetTTL: 3600,
		})
	}

	messages = append(messages, &Message{
		Role:    RoleUser,
		Content: &Content{Text: input},
	})

	stream, err := client.CreateChatCompletionStream(ctx, &ChatCompletionStreamRequest{
		Messages: messages,
		Model:    ModelMoonshot128K,
	})

	if err != nil {
		if parsed := ParseError(err); parsed != nil {
			log.Fatalln("("+parsed.Type+")", parsed.Message)
		}
		log.Fatalln(err)
	}

	defer stream.Close()
	for chunk := range stream.C {
		fmt.Printf("%s", chunk.GetDeltaContent())
	}
}

func chatWithCacheByMessage(
	ctx context.Context,
	client Client[moonshot],
	messages []*Message,
	cacheID string,
) {
	if cacheID != "" {
		messages = []*Message{
			{Role: RoleCache, Content: &Content{Cache: &ContextCacheOptions{CacheID: cacheID, ResetTTL: 3600}}},
			{Role: RoleUser, Content: &Content{Text: input}},
		}
	} else {
		messages = append(messages, &Message{
			Role:    RoleUser,
			Content: &Content{Text: input},
		})
	}

	stream, err := client.CreateChatCompletionStream(ctx, &ChatCompletionStreamRequest{
		Messages: messages,
		Model:    ModelMoonshot128K,
	})

	if err != nil {
		if parsed := ParseError(err); parsed != nil {
			log.Fatalln("("+parsed.Type+")", parsed.Message)
		}
		log.Fatalln(err)
	}

	defer stream.Close()
	for chunk := range stream.C {
		fmt.Printf("%s", chunk.GetDeltaContent())
	}
}

func RemoveAllFiles(ctx context.Context, client Client[moonshot]) error {
	files, err := client.ListFiles(ctx)
	if err != nil {
		return err
	}
	for _, file := range files.Data {
		if err = client.DeleteFile(ctx, file.ID); err != nil {
			return err
		}
	}
	return nil
}
