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
		Add(*ChatumClientDetails) *Client
		Remove(*Client)
	}
	bus struct {
		sync.Mutex
		numberOfClientsByUsername map[string]int
		clientsById               map[uuid.UUID]*Client
	}
)

func NewBus() Bus {
	return &bus{
		numberOfClientsByUsername: make(map[string]int),
		clientsById:               make(map[uuid.UUID]*Client),
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

func (b *bus) Add(d *ChatumClientDetails) *Client {
	b.Lock()
	defer b.Unlock()
	if b.numberOfClientsByUsername[d.Username] == 0 {
		b.BroadcastExceptUsername(d.Username, "I am online!")
	}
	b.numberOfClientsByUsername[d.Username] += 1
	client := NewClient(b, d)
	b.clientsById[d.Id] = client
	return client
}

func (b *bus) Remove(c *Client) {
	b.Lock()
	defer b.Unlock()
	b.numberOfClientsByUsername[c.Username] -= 1
	if b.numberOfClientsByUsername[c.Username] == 0 {
		b.BroadcastExceptUsername(c.Username, "I am offline!")
	}
	delete(b.clientsById, c.Id)
}
