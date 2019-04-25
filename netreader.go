package main

import (
	"net"
	"time"
)

type netReaderWithTimeout struct {
	connect net.Conn
	timeout int64
}

func (reader *netReaderWithTimeout) Read(buf []byte) (int, error) {
	if reader.timeout > 0 {
		reader.connect.SetReadDeadline(time.Now().Add(time.Duration(reader.timeout) * time.Second))
	}
	return reader.connect.Read(buf)
}

type netReader struct {
	rb *readBuffer
}

func newNetReader(connect net.Conn, timeout int64) *netReader {
	readerWithTimeout := &netReaderWithTimeout{
		connect,
		timeout}
	readBuffer := createReadBuffer(readerWithTimeout)
	return &netReader{readBuffer}
}

func (nr *netReader) readString(count int) (string, error) {
	buf, err := nr.rb.read(count)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
