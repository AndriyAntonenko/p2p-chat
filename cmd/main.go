package main

import (
	"fmt"
	"os"

	"github.com/AndriyAntonenko/my-peer/pkg/server"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	server.RunServer(os.Args[1])
}
