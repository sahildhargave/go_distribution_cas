package p2p

import "net"

// RPC represent any orbitary data that is being sent over the
// each transport between two node in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
