package server

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
)

type etagHandler struct {
	NextHandler Handler
}

func NewEtagHandler() *etagHandler {
	return &etagHandler{}
}

func (h *etagHandler) SetNext(handler Handler) Handler {
	h.NextHandler = handler
	return handler
}

func (h *etagHandler) Handle(conn net.Conn, request *Request, response *Response) error {
	hash := sha256.Sum256(response.Body)
	etag := `"` + hex.EncodeToString(hash[:]) + `"`

	ifNoMatch, ifNoMatchOk := request.Headers["If-None-Match"]

	if ifNoMatchOk && ifNoMatch == etag {
		// Etagが一致する場合、304 Not Modifiedを返す
		response.StatusCode = 304
		response.Body = nil
		return nil
	}

	response.Headers["Etag"] = etag
	response.Headers["Cache-Control"] = "public, max-age=0"
	if h.NextHandler != nil {
		return h.NextHandler.Handle(conn, request, response)
	}

	return nil
}
