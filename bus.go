package main

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"sync"
)

type (
	Bus interface {
		BroadcastExceptUUID(uuid.UUID, *chatum.ServerSideEvent)
		BroadcastExceptUsername(username, message string)
		Add(string, uuid.UUID, chatum.Chatum_CommunicateServer)
		Remove(string, uuid.UUID)
	}
	bus struct {
		sync.Mutex
		numberOfClientsByUsername map[string]int
		clientsById               map[uuid.UUID]Client
	}
)

func NewBus() Bus {
	return &bus{
		numberOfClientsByUsername: make(map[string]int),
		clientsById:               make(map[uuid.UUID]Client),
	}
}

func (b *bus) BroadcastExceptUsername(username, message string) {
	for _, client := range b.clientsById {
		if client.Username == username {
			continue
		}
		client.Send(NewMessage(username, message))
	}
}

func (b *bus) BroadcastExceptUUID(uid uuid.UUID, msg *chatum.ServerSideEvent) {
	for id, client := range b.clientsById {
		if id == uid {
			continue
		}
		client.Send(msg)
	}
}

func (b *bus) Add(username string, uid uuid.UUID, srv chatum.Chatum_CommunicateServer) {
	b.Lock()
	defer b.Unlock()
	if b.numberOfClientsByUsername[username] == 0 {
		b.BroadcastExceptUsername(username, "I am online!")
	}
	b.numberOfClientsByUsername[username] += 1
	b.clientsById[uid] = Client{srv, username}
}

func (b *bus) Remove(username string, uid uuid.UUID) {
	b.Lock()
	defer b.Unlock()
	b.numberOfClientsByUsername[username] -= 1
	if b.numberOfClientsByUsername[username] == 0 {
		b.BroadcastExceptUsername(username, "I am offline!")
	}
	delete(b.clientsById, uid)
}
