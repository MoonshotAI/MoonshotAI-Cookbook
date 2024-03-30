package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testKey    = "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	testServer = newTestServer()
	testServer.StartTLS()
	code := m.Run()
	testServer.Close()
	os.Exit(code)
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.Handle("/models", &handler{
		KeyResponse: keyResponse{
			plainTextKey: {
				StatusCode:   500,
				ContentType:  "text/plain",
				ResponseBody: "plain_text",
			},
			jsonArrayKey: {
				StatusCode:   500,
				ContentType:  "application/json",
				ResponseBody: "[]",
			},
		},
		ResponseData: &Models{
			Data: []struct {
				ID         string            `json:"id"`
				Object     string            `json:"object"`
				OwnedBy    string            `json:"owned_by"`
				Permission []json.RawMessage `json:"permission"`
			}{
				{
					ID:         "moonshot-v1-8k",
					Object:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					OwnedBy:    "moonshot",
					Permission: nil,
				},
			},
			Object: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
	})
	mux.Handle("/chat/completions/", &handler{
		KeyResponse: keyResponse{
			streamStopKey: {
				StatusCode:  200,
				ContentType: "text/event-stream",
				ResponseBody: `data: {"id": "xxxx-xxxx-xxxx", "model": "moonshot-v1-8k", "object": "xxxx", "created": 1700000000, "choices": [{"index": 0, "delta": {"role": "assistant", "content": "h", "tool_calls": [{"index": 0, "id": "xxxx", "type": "function", "function": {"name": "test", "arguments": "te"}}], "finish_reason": ""}}]}` + "\n" +
					`data: {"id": "xxxx-xxxx-xxxx", "model": "moonshot-v1-8k", "object": "xxxx", "created": 1700000000, "choices": [{"index": 0, "delta": {"role": "assistant", "content": "i", "tool_calls": [{"index": 0, "id": "", "type": "", "function": {"name": "", "arguments": "st"}}], "finish_reason": "stop"}}]}` + "\n" +
					`data: [DONE]` + "\n",
			},
			streamDoneKey: {
				StatusCode:  200,
				ContentType: "text/event-stream",
				ResponseBody: `data: {"id": "xxxx-xxxx-xxxx", "model": "moonshot-v1-8k", "object": "xxxx", "created": 1700000000, "choices": [{"index": 0, "delta": {"role": "assistant", "content": "h", "tool_calls": [{"index": 0, "id": "xxxx", "type": "function", "function": {"name": "test", "arguments": "te"}}], "finish_reason": ""}}]}` + "\n" +
					`data: {"id": "xxxx-xxxx-xxxx", "model": "moonshot-v1-8k", "object": "xxxx", "created": 1700000000, "choices": [{"index": 0, "delta": {"role": "assistant", "content": "i", "tool_calls": [{"index": 0, "id": "", "type": "", "function": {"name": "", "arguments": "st"}}], "finish_reason": "length"}}]}` + "\n" +
					`data: [DONE]` + "\n",
			},
			invalidStreamTypeKey: {
				StatusCode:   200,
				ContentType:  "application/json",
				ResponseBody: "data: [DONE]",
			},
			invalidStreamContentKey: {
				StatusCode:   200,
				ContentType:  "text/event-stream",
				ResponseBody: "data: []",
			},
		},
		ResponseData: &Completion{
			ID:                "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			Model:             "moonshot-v1-8k",
			Object:            "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			Created:           1700000000,
			SystemFingerprint: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			Choices: []struct {
				Index        NullableType[int]    `json:"index"`
				Message      *Message             `json:"message"`
				FinishReason NullableType[string] `json:"finish_reason"`
				LogProbs     json.RawMessage      `json:"logprobs"`
			}{
				{
					Index: "0",
					Message: &Message{
						Role:      RoleAssistant,
						Content:   &Content{Text: "hi"},
						ToolCalls: nil,
					},
					FinishReason: "stop",
					LogProbs:     nil,
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		},
	})
	mux.Handle("/files", http.HandlerFunc(uploadFile))
	return httptest.NewUnstartedServer(mux)
}

var (
	plainTextKey            = "plain_text_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	jsonArrayKey            = "json_array_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	streamStopKey           = "stream_stop_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	streamDoneKey           = "stream_done_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	invalidStreamTypeKey    = "invalid_stream_type_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	invalidStreamContentKey = "invalid_stream_content_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)

type (
	returnType struct {
		StatusCode   int
		ContentType  string
		ResponseBody string
	}
	keyResponse map[string]*returnType
)

type handler struct {
	KeyResponse  keyResponse
	ResponseData any
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var errorResponse struct {
		Error *Error `json:"error"`
	}
	if rt, ok := checkAuthorization(r.Header, h.KeyResponse); !ok {
		errorResponse.Error = &Error{
			Message: "Unauthorized",
			Type:    "unauthorized_error",
			Param:   "unauthorized_error",
			Code:    "unauthorized_error",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&errorResponse)
		return
	} else if rt != nil {
		w.Header().Set("Content-Type", rt.ContentType)
		w.WriteHeader(rt.StatusCode)
		io.WriteString(w, rt.ResponseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.ResponseData)
	return
}

func checkAuthorization(header http.Header, keyResponse keyResponse) (*returnType, bool) {
	authorization := header.Get("Authorization")
	authorization = strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer"))
	if authorization == testKey {
		return nil, true
	}
	for k, r := range keyResponse {
		if authorization == k {
			return r, true
		}
	}
	return nil, false
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, ok := checkAuthorization(r.Header, keyResponse{}); !ok {
		var unauthorizedError = struct {
			Error *Error `json:"error"`
		}{
			Error: &Error{
				Message: "Unauthorized",
				Type:    "unauthorized_error",
				Param:   "unauthorized_error",
				Code:    "unauthorized_error",
			},
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&unauthorizedError)
		return
	}
	var internalServerError struct {
		Error *Error `json:"error"`
	}
	internalServerError.Error = &Error{
		Message: "Internal Server Error",
		Type:    "internal_server_error",
		Param:   "internal_server_error",
		Code:    "internal_server_error",
	}
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&internalServerError)
		return
	}
	purpose := r.FormValue("purpose")
	if purpose == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&internalServerError)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&internalServerError)
		return
	}
	defer file.Close()
	if header.Filename != "client_test.go" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&internalServerError)
		return
	}
	hasher := md5.New()
	io.Copy(hasher, file)
	srcHashCode := hasher.Sum(nil)
	hasher.Reset()
	file, _ = os.Open("client_test.go")
	defer file.Close()
	io.Copy(hasher, file)
	tgtHashCode := hasher.Sum(nil)
	if !bytes.Equal(tgtHashCode, srcHashCode) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&internalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&File{
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Object:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Bytes:     1024,
		CreatedAt: 1700000000,
		Filename:  "client_test.go",
		Purpose:   purpose,
	})
}

