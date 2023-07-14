package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cnnrznn/goargs"

	"github.com/cnnrznn/goftp/client"
	"github.com/cnnrznn/goftp/server"
)

type Main struct {
	isServer, isClient bool
	filename           string
	toURI              string
}

func main() {
	addr := "localhost:9751"
	main := &Main{}
	stop := make(chan error)

	if err := parseArgs(main); err != nil {
		fmt.Printf("Bad args: %v\n", err)
		return
	}

	if main.isServer {
		// start server
		go server.New(addr, main.filename).Run(stop)
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
	parser := goargs.New()
	parser.Required("-type", goargs.One)
	parser.Required("-to", goargs.One)
	parser.Required("-fn", goargs.One)

	vals, err := parser.Parse()
	if err != nil {
		return err
	}

	if vals["-type"][0] == "server" {
		main.isServer = true
	}

	if len(main.toURI) == 0 {
		return fmt.Errorf("specify a destination URI")
	}

	if len(main.filename) == 0 {
		return fmt.Errorf("specify a filename to send/receive")
	}

	return nil
}
