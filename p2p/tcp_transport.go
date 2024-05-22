package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeers represenr the remote node ove a TCP established connection.
type TCPPeer struct {
	conn net.Conn

	// if we dial and retrieve a conn => outbound == true
	// if we accept  and retrive a conn -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// TCPTransport represents a TCP transport layer for peer-to-peer communication.
type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc HandshakeFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// NewTCPTransport creates a new TCPTransport.
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		handshakeFunc: func(any) error { return nil },
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incomming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.handshakeFunc(conn); err != nil {

	}
	msg := &Temp{}

	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}

	}
	fmt.Printf("new incomming connection %+v\n", peer)
}
