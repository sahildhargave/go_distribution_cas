package p2p

import (
	"fmt"
	"net"
)

// TCPPeers represenr the remote node ove a TCP established connection.
type TCPPeer struct {
	conn net.Conn

	// if we dial and retrieve a conn => outbound == true
	// if we accept  and retrive a conn -> outbound == false
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

//CLOSE IMPLEMENTS THE PEER INTERFACE

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// TCPTransport represents a TCP transport layer for peer-to-peer communication.
type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	handshakeFunc HandshakeFunc
	rpcch         chan RPC
}

// NewTCPTransport creates a new TCPTransport.
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
		handshakeFunc:    opts.HandshakeFunc, // Initialize the handshakeFunc
		// Initialize the peers map
	}
}

// Consume implements the transport interface , which will return read-only channel
// foe reading the incoming message recived from another peer

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch

}
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
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

//type Temp struct{}

//func (t *TCPTransport) handleConn(conn net.Conn) {
//	var err error
//
//	defer func() {
//		fmt.Printf("dropping peer connection: %s", err)
//		conn.Close()
//	}()
//	peer := NewTCPPeer(conn, true)
//
//	if err = t.handshakeFunc(peer); err != nil {
//		return
//	}
//
//	if t.OnPeer != nil {
//		if err := t.OnPeer(peer); err != nil {
//			return
//		}
//	}
//
//	//Read Loop
//	rpc := RPC{}
//	//buf := make([]byte, 2000)
//	for {
//		//n, err := conn.Read(buf)
//		//if err != nil{
//		//	fmt.Printf("TCP error: %s\n", err)
//		//}
//
//		err = t.Decoder.Decode(conn, &rpc)
//		//fmt.Println(reflect.TypeOf(err))
//		//panic(err)
//		//if err == net.ErrClosed{
//		//	return
//		//}
//		if err != nil {
//			// fmt.Printf("TCP read error: %s\n", err)
//			//continue
//			return
//		}
//		rpc.From = conn.RemoteAddr()
//		t.rpcch <- rpc
//	}
//
//}
func (t *TCPTransport) handleConn(conn net.Conn) {
    var err error

    defer func() {
        if err != nil {
            fmt.Printf("dropping peer connection: %s\n", err)
        }
        conn.Close()
    }()

    peer := NewTCPPeer(conn, true)

    if err = t.handshakeFunc(peer); err != nil {
        fmt.Printf("handshake error: %s\n", err)
        return
    }

    if t.OnPeer != nil {
        if err := t.OnPeer(peer); err != nil {
            fmt.Printf("OnPeer error: %s\n", err)
            return
        }
    }

    // Read Loop
    rpc := RPC{}
    for {
        err = t.Decoder.Decode(conn, &rpc)
        if err != nil {
            fmt.Printf("TCP read error: %s\n", err)
            return
        }
        rpc.From = conn.RemoteAddr()
        t.rpcch <- rpc
    }
}
