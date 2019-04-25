package main

type stringReader interface {
	readString(int) (string, error)
}
