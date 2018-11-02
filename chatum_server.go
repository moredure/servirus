package main

import (
	"github.com/mikefaraponov/chatum"
)

type chatumServer struct {
	Bus
}

func (cs *chatumServer) Communicate(srv chatum.Chatum_CommunicateServer) error {
	details, err := ExtractClientDetails(srv)
	if err != nil {
		return err
	}
	client := cs.Add(details)
	defer client.Close()
	go client.ListenEvents()
	return client.HealthCheck()
}

func NewChatumServer(bus Bus) chatum.ChatumServer {
	return &chatumServer{bus}
}
