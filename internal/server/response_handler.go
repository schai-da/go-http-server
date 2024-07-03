package server

import (
	"fmt"
	"net"
)

const (
	chunkSize = 512
)

type responseHandler struct {
	next Handler
}

func NewResponseHandler() *responseHandler {
	return &responseHandler{}
}

func (h *responseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *responseHandler) Handle(conn net.Conn, request *Request, response *Response) error {
	if h.next == nil {
		return nil
	}

	// 処理結果をレスポンスに含めるため、先に次の処理を実行
	result := h.next.Handle(conn, request, response)
	if result != nil {
		response.StatusCode = 500
		response.Body = []byte("<h1>500 Internal Server Error</h1>")
	}

	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", response.StatusCode, statusText(response.StatusCode))
	for key, value := range response.Headers {
		fmt.Fprintf(conn, "%s: %s\r\n", key, value)
	}
	if response.Body != nil && len(response.Body) > chunkSize {
		return sendChunkedResponse(conn, *response)
	} else {
		return sendResponse(conn, *response)
	}
}

func sendResponse(conn net.Conn, response Response) error {
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(response.Body))
	fmt.Fprint(conn, "\r\n")
	_, err := conn.Write(response.Body)

	return err
}

func sendChunkedResponse(conn net.Conn, response Response) error {
	fmt.Fprint(conn, "Transfer-Encoding: chunked\r\n")
	for i := 0; i < len(response.Body); i += chunkSize {
		end := i + chunkSize
		if end > len(response.Body) {
			end = len(response.Body)
		}
		chunk := response.Body[i:end]

		// チャンクのサイズ
		if _, err := fmt.Fprintf(conn, "%x\r\n", len(chunk)); err != nil {
			return err
		}

		// チャンクの内容
		if _, err := conn.Write(chunk); err != nil {
			return err
		}

		// チャンクの終端
		if _, err := fmt.Fprintf(conn, "\r\n"); err != nil {
			return err
		}
	}

	// 全てのチャンクを送信したことを示す
	if _, err := fmt.Fprintf(conn, "0\r\n\r\n"); err != nil {
		return err
	}

	return nil
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 304:
		return "Not Modified"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status Code"
	}
}
