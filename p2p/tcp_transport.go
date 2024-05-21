package p2p

import (
	"net"
	"sync"
	"fmt"
)

// TCPPeers represenr the remote node ove a TCP established connection.
type TCPPeer struct{
	conn  net.Conn

	// if we dial and retrieve a conn => outbound == true
	// if we accept  and retrive a conn -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn :  conn,
		outbound : outbound,
	}
}
// TCPTransport represents a TCP transport layer for peer-to-peer communication.
type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// NewTCPTransport creates a new TCPTransport.
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}


func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener , err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}


func (t *TCPTransport) startAcceptLoop(){
	for{
		conn , err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn){
	fmt.Printf("new incomming connection %+V\n", conn)
}