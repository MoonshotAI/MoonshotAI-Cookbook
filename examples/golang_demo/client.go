package main

import (
	"context"
	"net/http"
	"time"
)

type Caller interface {
	BaseUrl() string
	Key() string
}

type Logger interface {
	Log(ctx context.Context, caller string, request *http.Request, response *http.Response, elapse time.Duration)
}

type CustomHTTPClient interface {
	Client() *http.Client
}

//go:generate go run -mod=mod github.com/x5iu/defc generate --features=api/nort,api/logx,api/error,api/future,api/client --func get_cache_options=getCacheOptions
type Client[C Caller] interface {
	// ListModels GET {{ $.Client.BaseUrl }}/models
	// Authorization: Bearer {{ $.Client.Key }}
	ListModels(ctx context.Context) (*Models, error)

	// EstimateTokenCount POST {{ $.Client.BaseUrl }}/tokenizers/estimate-token-count
	// Authorization: Bearer {{ $.Client.Key }}
	//
	// {{ $.request.ToJSON }}
	EstimateTokenCount(ctx context.Context, request *EstimateTokenCountRequest) (*EstimateTokenCount, error)

	// CheckBalance GET {{ $.Client.BaseUrl }}/users/me/balance
	// Authorization: Bearer {{ $.Client.Key }}
	CheckBalance(ctx context.Context) (*Balance, error)

	// CreateChatCompletion POST {{ $.Client.BaseUrl }}/chat/completions
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	// {{ $options := (get_cache_options $.ctx) -}}
	// {{ if $options }}X-Msh-Context-Cache: {{ $options.CacheID }}
	// {{ end }}{{ if $options }}{{ if $options.ResetTTL }}X-Msh-Context-Cache-Reset-TTL: {{ $options.ResetTTL }}
	// {{ end }}{{ end }}
	//
	// {{ $.request.ToJSON }}
	CreateChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*Completion, error)

	// CreateChatCompletionStream POST {{ $.Client.BaseUrl }}/chat/completions
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	// {{ $options := (get_cache_options $.ctx) -}}
	// {{ if $options }}X-Msh-Context-Cache: {{ $options.CacheID }}
	// {{ end }}{{ if $options }}{{ if $options.ResetTTL }}X-Msh-Context-Cache-Reset-TTL: {{ $options.ResetTTL }}
	// {{ end }}{{ end }}
	//
	// {{ $.request.ToJSON }}
	CreateChatCompletionStream(ctx context.Context, request *ChatCompletionStreamRequest) (*Stream, error)

	// CreateContextCache POST {{ $.Client.BaseUrl }}/caching
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	//
	// {{ $.request.ToJSON }}
	CreateContextCache(ctx context.Context, request *CreateContextCacheRequest) (*ContextCache, error)

	// RetrieveContextCache GET {{ $.Client.BaseUrl }}/caching/{{ $.cacheID }}
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	RetrieveContextCache(ctx context.Context, cacheID string) (*ContextCache, error)

	// DeleteContextCache DELETE {{ $.Client.BaseUrl }}/caching/{{ $.cacheID }}
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	DeleteContextCache(ctx context.Context, cacheID string) error

	// UploadFile POST {{ $.Client.BaseUrl }}/files
	// Content-Type: {{ $.request.ContentType }}
	// Authorization: Bearer {{ $.Client.Key }}
	UploadFile(ctx context.Context, request *UploadFileRequest) (*File, error)

	// ListFiles GET {{ $.Client.BaseUrl }}/files
	// Authorization: Bearer {{ $.Client.Key }}
	ListFiles(ctx context.Context) (*Files, error)

	// DeleteFile DELETE {{ $.Client.BaseUrl }}/files/{{ $.fileID }}
	// Authorization: Bearer {{ $.Client.Key }}
	DeleteFile(ctx context.Context, fileID string) error

	// RetrieveFileContent GET {{ $.Client.BaseUrl }}/files/{{ $.fileID }}/content
	// Authorization: Bearer {{ $.Client.Key }}
	RetrieveFileContent(ctx context.Context, fileID string) ([]byte, error)

	Inner() C
	response() *ResponseHandler
}
