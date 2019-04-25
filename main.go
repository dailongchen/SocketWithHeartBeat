package main

import (
	"fmt"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s\n", err)
		}
	}()

	argsWithoutProg := os.Args[1:]

	runAsServer := len(argsWithoutProg) == 0

	port := 2234
	if !runAsServer {
		serverAddress := argsWithoutProg[0]
		runClient(serverAddress, port)
	} else {
		runServer(port)
	}
}
