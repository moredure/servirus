package server

import (
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/servirus/bus"
)

type chatumServer struct {
	bus.Bus
}

func (cs *chatumServer) Communicate(srv chatum.Chatum_CommunicateServer) error {
	if client, err := cs.Add(srv); err != nil {
		return err
	} else {
		return client.Operate()
	}
}

func NewChatumServer(bus bus.Bus) chatum.ChatumServer {
	return &chatumServer{bus}
}
