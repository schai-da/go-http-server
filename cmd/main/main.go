package main

import (
	"go-http-server/internal/server"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	keepAliveHandler := server.NewKeepAliveHandler()
	responseHandler := server.NewResponseHandler()
	requestHandler := server.NewRequestHandler()
	lastModifiedHandler := server.NewLastModifiedHandler()
	fileHandler := server.NewFileHandler()
	etagHandler := server.NewEtagHandler()

	keepAliveHandler.SetNext(responseHandler).SetNext(requestHandler).SetNext(lastModifiedHandler).SetNext(fileHandler).SetNext(etagHandler)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		// ゴルーチンによる並列処理
		go func(conn net.Conn) {
			defer conn.Close()

			handleError := keepAliveHandler.Handle(conn, server.NewRequest(), server.NewResponse())
			if handleError != nil {
				log.Println("Failed to handle request:", handleError)
				return
			}
		}(conn)
	}
}
