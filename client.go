package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

func connClientHandler(c net.Conn) {
	if c == nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	connection := connected(ctx, c, 20, false)

	for i := 0; i < 15; i++ {
		if i > 5 && i < 10 {
			time.Sleep(5 * time.Second)
		}

		if err := connection.writeMessage(strconv.Itoa(i)); err != nil {
			logError.Printf("Failed to write, %s\n", err)
			break
		}

		_, ok, err := connection.readMessage()
		if !ok {
			break
		}
		if err != nil {
			logError.Printf("Failed to read, %s\n", err)
			break
		}
	}

	cancel()

	<-connection.closedChan
	logInfo.Println("Done")
}

func runClient(serverAddress string, port int) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", serverAddress, port), 2*time.Second)
	if err != nil {
		logError.Printf("Fail to connect, %s\n", err)
		return
	}

	connClientHandler(conn)
}
