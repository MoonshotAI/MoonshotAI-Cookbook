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

//go:generate go run -mod=mod "github.com/x5iu/defc" --mode=api --output=client.gen.go --features=api/nort,api/logx,api/error,api/future,api/client
type Client[C Caller] interface {
	// ListModels GET {{ $.Client.BaseUrl }}/models
	// Authorization: Bearer {{ $.Client.Key }}
	ListModels(ctx context.Context) (*Models, error)

	// CreateChatCompletion POST {{ $.Client.BaseUrl }}/chat/completions
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	//
	// {{ $.request.ToJSON }}
	CreateChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*Completion, error)

	// CreateChatCompletionStream POST {{ $.Client.BaseUrl }}/chat/completions
	// Content-Type: application/json
	// Authorization: Bearer {{ $.Client.Key }}
	//
	// {{ $.request.ToJSON }}
	CreateChatCompletionStream(ctx context.Context, request *ChatCompletionStreamRequest) (*Stream, error)

	// UploadFile POST {{ $.Client.BaseUrl }}/files
	// Content-Type: {{ $.request.ContentType }}
	// Authorization: Bearer {{ $.Client.Key }}
	UploadFile(ctx context.Context, request *UploadFileRequest) (*File, error)

	// RetrieveFileContent GET {{ $.Client.BaseUrl }}/files/{{ $.fileID }}/content
	// Authorization: Bearer {{ $.Client.Key }}
	RetrieveFileContent(ctx context.Context, fileID string) ([]byte, error)

	Inner() C
	response() *ResponseHandler
}
