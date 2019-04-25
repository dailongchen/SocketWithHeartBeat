package main

import (
	"bytes"
	"encoding/binary"
)

type messager struct {
	sr stringReader
	sw stringWriter
}

func (mg *messager) read() (string, error) {
	lenString, err := mg.sr.readString(4)
	if err != nil {
		return "", err
	}

	var messageLength int32
	reader := bytes.NewReader([]byte(lenString))
	if err := binary.Read(reader, binary.LittleEndian, &messageLength); err != nil {
		return "", err
	}

	return mg.sr.readString(int(messageLength))
}

func (mg *messager) write(content string) error {
	contentLength := len(content)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, int32(contentLength)); err != nil {
		return err
	}

	buf.WriteString(content)
	message := buf.String()
	return mg.sw.writeString(message)
}
