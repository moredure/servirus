package bus

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/mikefaraponov/servirus/common"
)

type bus struct {
	sync.Mutex
	numberOfClientsByUsername map[string]int
	clientsById               map[*Connection]bool
}

func (b *bus) BroadcastExceptUsername(msg *chatum.ServerSideEvent) {
	for client := range b.clientsById {
		if client.Username == msg.GetUsername() {
			continue
		}
		go client.Send(msg)
	}
}

func (b *bus) BroadcastExceptUUID(uid uuid.UUID, msg *chatum.ServerSideEvent) {
	for client := range b.clientsById {
		if client.Id == uid {
			continue
		}
		go client.Send(msg)
	}
}

func (b *bus) Add(srv chatum.Chatum_CommunicateServer) (*Connection, error) {
	b.Lock()
	defer b.Unlock()
	d, err := common.ExtractClientDetails(srv)
	if err != nil {
		return nil, err
	}
	client := NewClient(b, d)
	if b.numberOfClientsByUsername[client.Username] == 0 {
		client.BroadcastExceptSelfUsername("I am online!")
	}
	b.numberOfClientsByUsername[client.Username] += 1
	b.clientsById[client] = true
	return client, nil
}

func (b *bus) Remove(c *Connection) {
	b.Lock()
	defer b.Unlock()
	b.numberOfClientsByUsername[c.Username] -= 1
	if b.numberOfClientsByUsername[c.Username] == 0 {
		c.BroadcastExceptSelfUsername("I am offline!")
	}
	delete(b.clientsById, c)
}

func NewBus() Bus {
	return &bus{
		numberOfClientsByUsername: make(map[string]int),
		clientsById:               make(map[*Connection]bool),
	}
}
