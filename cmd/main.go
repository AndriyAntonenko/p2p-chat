package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndriyAntonenko/my-peer/pkg/topology"
	"github.com/AndriyAntonenko/my-peer/pkg/utils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s name host1:port1 [host2:port2], [host3:port3]...", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]
	address := os.Args[2]
	peers := os.Args[3:]

	go runNetwork(name, address, peers)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// block main goroutine
	<-quit
}

// @TODO: Create Text based user interface with list of users, messages box, and input field
// Lib to use https://github.com/gdamore/tcell
func readInput(network *topology.Topology) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		network.Broadcast(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		utils.HandleFatalError("scanner error", err)
	}
}

func runNetwork(name string, address string, peers []string) {
	network := topology.NewTopology(name, address)

	network.OnMessage(func(msg topology.Message) {
		log.Printf("[%s]-> %s\n", msg.AuthorName, msg.Content)
	})

	network.Listen(address)
	if peers != nil && len(peers) > 0 {
		network.AddPeers(peers)
		network.Broadcast(fmt.Sprintf("Hello from %s!", network.GetMe()))
	}

	go readInput(network)
}
