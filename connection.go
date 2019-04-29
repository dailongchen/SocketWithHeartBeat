package main

import (
	"context"
	"io"
	"net"
	"time"
)

const heartBeatMessage string = "heart beat"

type connection struct {
	connect    net.Conn
	mg         messager
	closedChan chan bool
	readChan   chan string
	writeChan  chan string
	errChan    chan error
	stopRead   chan bool
	stopWrite  chan bool
	err        error
}

func connected(ctx context.Context, connect net.Conn, timeout int64, heartBeat bool) *connection {
	reader := newNetReader(connect, timeout)
	writer := newNetWriter(connect, timeout)
	mg := messager{reader, writer}

	closedChan := make(chan bool)

	readChan := make(chan string)
	writeChan := make(chan string)

	errChan := make(chan error)

	stopRead := make(chan bool)
	stopWrite := make(chan bool)

	c := &connection{
		connect,
		mg,
		closedChan,
		readChan,
		writeChan,
		errChan,
		stopRead,
		stopWrite,
		nil}

	c.run(ctx, heartBeat)

	go func() {
		select {
		case <-ctx.Done():
			break
		case err := <-errChan:
			c.err = err
			if err == io.EOF {
				logInfo.Println("closed by peer...")
			} else {
				logError.Printf("Error %s\n", err)
			}
			break
		}

		close(readChan)
		close(writeChan)
		c.connect.Close()

		endedRoutineCount := 0
		for {
			select {
			case <-stopRead:
				logInfo.Println("Stopped read routine")
				endedRoutineCount++
			case <-stopWrite:
				logInfo.Println("Stopped write routine")
				endedRoutineCount++
			case <-errChan:
				break
			}

			if endedRoutineCount >= 2 {
				break
			}
		}

		close(stopRead)
		close(stopWrite)
		close(errChan)
		close(closedChan)

		logInfo.Println("Connection closed")
	}()

	return c
}

func (c *connection) run(parentCtx context.Context, heartBeat bool) {
	ctx, cancel := context.WithCancel(parentCtx)

	// read routine
	go func() {
	EndFor:
		for {
			select {
			case <-ctx.Done():
				break EndFor
			default:
				message, err := c.mg.read()
				if err != nil {
					cancel()
					c.errChan <- err
					break EndFor
				}
				logInfo.Printf("read: %s\n", message)

				c.readChan <- message
			}
		}

		c.stopRead <- true
	}()

	// write routine
	go func() {
	EndFor:
		for {
			sendingMessage := ""

			select {
			case <-ctx.Done():
				break EndFor
			case message, ok := <-c.writeChan:
				if ok {
					sendingMessage = message
				}
			case <-time.After(2 * time.Second):
				if heartBeat {
					sendingMessage = heartBeatMessage
				}
			}

			if len(sendingMessage) > 0 {
				logInfo.Printf("writing: %s\n", sendingMessage)
				if err := c.mg.write(sendingMessage); err != nil {
					cancel()
					c.errChan <- err
					break EndFor
				}
				logInfo.Println("write done")
			}
		}

		c.stopWrite <- true
	}()
}

func (c *connection) readMessage() (string, bool, error) {
	if c.err != nil {
		return "", false, c.err
	}

	for {
		message, ok := <-c.readChan
		if !ok {
			return "", false, c.err
		}
		if heartBeatMessage != message {
			return message, ok, c.err
		}
	}
}

func (c *connection) writeMessage(content string) error {
	if c.err != nil {
		return c.err
	}

	c.writeChan <- content
	return c.err
}
