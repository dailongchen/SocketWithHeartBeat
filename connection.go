package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

const heartBeatMessage string = "heart beat"

type connection struct {
	connect         net.Conn
	mg              messager
	closedChan      chan bool
	readChan        chan string
	writeChan       chan string
	errChan         chan error
	endReadRoution  chan bool
	endWriteRoutine chan bool
	err             error
}

func connected(ctx context.Context, connect net.Conn, timeout int64, heartBeat bool) *connection {
	reader := newNetReader(connect, timeout)
	writer := newNetWriter(connect, timeout)
	mg := messager{reader, writer}

	closedChan := make(chan bool)

	readChan := make(chan string)
	writeChan := make(chan string)

	errChan := make(chan error)

	endReadRoution := make(chan bool)
	endWriteRoutine := make(chan bool)

	c := &connection{
		connect,
		mg,
		closedChan,
		readChan,
		writeChan,
		errChan,
		endReadRoution,
		endWriteRoutine,
		nil}

	c.run(ctx, heartBeat)

	go func() {
		select {
		case <-ctx.Done():
			break
		case err := <-errChan:
			c.err = err
			if err == io.EOF {
				fmt.Println("closed by peer...")
			} else {
				fmt.Printf("Error %s\n", err)
			}
			break
		}

		close(readChan)
		close(writeChan)
		c.connect.Close()

		endedRoutineCount := 0
		for {
			select {
			case <-endReadRoution:
				endedRoutineCount++
			case <-endWriteRoutine:
				endedRoutineCount++
			case <-errChan:
				break
			}

			if endedRoutineCount >= 2 {
				break
			}
		}

		close(endReadRoution)
		close(endWriteRoutine)
		close(errChan)

		fmt.Println("Connection closed")
		close(closedChan)
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
				fmt.Println("reading")
				message, err := c.mg.read()
				if err != nil {
					cancel()
					c.errChan <- err
					break EndFor
				}
				fmt.Printf("read: %s\n", message)

				c.readChan <- message
			}
		}

		fmt.Println("End of read routine")
		c.endReadRoution <- true
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
				fmt.Printf("writing: %s\n", sendingMessage)
				if err := c.mg.write(sendingMessage); err != nil {
					cancel()
					c.errChan <- err
					break EndFor
				}
				fmt.Println("write done")
			}
		}

		fmt.Println("End of write routine")
		c.endWriteRoutine <- true
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
