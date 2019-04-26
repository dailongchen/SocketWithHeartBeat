package main

import (
	"context"
	"fmt"
	"net"
)

func connEchoHandler(c net.Conn) {
	if c == nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connection := connected(ctx, c, 20, true)

	for {
		content, ok, err := connection.readMessage()
		if !ok {
			break
		}
		if err != nil {
			logError.Printf("Failed to read, %s\n", err)
			break
		}

		if err := connection.writeMessage(content); err != nil {
			logError.Printf("Failed to write, %s\n", err)
			break
		}
	}
}

func runServer(port int) {
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logError.Printf("Fail to start server, %s\n", err)
		return
	}

	fmt.Println("Server Started ...")

	for {
		conn, err := server.Accept()
		if err != nil {
			logError.Printf("Fail to connect, %s\n", err)
			break
		}

		logInfo.Println("Accepted new connection")
		go connEchoHandler(conn)
	}
}
