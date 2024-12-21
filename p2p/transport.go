package p2p

// Peer is interface that represents the remote connection
type Peer interface {
}

// Transport is anything that handeles the
// communication between the nodes in the network this can be of form TCP, UDP websockets...
type Transport struct {
	ListenAndAccept func() error
}
