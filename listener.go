package main

import (
	"net"
)

func NewListener(port ServerPort) (net.Listener, error) {
	return net.Listen("tcp", string(port))
}
