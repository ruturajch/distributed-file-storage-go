package p2p

import "net"

// Message represents the message holds any arbitrary data
// each transport between two nodes in network
type RPC struct {
	From    net.Addr
	Payload []byte
}
