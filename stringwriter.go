package main

type stringWriter interface {
	writeString(string) error
}
