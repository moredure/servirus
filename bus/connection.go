package bus

import (
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/servirus/common"
	"io"
	"time"
)

type Connection struct {
	*common.ChatumClientDetails
	bus Bus
	ponger chan bool
	closer chan error
	pinger *time.Ticker
	messages chan *chatum.ClientSideEvent
}

func (c *Connection) GetMessages() {
	for {
		msg, err := c.Recv()
		if err == io.EOF {
			c.closer <- nil
			return
		} else if err != nil {
			c.closer <- err
			return
		}
		c.messages <- msg
	}
}

func (c *Connection) Listen() {
	for msg := range c.messages {
		switch msg.GetType() {
		case chatum.EventType_DEFAULT:
			go c.BroadcastExceptSelfUUID(msg.GetMessage())
		case chatum.EventType_PONG:
			c.ponger <- true
		default:
		}
	}
}

func (c *Connection) Operate() error {
	defer c.Close()
	go c.GetMessages()
	go c.Listen()
	return c.PingPong()
}

func (c *Connection) PingPong() error {
	for {
		select {
		case err := <-c.closer:
			return err
		case <-c.pinger.C:
		}

		go c.Send(common.NewPingMessage())

		select {
		case <-time.After(common.PingPongTimeout):
			return common.PingPongTimeoutErr
		case <-c.ponger:
		}
	}
}

func (c *Connection) BroadcastExceptSelfUsername(msg string) {
	c.bus.BroadcastExceptUsername(c.newMessage(msg))
}

func (c *Connection) BroadcastExceptSelfUUID(msg string) {
	c.bus.BroadcastExceptUUID(c.Id, c.newMessage(msg))
}

func (c *Connection) Close() {
	close(c.ponger)
	close(c.closer)
	close(c.messages)
	c.pinger.Stop()
	c.bus.Remove(c)
}

func (c *Connection) newMessage(msg string) *chatum.ServerSideEvent {
	return common.NewMessage(c.Username, msg)
}

func NewClient(b Bus, d *common.ChatumClientDetails) *Connection {
	return &Connection{
		ChatumClientDetails: d,
		bus:                 b,
		messages:            make(chan *chatum.ClientSideEvent),
		pinger:              time.NewTicker(common.PingPongInterval),
		ponger:              make(chan bool, 1),
		closer:              make(chan error, 1),
	}
}
