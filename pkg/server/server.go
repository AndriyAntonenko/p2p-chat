package server

import (
	"fmt"
	"net"

	"github.com/AndriyAntonenko/my-peer/pkg/utils"
)

func RunServer(address string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	utils.HandleFatalError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	utils.HandleFatalError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		_, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		fmt.Println(string(buf[0:]))
	}
}
