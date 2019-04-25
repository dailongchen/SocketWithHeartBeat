package main

import (
	"net"
	"time"
)

type netWriterWithTimeout struct {
	connect net.Conn
	timeout int64
}

func (writer *netWriterWithTimeout) Write(content []byte) (int, error) {
	if len(content) == 0 {
		return 0, nil
	}

	if writer.timeout > 0 {
		writer.connect.SetWriteDeadline(time.Now().Add(time.Duration(writer.timeout) * time.Second))
	}
	return writer.connect.Write(content)
}

type netWriter struct {
	writer netWriterWithTimeout
}

func newNetWriter(connect net.Conn, timeout int64) *netWriter {
	writerWithTimeout := netWriterWithTimeout{
		connect,
		timeout}

	return &netWriter{writerWithTimeout}
}

func (nw *netWriter) writeString(content string) error {
	_, err := nw.writer.Write([]byte(content))

	return err
}
