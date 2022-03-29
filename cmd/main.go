package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndriyAntonenko/my-peer/pkg/peer"
	"github.com/AndriyAntonenko/my-peer/pkg/topology"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s name host1:port1 [host2:port2], [host3:port3]...", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]
	address := os.Args[2]
	peers := os.Args[3:]

	network := topology.NewFullyConnectedTopology(name, address)

	network.OnConnection(func(id string, peer peer.Peer) {
		fmt.Fprintf(os.Stdout, "%s connected\n", id)
		peer.Write(fmt.Sprintf("Hello from %s!", network.GetMe()))
	})

	network.OnMessage(func(id string, msg string) {
		fmt.Fprintf(os.Stdout, "[%s]-> %s\n", id, msg)
	})

	network.Listen(address)
	if peers != nil && len(peers) > 0 {
		network.AddPeers(peers)
		network.Broadcast(fmt.Sprintf("Hello from %s!", network.GetMe()))
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// block main goroutine
	<-quit
}
