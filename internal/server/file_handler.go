package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	PublicDir = "public"
)

type fileHandler struct {
	next Handler
}

func NewFileHandler() *fileHandler {
	return &fileHandler{}
}

func (h *fileHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *fileHandler) Handle(conn net.Conn, request *Request, response *Response) error {

	if request.Method == "GET" {

		filePath := fmt.Sprintf("%s%s", PublicDir, request.Path)
		if strings.HasSuffix(filePath, "/") {
			filePath += "index.html"
		}

		file, err := os.Open(filePath)
		if err != nil {
			response.StatusCode = 404
			response.Body = []byte("<h1>404 Not Found</h1>")
			return nil
		}
		defer file.Close()

		body, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		response.StatusCode = 200
		response.Headers["Content-Type"] = "text/html"
		response.Body = body
	} else {
		response.StatusCode = 405
		return nil
	}
	if h.next != nil {
		return h.next.Handle(conn, request, response)
	}
	return nil
}
