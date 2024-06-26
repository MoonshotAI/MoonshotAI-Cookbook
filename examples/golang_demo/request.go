package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"runtime"
	"sync"
)

type ChatCompletionRequest struct {
	Messages         []*Message            `json:"messages"`
	Model            string                `json:"model"`
	LogProbs         bool                  `json:"logprobs,omitempty"`
	MaxTokens        int                   `json:"max_tokens,omitempty"`
	N                int                   `json:"n,omitempty"`
	ResponseFormat   ResponseFormat        `json:"response_format,omitempty"`
	Seed             int                   `json:"seed,omitempty"`
	Temperature      NullableType[float64] `json:"temperature,omitempty"`
	TopP             NullableType[float64] `json:"top_p,omitempty"`
	PresencePenalty  float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64               `json:"frequency_penalty,omitempty"`
	Tools            []*Tool               `json:"tools,omitempty"`
	ToolChoice       ToolChoice            `json:"tool_choice,omitempty"`
}

func (r *ChatCompletionRequest) ToJSON() (string, error) {
	return toJSON(r)
}

type ChatCompletionStreamRequest ChatCompletionRequest

func (r *ChatCompletionStreamRequest) ToJSON() (string, error) {
	type ChatCompletionStreamRequestWrapper struct {
		*ChatCompletionStreamRequest
		Stream bool `json:"stream"`
	}
	return toJSON(&ChatCompletionStreamRequestWrapper{
		ChatCompletionStreamRequest: r,
		Stream:                      true,
	})
}

type EstimateTokenCountRequest ChatCompletionRequest

func (r *EstimateTokenCountRequest) ToJSON() (string, error) {
	type EstimateTokenCountRequestWrapper struct {
		Model    string     `json:"model"`
		Messages []*Message `json:"messages"`
	}
	return toJSON(&EstimateTokenCountRequestWrapper{
		Model:    r.Model,
		Messages: r.Messages,
	})
}

type CreateContextCacheRequest struct {
	Messages    []*Message        `json:"messages"`
	Model       string            `json:"model"`
	Tools       []*Tool           `json:"tools"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ExpiredAt   int               `json:"expiredAt,omitempty"`
	TTL         int               `json:"ttl,omitempty"`
}

func (r *CreateContextCacheRequest) ToJSON() (string, error) {
	return toJSON(r)
}

type ContextCacheOptions struct {
	CacheID  string
	ResetTTL int
}

type contextKeyCacheID struct{}

func withCacheOptions(ctx context.Context, options *ContextCacheOptions) context.Context {
	return context.WithValue(ctx, contextKeyCacheID{}, options)
}

func getCacheOptions(ctx context.Context) *ContextCacheOptions {
	cv := ctx.Value(contextKeyCacheID{})
	if cv == nil {
		return nil
	}
	options, ok := cv.(*ContextCacheOptions)
	if !ok {
		return nil
	}
	return options
}

func toJSON(obj any) (string, error) {
	bytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type UploadFileRequest struct {
	File     io.Reader
	Filename string
	Purpose  string

	once        sync.Once
	reader      *io.PipeReader
	errorChan   <-chan error
	contentType string
}

func (r *UploadFileRequest) init() {
	r.once.Do(func() {
		const (
			formFieldPurpose = "purpose"
			formFieldFile    = "file"
		)

		var (
			err                    error
			pipeReader, pipeWriter = io.Pipe()
			writer                 = multipart.NewWriter(pipeWriter)
			errorChan              = make(chan error, 1)
		)

		r.reader = pipeReader
		r.contentType = writer.FormDataContentType()
		r.errorChan = errorChan

		// To prevent goroutine leaks, utilize the SetFinalizer function to tidy up the Producer-side goroutine
		// upon the completion of the UploadFileRequest lifecycle.
		runtime.SetFinalizer(r, func(*UploadFileRequest) { pipeReader.Close() })

		go func() {
			defer func(ch chan<- error) {
				if err != nil {
					if !errors.Is(err, io.ErrClosedPipe) {
						ch <- err
					}
				}
				close(ch)
				pipeWriter.Close()
			}(errorChan)
			if err = writer.WriteField(formFieldPurpose, r.Purpose); err != nil {
				return
			}
			var filename string
			if naming, ok := r.File.(interface{ Name() string }); ok { // os.File has Name() method
				filename = naming.Name()
			}
			if r.Filename != "" {
				filename = r.Filename
			}
			var part io.Writer
			part, err = writer.CreateFormFile(formFieldFile, filename)
			if err != nil {
				return
			}
			if _, err = io.Copy(part, r.File); err != nil {
				return
			}
			if err = writer.Close(); err != nil {
				return
			}
			if err = pipeWriter.Close(); err != nil {
				return
			}
		}()
	})
}

func (r *UploadFileRequest) ContentType() string {
	r.init()
	return r.contentType
}

func (r *UploadFileRequest) Read(p []byte) (n int, err error) {
	r.init()
	for {
		var (
			alreadyEOF bool
			chanClosed bool
		)
		{
			var (
				chanActive bool
			)
			select {
			case err, chanActive = <-r.errorChan:
				chanClosed = !chanActive
				if chanActive {
					if err != nil {
						return n, err
					}
				}
			default:
			}
		}
		if alreadyEOF && chanClosed {
			return n, io.EOF
		}
		n, err = r.reader.Read(p)
		if err != nil && errors.Is(err, io.EOF) && !chanClosed {
			alreadyEOF = true
			continue
		}
		return n, err
	}
}
