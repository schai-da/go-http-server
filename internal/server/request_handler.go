package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type requestHandler struct {
	next Handler
}

func NewRequestHandler() *requestHandler {
	return &requestHandler{}
}

func (h *requestHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *requestHandler) Handle(conn net.Conn, request *Request, response *Response) error {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) < 2 {
		return fmt.Errorf("invalid request line")
	}

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // 空行が見つかったらヘッダーの読み取りを終了
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) != 2 {
			continue // key-value形式でない
		}
		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])
		headers[key] = value
	}

	request.Method = parts[0]
	request.Path = parts[1]
	request.Headers = headers

	log.Printf("method: %s, path: %s\n", request.Method, request.Path)

	if h.next != nil {
		return h.next.Handle(conn, request, response)
	}

	return nil
}
