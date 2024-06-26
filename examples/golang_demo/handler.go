package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ResponseHandler struct {
	Caller   string
	Response *http.Response
}

func (r *ResponseHandler) FromResponse(caller string, response *http.Response) (err error) {
	r.Caller = caller
	r.Response = response
	return nil
}

var ErrNotEventStream = errors.New("moonshot: response Content-Type is not " + httpContentTypeSSE + " while Stream=true")

func (r *ResponseHandler) ScanValues(values ...any) (err error) {
	if len(values) == 0 {
		return nil
	}
	val := values[0]
	switch obj := val.(type) {
	case *Stream:
		if r.Caller == CallerCreateChatCompletionStream {
			if !isContentType(r.Response.Header, httpContentTypeSSE) {
				r.Response.Body.Close()
				return ErrNotEventStream
			}
			ch := make(chan *Chunk)
			ec := make(chan error, 1)
			obj.C = ch
			obj.error = ec
			go readStream(r.Response.Body, ch, ec)
			return nil
		}
	case *[]byte:
		defer r.Response.Body.Close()
		switch r.Caller {
		case CallerRetrieveFileContent:
			*obj, err = io.ReadAll(r.Response.Body)
			return err
		}
	}
	defer r.Response.Body.Close()
	decoder := json.NewDecoder(r.Response.Body)
	return decoder.Decode(val)
}

var (
	serverSentEventData  = []byte("data")
	serverSentEventDone  = []byte("[DONE]")
	serverSentEventColon = []byte{':'}
)

func readStream(r io.ReadCloser, ch chan<- *Chunk, ec chan<- error) {
	defer r.Close()
	defer close(ec)
	defer close(ch)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		data := scanner.Bytes()
		before, after, found := bytes.Cut(data, serverSentEventColon)
		if found {
			before, after = bytes.TrimSpace(before), bytes.TrimSpace(after)
			if bytes.Equal(before, serverSentEventData) {
				if bytes.Equal(after, serverSentEventDone) {
					break
				}
				var chunk Chunk
				if err := json.Unmarshal(after, &chunk); err != nil {
					ec <- err
					return
				}
				ch <- &chunk
				if chunk.GetFinishReason() == FinishReasonStop {
					break
				}
			}
		}
	}
	if scannerErr := scanner.Err(); scannerErr != nil {
		ec <- scannerErr
	}
}

func (r *ResponseHandler) Err() error  { return nil }
func (r *ResponseHandler) Break() bool { return true }

type getResponse interface {
	Response() *http.Response
}

func CloseErrorResponseBody(err error) {
	if gr, ok := err.(getResponse); ok {
		response := gr.Response()
		response.Body.Close()
	}
}

const (
	httpContentTypeJSON = "application/json"
	httpContentTypeSSE  = "text/event-stream"
)

var (
	errorPrefix = []byte(`{"error"`)
)

func ParseError(err error) *Error {
	if gr, ok := err.(getResponse); ok {
		response := gr.Response()
		if isContentType(response.Header, httpContentTypeJSON) {
			body, rdErr := io.ReadAll(response.Body)
			if rdErr != nil {
				return nil
			}
			response.Body.Close()
			first8 := make([]byte, 0, 8)
			for _, ch := range body {
				if !(ch == ' ' || ch == '\r' || ch == '\n' || ch == '\t') {
					first8 = append(first8, ch)
					if len(first8) >= 8 {
						break
					}
				}
			}
			if bytes.Equal(first8, errorPrefix) {
				var errorObject struct {
					Error *Error `json:"error"`
				}
				if decErr := json.Unmarshal(body, &errorObject); decErr == nil && errorObject.Error != nil {
					return errorObject.Error
				}
			}
			response.Body = io.NopCloser(bytes.NewReader(body))
		}
	}
	return nil
}

func isContentType(header http.Header, targetContentType string) bool {
	for _, headerContentType := range header.Values("Content-Type") {
		for i, ch := range headerContentType {
			if ch == ' ' || ch == ';' {
				headerContentType = headerContentType[:i]
				break
			}
		}
		if headerContentType == targetContentType {
			return true
		}
	}
	return false
}