type testCaller struct {
	baseUrl string
	key     string
	client  *http.Client
}

func (tc *testCaller) BaseUrl() string      { return tc.baseUrl }
func (tc *testCaller) Key() string          { return tc.key }
func (tc *testCaller) Client() *http.Client { return tc.client }

func newClient(key string) Client[*testCaller] {
	return NewClient[*testCaller](&testCaller{
		baseUrl: testServer.URL,
		key:     key,
		client:  testServer.Client(),
	})
}

func TestClient(t *testing.T) {
	t.Run("error", testClientError)
	t.Run("models", testClientModels)
	t.Run("chat", testClientChat)
	t.Run("file", testClientFile)
}

func testClientError(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		brokenClient := newClient(plainTextKey)
		_, err := brokenClient.ListModels(context.Background())
		if err == nil {
			t.Errorf("brokenClient.ListModels: expects errors, got nothing")
			return
		}
		defer CloseErrorResponseBody(err)
		parsedError := ParseError(err)
		if parsedError != nil {
			serialized, _ := json.Marshal(parsedError)
			t.Errorf("ParseError: not Error, but got => %s", string(serialized))
			return
		}
		response := err.(getResponse).Response()
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		if !bytes.Equal(body, []byte("plain_text")) {
			t.Errorf("http.Response.Body: %s != plain_text", string(body))
			return
		}
	})
	t.Run("2", func(t *testing.T) {
		brokenClient := newClient(jsonArrayKey)
		_, err := brokenClient.ListModels(context.Background())
		if err == nil {
			t.Errorf("brokenClient.ListModels: expects errors, got nothing")
			return
		}
		defer CloseErrorResponseBody(err)
		parsedError := ParseError(err)
		if parsedError != nil {
			serialized, _ := json.Marshal(parsedError)
			t.Errorf("ParseError: not Error, but got => %s", string(serialized))
			return
		}
		response := err.(getResponse).Response()
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		if !bytes.Equal(body, []byte("[]")) {
			t.Errorf("http.Response.Body: %s != []", string(body))
			return
		}
	})
}

func testClientModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newClient(testKey)
		models, err := client.ListModels(context.Background())
		if err != nil {
			t.Errorf("client.ListModels: %s", err)
			return
		}
		if len(models.Data) != 1 || models.Data[0].ID != "moonshot-v1-8k" {
			serialized, _ := json.Marshal(models)
			t.Errorf("client.ListModels: unexpected Models object => %s", string(serialized))
			return
		}
	})
	t.Run("fail", func(t *testing.T) {
		noKeyClient := newClient("")
		_, err := noKeyClient.ListModels(context.Background())
		if err == nil {
			t.Errorf("noKeyClient.ListModels: expects errors, got nothing")
			return
		}
		defer CloseErrorResponseBody(err)
		parsedError := ParseError(err)
		if parsedError == nil {
			t.Errorf("ParseError: expects Error, got => %s", err)
			return
		}
		if parsedError.Type != "unauthorized_error" {
			t.Errorf("parsedError.Type != unauthorized_error")
			return
		}
		if parsedError.Error() != "moonshot: Unauthorized" {
			t.Errorf("parsedError.Error() != Unauthorized")
			return
		}
	})
}

func testClientChat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newClient(testKey)
		completion, err := client.CreateChatCompletion(context.Background(), &ChatCompletionRequest{
			Messages: []*Message{
				{
					Role:    RoleUser,
					Content: &Content{Text: "hi"},
				},
			},
			Model: "moonshot-v1-8k",
		})
		if err != nil {
			t.Errorf("client.CreateChatCompletion: %s", err)
			return
		}
		if completion.Model != "moonshot-v1-8k" || len(completion.Choices) != 1 ||
			completion.GetFinishReason() != FinishReasonStop || completion.GetMessageRole() != RoleAssistant ||
			completion.GetMessageContent() != "hi" || completion.GetPromptTokens() != 10 ||
			completion.GetCompletionTokens() != 5 || completion.GetTotalTokens() != 15 ||
			completion.GetPromptTokens()+completion.GetCompletionTokens() != completion.GetTotalTokens() {
			serialized, _ := json.Marshal(completion)
			t.Errorf("client.CreateChatCompletion: unexpected Completion object => %s", string(serialized))
			return
		}
		t.Run("stream", func(t *testing.T) {
			t.Run("stop", func(t *testing.T) {
				streamClient := newClient(streamStopKey)
				stream, err := streamClient.CreateChatCompletionStream(context.Background(), &ChatCompletionStreamRequest{
					Messages: []*Message{
						{
							Role:    RoleUser,
							Content: &Content{Text: "hi"},
						},
					},
					Model:          "moonshot-v1-8k",
					ToolChoice:     "test",
					ResponseFormat: ResponseFormatText,
				})
				if err != nil {
					t.Errorf("streamClient.CreateChatCompletionStream: %s", err)
					return
				}
				defer stream.Close()
				message := stream.CollectMessage()
				if streamErr := stream.Err(); streamErr != nil {
					t.Errorf("stream.Err(): %s", streamErr)
					return
				}
				if message.Role != RoleAssistant || message.Content.Text != "hi" ||
					message.ToolCalls[0].Type != ToolTypeFunction || message.ToolCalls[0].Function.Name != "test" ||
					message.ToolCalls[0].Function.Arguments != "test" {
					serialized, _ := json.Marshal(message)
					t.Errorf("streamClient.CreateChatCompletionStream: unexpected Stream Message => %s", string(serialized))
					return
				}
			})
			t.Run("done", func(t *testing.T) {
				streamClient := newClient(streamDoneKey)
				stream, err := streamClient.CreateChatCompletionStream(context.Background(), &ChatCompletionStreamRequest{
					Messages: []*Message{
						{
							Role:    RoleUser,
							Content: &Content{Text: "hi"},
						},
					},
					Model:          "moonshot-v1-8k",
					ToolChoice:     "test",
					ResponseFormat: ResponseFormatText,
				})
				if err != nil {
					t.Errorf("streamClient.CreateChatCompletionStream: %s", err)
					return
				}
				defer stream.Close()
				message := stream.CollectMessage()
				if streamErr := stream.Err(); streamErr != nil {
					t.Errorf("stream.Err(): %s", streamErr)
					return
				}
				if message.Role != RoleAssistant || message.Content.Text != "hi" ||
					message.ToolCalls[0].Type != ToolTypeFunction || message.ToolCalls[0].Function.Name != "test" ||
					message.ToolCalls[0].Function.Arguments != "test" {
					serialized, _ := json.Marshal(message)
					t.Errorf("streamClient.CreateChatCompletionStream: unexpected Stream Message => %s", string(serialized))
					return
				}
			})
			t.Run("error", func(t *testing.T) {
				t.Run("type", func(t *testing.T) {
					invalidStreamClient := newClient(invalidStreamTypeKey)
					_, streamErr := invalidStreamClient.CreateChatCompletionStream(context.Background(), &ChatCompletionStreamRequest{
						Messages: []*Message{
							{
								Role:    RoleUser,
								Content: &Content{Text: "hi"},
							},
						},
						Model: "moonshot-v1-8k",
					})
					if streamErr == nil {
						t.Errorf("invalidStreamClient.CreateChatCompletionStream: expects errors, got nothing")
						return
					}
					if !errors.Is(streamErr, ErrNotEventStream) {
						t.Errorf("streamErr != ErrNotEventStream => %s", streamErr)
						return
					}
				})
				t.Run("content", func(t *testing.T) {
					invalidStreamClient := newClient(invalidStreamContentKey)
					stream, streamErr := invalidStreamClient.CreateChatCompletionStream(context.Background(), &ChatCompletionStreamRequest{
						Messages: []*Message{
							{
								Role:    RoleUser,
								Content: &Content{Text: "hi"},
							},
						},
						Model: "moonshot-v1-8k",
					})
					if streamErr != nil {
						t.Errorf("invalidStreamClient.CreateChatCompletionStream: %s", streamErr)
						return
					}
					defer stream.Close()
					if streamErr = stream.Err(); streamErr != nil {
						t.Errorf("stream.Err(): %s", streamErr)
						return
					}
					stream.CollectMessage()
					if streamErr = stream.Err(); streamErr == nil {
						t.Errorf("stream.Err() expects errors, got nothing")
						return
					}
				})
			})
		})
	})
}

func testClientFile(t *testing.T) {
	client := newClient(testKey)
	testFile, err := os.Open("client_test.go")
	if err != nil {
		t.Errorf("os.Open: %s", err)
		return
	}
	defer testFile.Close()
	file, err := client.UploadFile(context.Background(), &UploadFileRequest{
		File:     testFile,
		Filename: "client_test.go",
		Purpose:  "fine-tune",
	})
	if err != nil {
		t.Errorf("client.UploadFile: %s", err)
		return
	}
	if file.Filename != "client_test.go" || file.Purpose != "fine-tune" {
		serialized, _ := json.Marshal(file)
		t.Errorf("client.UploadFile: unexpected File object => %s", string(serialized))
		return
	}
}

func TestGenericNew(t *testing.T) {
	if __ClientNew[*ResponseHandler]() == nil {
		t.Errorf("__ClientNew[*ResponseHandler]() got nil value")
		return
	}
}
