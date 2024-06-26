package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type closedReader struct{}

func (*closedReader) Read([]byte) (int, error) { return 0, os.ErrClosed }

func TestUploadFileRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := &UploadFileRequest{
			File:    strings.NewReader("test"),
			Purpose: "fine-tune",
		}
		if _, err := io.ReadAll(r); err != nil {
			t.Errorf("UploadFileRequest: %s", err)
			return
		}
	})
	t.Run("error", func(t *testing.T) {
		r := &UploadFileRequest{
			File:    &closedReader{},
			Purpose: "fine-tune",
		}
		if _, err := io.ReadAll(r); err == nil {
			t.Errorf("UploadFileRequest: expects errors, got nil")
			return
		} else if !errors.Is(err, os.ErrClosed) {
			t.Errorf("UploadFileRequest: expects os.ErrClosed, got => %s", err)
			return
		}
	})
}
