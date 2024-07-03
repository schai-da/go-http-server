package server

import (
	"net"
	"time"
)

const (
	timeout = 5 * time.Second
)

type keepAliveHandler struct {
	next Handler
}

func NewKeepAliveHandler() *keepAliveHandler {
	return &keepAliveHandler{}
}

func (h *keepAliveHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *keepAliveHandler) Handle(conn net.Conn, request *Request, response *Response) error {
	defer conn.Close()

	if h.next != nil {
		for {
			err := h.next.Handle(conn, request, response)

			if err != nil {
				return err // エラーにはタイムアウトを含む
			}

			// Keep-Aliveヘッダーの確認
			keepAlive := request.Headers["Connection"] == "keep-alive"

			if !keepAlive {
				return nil
			}
			// keep-aliveの場合はループ処理を継続
			conn.SetReadDeadline(time.Now().Add(timeout))
			request = NewRequest()
			response = NewResponse()

			continue
		}
	} else {
		return nil
	}
}
