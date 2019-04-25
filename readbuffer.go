package main

import (
	"bytes"
	"io"
)

type readBuffer struct {
	buffer bytes.Buffer
	reader io.Reader
}

func createReadBuffer(reader io.Reader) *readBuffer {
	return &readBuffer{reader: reader}
}

func (rb *readBuffer) read(count int) ([]byte, error) {
	if count == 0 {
		return nil, nil
	}

	lr := io.LimitReader(rb.reader, int64(count))
	for rb.buffer.Len() < count {
		_, err := rb.buffer.ReadFrom(lr)
		if err != nil {
			return nil, err
		}
	}

	strBuf := make([]byte, count)
	if _, err := rb.buffer.Read(strBuf); err != nil {
		return nil, err
	}

	return strBuf, nil
}
