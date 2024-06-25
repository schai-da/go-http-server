package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	bufferSize = 1024
)

type HttpServer struct {
	Address string
}

func (s *HttpServer) Start() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("http server is listening on %s ...\n", s.Address)
	defer fmt.Println("... http server stopped")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, bufferSize)
	_, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Failed to read request:", err)
		return
	}
	// リクエスト内容をログに出力
	log.Printf("Received: %s", buffer)

	// リクエスト内容をそのままクライアントに返す
	_, err = conn.Write(buffer)
	if err != nil {
		log.Println("Failed to write response:", err)
		return
	}
}
