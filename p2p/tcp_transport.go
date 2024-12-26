package p2p

import (
	"fmt"
	"net"
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

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// consume implements which will return read only channel recieved from another peer
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
		fmt.Printf("New incoming connection\n %+v\n", conn)
		go t.handelConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handelConn(conn net.Conn) {
	var err error
	peer := NewTCPPeer(conn, true)

	defer func() {
		fmt.Printf("dropping peer connection %s\n", err)
		conn.Close()
	}()
	peer = NewTCPPeer(conn, true)
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("Handshake error: %s\n", err)
		conn.Close()
		return
	}
	// Read Loop
	msg := RPC{}
	for {
		err := t.Decoder.Decode(conn, &msg)

		if err == net.ErrClosed {
			return
		}
		if err != nil {
			fmt.Printf("TCP read error: %sn", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		t.rpcch <- msg
		fmt.Printf("Message: %+v\n", msg)
		// msg := buf[:n]
	}
}
