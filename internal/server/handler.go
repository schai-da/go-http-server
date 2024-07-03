package server

import "net"

type Handler interface {
	SetNext(handler Handler) Handler
	Handle(conn net.Conn, request *Request, response *Response) error
}
