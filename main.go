package main

import (
	"fmt"
	"log"

	"github.com/anthdm/foreverstore/p2p"
)

func OnPeer(peer p2p.Peer) error {
	//fmt.Println("Failed to load peer func")
	peer.Close()
	return nil // Returning nil to indicate success, or an error to indicate failure.
}
func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
