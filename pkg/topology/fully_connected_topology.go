package topology

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/AndriyAntonenko/my-peer/pkg/peer"
	"github.com/AndriyAntonenko/my-peer/pkg/utils"
)

type ConnectionHandler = func(id string, peer peer.Peer)
type PeerMessageHandler = func(id string, msg string)
type FullyConnectedTopology struct {
	me            string
	serverAddress string
	peers         map[string]*peer.Peer
	server        *net.TCPListener

	connectionHandlers  []ConnectionHandler
	peerMessageHandlers []PeerMessageHandler
}

func NewFullyConnectedTopology(me string, serverAddress string) *FullyConnectedTopology {
	topology := FullyConnectedTopology{
		me:            me,
		serverAddress: serverAddress,
		peers:         make(map[string]*peer.Peer),
	}

	return &topology
}

func (t *FullyConnectedTopology) GetMe() string {
	return t.me
}

func (t *FullyConnectedTopology) Listen(address string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	utils.HandleFatalError("cannot resolve server address", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	utils.HandleFatalError("cannot listen server address", err)

	t.server = listener
	go t.ListenIncomingPeers()
}

func (t *FullyConnectedTopology) ListenIncomingPeers() {
	for {
		conn, err := t.server.AcceptTCP()
		if err != nil {
			continue
		}

		newPeer := peer.NewPeerFromSocket(conn)
		peerId := t.initPeer(newPeer)

		if t.connectionHandlers != nil {
			for _, handler := range t.connectionHandlers {
				handler(peerId, *newPeer)
			}
		}
	}
}

func (t *FullyConnectedTopology) AddPeers(peers []string) {
	for _, peer := range peers {
		t.addPeer(peer)
	}
}

func (t *FullyConnectedTopology) OnMessage(handler PeerMessageHandler) {
	if t.peerMessageHandlers == nil {
		t.peerMessageHandlers = make([]PeerMessageHandler, 0)
	}
	t.peerMessageHandlers = append(t.peerMessageHandlers, handler)

	for _, peer := range t.peers {
		peer.OnMessage(handler)
	}
}

func (t *FullyConnectedTopology) OnConnection(handler ConnectionHandler) {
	if t.connectionHandlers == nil {
		t.connectionHandlers = make([]func(id string, peer peer.Peer), 0)
	}
	t.connectionHandlers = append(t.connectionHandlers, handler)
}

func (t *FullyConnectedTopology) Broadcast(msg string) {
	for _, p := range t.peers {
		p.Write(msg)
	}
}

func (t *FullyConnectedTopology) initPeer(p *peer.Peer) string {
	p.Connect()
	peerId := p.GetRemoteAddress()
	t.peers[peerId] = p

	if t.peerMessageHandlers != nil && len(t.peerMessageHandlers) > 0 {
		for _, handler := range t.peerMessageHandlers {
			p.OnMessage(handler)
		}
	}

	return peerId
}

func (t *FullyConnectedTopology) addPeer(address string) {
	parts := strings.Split(address, ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse address %s", address)
		return
	}
	host := parts[0]

	newPeer := peer.NewPeer(host, uint16(port))
	t.initPeer(newPeer)
}
