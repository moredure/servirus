package server

import (
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/servirus/bus"
	"github.com/mikefaraponov/servirus/common"
)

type chatumServer struct {
	bus.Bus
}

func (cs *chatumServer) Communicate(srv chatum.Chatum_CommunicateServer) error {
	details, err := common.ExtractClientDetails(srv)
	if err != nil {
		return err
	}
	client := cs.Add(details)
	defer client.Close()
	go client.Listen()
	return client.PingPong()
}

func NewChatumServer(bus bus.Bus) chatum.ChatumServer {
	return &chatumServer{bus}
}
