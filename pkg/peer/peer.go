package peer

import (
	"fmt"
	"net"
	"os"

	"github.com/AndriyAntonenko/my-peer/pkg/utils"
)

type IncomingMessageHandler = func(id string, msg string)
type Peer struct {
	host string
	port uint16
}

func NewPeer(host string, port uint16) *Peer {
	peer := Peer{host: host, port: port}
	return &peer
}

func NewPeerFromSocket(conn *net.TCPConn) *Peer {
	address := utils.SplitAddress(conn.RemoteAddr().String())
	peer := Peer{host: address.Host, port: address.Port}
	return &peer
}

func (p *Peer) Write(msg string) error {
	address := fmt.Sprintf("%s:%d", p.host, p.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot resolve peer address %s", err.Error())
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	return err
}
