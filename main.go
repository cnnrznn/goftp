package main

import (
	"fmt"
	"os"
)

type Main struct {
	isServer, isClient bool
}

func main() {
	main := &Main{}

	if err := parseArgs(main); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func parseArgs(main *Main) error {
	args := os.Args

	for _, arg := range args {
		switch arg {
		case "-s":
			main.isServer = true
		case "-c":
			main.isClient = true
		case "-b":
			main.isClient = true
			main.isServer = true
		}
	}

	if !main.isClient && !main.isServer {
		return fmt.Errorf("Bad args, specify server, or client, or both")
	}

	return nil
}
