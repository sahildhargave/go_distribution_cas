package p2p

//	type Handshaker interface {
//		Handshake() error
//	}
//
// type DefaultHandshaker struct{}

// var ErrInvalidHandshake = errors.New("Invalid Handshake")
// / HandshakeFunc.... ?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
