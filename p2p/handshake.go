package p2p

// HandshakeFunc is a function
type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error { return nil }
