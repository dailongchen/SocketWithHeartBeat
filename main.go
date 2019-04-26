package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	logInfo    *log.Logger
	logWarning *log.Logger
	logError   *log.Logger
)

func main() {
	errFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("failed to open log file", err)
	}

	logInfo = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	logWarning = log.New(os.Stdout, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	logError = log.New(io.MultiWriter(os.Stderr, errFile), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

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
