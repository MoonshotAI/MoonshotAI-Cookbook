//go:build !generate_models_file
// +build !generate_models_file

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	_ Logger           = moonshot{}
	_ CustomHTTPClient = moonshot{}
)

type moonshot struct {
	host   string
	key    string
	client *http.Client
	log    func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration)
}

func (m moonshot) BaseUrl() string      { return m.host }
func (m moonshot) Key() string          { return m.key }
func (m moonshot) Client() *http.Client { return m.client }

func (m moonshot) Log(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
	m.log(ctx, caller, request, response, elapse)
}

func runDemo() error {
	ctx := context.Background()

	client := NewClient[moonshot](moonshot{
		host:   "https://api.moonshot.cn/v1",
		key:    os.Getenv("MOONSHOT_API_KEY"),
		client: http.DefaultClient,
		log: func(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration) {
			log.Printf("[%s] %s %s", caller, request.URL, elapse)
		},
	})

	completion, err := client.CreateChatCompletion(ctx, &ChatCompletionRequest{
		Messages: []*Message{
			{
				Role:    RoleSystem,
				Content: &Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    RoleUser,
				Content: &Content{Text: "你好，我叫李雷，1+1等于多少？"},
			},
		},
		Model:       ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	fmt.Println(completion.GetMessageContent())

	stream, err := client.CreateChatCompletionStream(ctx, &ChatCompletionStreamRequest{
		Messages: []*Message{
			{
				Role:    RoleSystem,
				Content: &Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    RoleUser,
				Content: &Content{Text: "写一个小故事，讲的是一个叫“龙猫”的勇士积极抵抗魔族入侵，保卫 Kimi 女神。"},
			},
		},
		Model:       ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	defer stream.Close()
	for chunk := range stream.C {
		fmt.Printf("%s", chunk.GetDeltaContent())
	}
	fmt.Println("")

	if err = stream.Err(); err != nil {
		return err
	}

	stream, err = client.CreateChatCompletionStream(ctx, &ChatCompletionStreamRequest{
		Messages: []*Message{
			{
				Role:    RoleSystem,
				Content: &Content{Text: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
			},
			{
				Role:    RoleUser,
				Content: &Content{Text: "写一个小故事，讲的是有一个叫“龙猫”的人，每天会在各个群聊里游荡，挑选一些感兴趣的话题回复，每个群都以得到龙猫老师的回复为荣，请写一个跌宕起伏的剧情，讲述“龙猫”与各个群聊的爱恨情仇。"},
			},
		},
		Model:       ModelMoonshot8K,
		MaxTokens:   4096,
		N:           1,
		Temperature: "0.3",
	})

	if err != nil {
		return err
	}

	defer stream.Close()
	message := stream.CollectMessage()
	fmt.Println(message.Content.Text)

	if err = stream.Err(); err != nil {
		return err
	}

	pdf, err := os.Open("moonshot.pdf")
	if err != nil {
		return err
	}

	defer pdf.Close()

	file, err := client.UploadFile(ctx, &UploadFileRequest{
		File:    pdf,
		Purpose: "file-extract",
	})

	if err != nil {
		return err
	}

	log.Printf("file_id=%q; status=%s", file.ID, file.Status)

	content, err := client.RetrieveFileContent(ctx, file.ID)
	if err != nil {
		return err
	}

	fmt.Println(string(content))

	return nil
}

func main() {
	if err := runDemo(); err != nil {
		log.Fatalln(err)
	}
}
