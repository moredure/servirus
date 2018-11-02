package server

import (
	"net"
	"github.com/mikefaraponov/servirus/common"
)

func NewListener(port common.ServerPort) (net.Listener, error) {
	return net.Listen("tcp", string(port))
}
