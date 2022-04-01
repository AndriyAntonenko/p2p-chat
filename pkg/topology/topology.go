package topology

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/AndriyAntonenko/my-peer/pkg/peer"
	"github.com/AndriyAntonenko/my-peer/pkg/utils"
)

type PeerMessageHandler = func(msg Message)
type Message struct {
	Kind          string `json:"kind"`
	AuthorName    string `json:"authorName"`
	AuthorAddress string `json:"authorAddress"`
	Content       string `json:"content"`
}

type Topology struct {
	me            string
	serverAddress string
	peers         map[string]*peer.Peer
	server        *net.TCPListener

	peerMessageHandlers []PeerMessageHandler
}

func NewTopology(me string, serverAddress string) *Topology {
	topology := Topology{
		me:            me,
		serverAddress: serverAddress,
		peers:         make(map[string]*peer.Peer),
	}

	return &topology
}

func (t *Topology) GetMe() string {
	return t.me
}

func (t *Topology) Listen(address string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	utils.HandleFatalError("cannot resolve server address", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	utils.HandleFatalError("cannot listen server address", err)

	t.server = listener
	go t.ListenIncomingPeers()
}

func (t *Topology) ListenIncomingPeers() {
	for {
		conn, err := t.server.AcceptTCP()
		if err != nil {
			continue
		}

		go t.handleMessage(conn)
	}
}

func (t *Topology) handleMessage(conn *net.TCPConn) {
	defer conn.Close()

	for {
		var msg Message
		decoder := json.NewDecoder(conn)
		if err := decoder.Decode(&msg); err != nil {
			return
		}

		if msg.Kind == "INTRO" {
			parts := utils.SplitAddress(msg.AuthorAddress)
			peerId := fmt.Sprintf("%s:%s", msg.AuthorName, msg.AuthorAddress)
			if _, alreadyInList := t.peers[peerId]; alreadyInList {
				return
			}
			newPeer := peer.NewPeer(parts.Host, parts.Port)
			t.peers[peerId] = newPeer
			t.sendIntroduceMessage(newPeer)
			return
		}

		if msg.Kind == "PLAIN" {
			peerId := fmt.Sprintf("%s:%s", msg.AuthorName, msg.AuthorAddress)
			if _, ok := t.peers[peerId]; !ok {
				return
			}
			for _, handler := range t.peerMessageHandlers {
				handler(msg)
			}
		}
	}
}

func (t *Topology) AddPeers(peers []string) {
	for _, peer := range peers {
		t.addPeer(peer)
	}
}

func (t *Topology) OnMessage(handler PeerMessageHandler) {
	if t.peerMessageHandlers == nil {
		t.peerMessageHandlers = make([]PeerMessageHandler, 0)
	}
	t.peerMessageHandlers = append(t.peerMessageHandlers, handler)

}

func (t *Topology) Broadcast(msg string) error {
	jsonMsg, err := json.Marshal(t.buildPlainMessage(msg))
	if err != nil {
		return err
	}

	for _, p := range t.peers {
		p.Write(string(jsonMsg))
	}
	return nil
}

func (t *Topology) addPeer(address string) {
	parts := strings.Split(address, ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse address %s", address)
		return
	}
	host := parts[0]

	newPeer := peer.NewPeer(host, uint16(port))
	t.sendIntroduceMessage(newPeer)
}

func (t *Topology) sendIntroduceMessage(p *peer.Peer) {
	msg, err := json.Marshal(t.buildIntroduceMessage())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	p.Write(string(msg))
}

func (t *Topology) buildIntroduceMessage() Message {
	return Message{
		AuthorName:    t.me,
		AuthorAddress: t.serverAddress,
		Kind:          "INTRO",
	}
}

func (t *Topology) buildPlainMessage(content string) Message {
	return Message{
		AuthorName:    t.me,
		AuthorAddress: t.serverAddress,
		Kind:          "PLAIN",
		Content:       content,
	}
}
