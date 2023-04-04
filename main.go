package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cnnrznn/goftp/client"
	"github.com/cnnrznn/goftp/server"
)

type Main struct {
	isServer, isClient bool
	filename           string
	toURI              string
}

func main() {
	main := &Main{}
	stop := make(chan error)

	if err := parseArgs(main); err != nil {
		fmt.Printf("Bad args: %v\n", err)
		return
	}

	if main.isServer {
		// start server
		go server.New(main.filename).Run(stop)
	}

	if main.isClient {
		// start client
		go client.New(main.filename, main.toURI).Run(stop)
	}

	for err := range stop {
		if err != nil {
			log.Fatalln(err)
		}
	}

	os.Exit(0)
}

// Replace this piece of trash with something better.
func parseArgs(main *Main) error {
	args := os.Args[1:]

	for i, arg := range args {
		switch arg {
		case "-s":
			main.isServer = true
		case "-c":
			main.isClient = true
		default:
			main.toURI = args[i]
			main.filename = args[i+1]
			break
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
