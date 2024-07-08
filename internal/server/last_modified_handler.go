package server

import (
	"fmt"
	"net"
	"os"
	"time"
)

type lastModifiedHandler struct {
	NextHandler Handler
}

func NewLastModifiedHandler() *lastModifiedHandler {
	return &lastModifiedHandler{}
}

func (h *lastModifiedHandler) SetNext(handler Handler) Handler {
	h.NextHandler = handler
	return handler
}

func (h *lastModifiedHandler) Handle(conn net.Conn, request *Request, response *Response) error {
	err := func() error {
		if request.Method != "GET" {
			return nil
		}
		filePath := fmt.Sprintf("%s%s", PublicDir, request.Path)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil // ファイルが存在しない場合
		}
		lastModified := fileInfo.ModTime()
		response.Headers["Last-Modified"] = lastModified.Format(time.RFC1123)

		ifModifiedSince := request.Headers["If-Modified-Since"]
		if ifModifiedSince == "" {
			return nil
		}

		ifModifiedSinceTime, err := time.Parse(time.RFC1123, ifModifiedSince)
		if err != nil {
			return err
		}

		if lastModified.Before(ifModifiedSinceTime) || lastModified.Equal(ifModifiedSinceTime) {
			response.StatusCode = 304
		}
		return nil
	}()

	if err != nil {
		return err
	}

	if h.NextHandler != nil {
		return h.NextHandler.Handle(conn, request, response)
	}

	return nil
}
