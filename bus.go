package main

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"sync"
	"fmt"
)

type (
	Bus interface {
		BroadcastExceptUUID(uuid.UUID, *chatum.ServerSideEvent)
		BroadcastExceptUsername(string, *chatum.ServerSideEvent)
		Add(string, uuid.UUID, chatum.Chatum_CommunicateServer)
		Remove(string, uuid.UUID)
	}
	bus struct {
		sync.Mutex
		numberOfClientsByUsername map[string]int
		clientsById               map[uuid.UUID]Client
	}
)

type Client struct {
	Srv      chatum.Chatum_CommunicateServer
	Username string
}

func NewBus() Bus {
	return &bus{
		numberOfClientsByUsername: make(map[string]int),
		clientsById:               make(map[uuid.UUID]Client),
	}
}

func (b *bus) BroadcastExceptUsername(username string, msg *chatum.ServerSideEvent) {
	for _, v := range b.clientsById {
		if v.Username == username {
			continue
		}
		v.Srv.Send(msg)
	}
}

func (b *bus) BroadcastExceptUUID(uid uuid.UUID, msg *chatum.ServerSideEvent) {
	for k, v := range b.clientsById {
		if k == uid {
			continue
		}
		v.Srv.Send(msg)
	}
}

func (b *bus) Add(username string, uid uuid.UUID, srv chatum.Chatum_CommunicateServer) {
	fmt.Println("add", username, uid)
	b.Lock()
	defer b.Unlock()
	if b.numberOfClientsByUsername[username] == 0 {
		b.BroadcastExceptUsername(username, &chatum.ServerSideEvent{
			Username: username,
			Message:  "I am online!",
		})
	}
	b.numberOfClientsByUsername[username] += 1
	b.clientsById[uid] = Client{srv, username}
}

func (b *bus) Remove(username string, uid uuid.UUID) {
	fmt.Println("remove", username, uid)
	b.Lock()
	defer b.Unlock()
	b.numberOfClientsByUsername[username] -= 1
	if b.numberOfClientsByUsername[username] == 0 {
		b.BroadcastExceptUsername(username, &chatum.ServerSideEvent{
			Username: username,
			Message:  "I am offline!",
		})
	}
	delete(b.clientsById, uid)
}
