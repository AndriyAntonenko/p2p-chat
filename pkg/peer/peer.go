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
	conn *net.TCPConn

	incomingMessageHandlers []IncomingMessageHandler
}

func NewPeer(host string, port uint16) *Peer {
	peer := Peer{host: host, port: port}
	return &peer
}

func NewPeerFromSocket(conn *net.TCPConn) *Peer {
	address := utils.SplitAddress(conn.RemoteAddr().String())
	peer := Peer{host: address.Host, port: address.Port, conn: conn}
	return &peer
}

func (p *Peer) Connect() {
	if p.conn != nil {
		go p.handleMessage(p.conn)
		return
	}

	address := fmt.Sprintf("%s:%d", p.host, p.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot resolve peer address %s", err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	p.conn = conn

	go p.handleMessage(p.conn)
}

func (p *Peer) Write(msg string) error {
	_, err := p.conn.Write([]byte(msg))
	return err
}

func (p *Peer) GetRemoteAddress() string {
	return p.conn.RemoteAddr().String()
}

func (p *Peer) OnMessage(handler IncomingMessageHandler) {
	if p.incomingMessageHandlers == nil {
		p.incomingMessageHandlers = make([]func(id string, msg string), 0)
	}
	p.incomingMessageHandlers = append(p.incomingMessageHandlers, handler)
}

func (p *Peer) handleMessage(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		if p.incomingMessageHandlers != nil && n > 0 {
			for _, handler := range p.incomingMessageHandlers {
				handler(p.GetRemoteAddress(), string(buf[0:]))
			}
		}
	}
}
