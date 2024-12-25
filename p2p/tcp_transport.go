package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCP peer represents a remote node established connection.
type TCPPeer struct {
	//conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrive a conn => outbound == true
	// if we accept and retrive a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
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
		fmt.Printf("New incoming connection\n %+v\n", conn)
		go t.handelConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handelConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake error: %s\n", err)
		conn.Close()
		return
	}
	// Read Loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %sn", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		fmt.Printf("Message: %+v\n", msg)
		// msg := buf[:n]
	}
	// buf := make([]byte, 2000)
	// for {
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		fmt.Printf("TCP read error: %s\n", err)
	// 		break
	// 	}
	// 	fmt.Printf("Received: %s\n", buf[:n])
	// }
}
