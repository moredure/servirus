package bus

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/mikefaraponov/servirus/common"
)

type (
	Bus interface {
		BroadcastExceptUUID(uuid.UUID, *chatum.ServerSideEvent)
		BroadcastExceptUsername(*chatum.ServerSideEvent)
		Add(*common.ChatumClientDetails) *Client
		Remove(*Client)
	}
	bus struct {
		sync.Mutex
		numberOfClientsByUsername map[string]int
		clientsById               map[uuid.UUID]*Client
	}
)

func (b *bus) BroadcastExceptUsername(msg *chatum.ServerSideEvent) {
	for _, client := range b.clientsById {
		if client.Username == msg.GetUsername() {
			continue
		}
		go client.Send(msg)
	}
}

func (b *bus) BroadcastExceptUUID(uid uuid.UUID, msg *chatum.ServerSideEvent) {
	for id, client := range b.clientsById {
		if id == uid {
			continue
		}
		go client.Send(msg)
	}
}

func (b *bus) Add(d *common.ChatumClientDetails) *Client {
	client := NewClient(b, d)

	b.Lock()
	defer b.Unlock()

	if b.numberOfClientsByUsername[client.Username] == 0 {
		client.BroadcastExceptSelfUsername("I am online!")
	}
	b.numberOfClientsByUsername[client.Username] += 1
	b.clientsById[client.Id] = client

	return client
}

func (b *bus) Remove(c *Client) {
	b.Lock()
	defer b.Unlock()

	b.numberOfClientsByUsername[c.Username] -= 1
	if b.numberOfClientsByUsername[c.Username] == 0 {
		c.BroadcastExceptSelfUsername("I am offline!")
	}
	delete(b.clientsById, c.Id)
}

func NewBus() Bus {
	return &bus{
		numberOfClientsByUsername: make(map[string]int),
		clientsById:               make(map[uuid.UUID]*Client),
	}
}
