package main

import (
	"fmt"
	"log"
	"os"
)

type Main struct {
	isServer, isClient bool
	filename           string
	toURI              string
}

func main() {
	main := &Main{}

	if err := parseArgs(main); err != nil {
		log.Fatalf("Bad args: %v", err)
	}

	if main.isServer {
		// start server
	}

	if main.isClient {
		// start client
	}

	os.Exit(0)
}

func parseArgs(main *Main) error {
	args := os.Args

	for i, arg := range args {
		switch arg {
		case "-s":
			main.isServer = true
		case "-c":
			main.isClient = true
		case "-b":
			main.isClient = true
			main.isServer = true
		default:
			main.toURI = args[i]
			main.filename = args[i+1]
		}
	}

	if !main.isClient && !main.isServer {
		return fmt.Errorf("specify server, client, or both")
	}

	if len(main.toURI) == 0 {
		return fmt.Errorf("specify a destination URI")
	}

	if len(main.filename) == 0 {
		return fmt.Errorf("specify a filename to send/receive")
	}

	return nil
}
