package main

import (
	"github.com/mikefaraponov/chatum"
	"io"
	"time"
)

type Client struct {
	Bus
	*ChatumClientDetails
	ponger chan bool
	closer chan error
	pinger *time.Ticker
}

func (c *Client) ListenEvents() {
	for {
		msg, err := c.Recv()
		if err == io.EOF {
			c.closer <- nil
			return
		}
		if err != nil {
			c.closer <- err
			return
		}
		switch msg.GetType() {
		case chatum.EventType_DEFAULT:
			go c.BroadcastExceptUUID(c.Id, NewMessage(c.Username, msg.GetMessage()))
		case chatum.EventType_PONG:
			c.ponger <- true
		default:
		}
	}
}

func (c *Client) HealthCheck() error {
	for {
		select {
		case err := <-c.closer:
			return err
		case <-c.pinger.C:
		}

		if err := c.Send(NewPingMessage()); err != nil {
			return err
		}

		select {
		case <-time.After(PingPongTimeout):
			return PingPongTimeoutErr
		case <-c.ponger:
		}
	}
}

func (c *Client) Close() {
	close(c.ponger)
	close(c.closer)
	c.pinger.Stop()
	c.Remove(c)
}

func NewClient(b Bus, d *ChatumClientDetails) *Client {
	return &Client{
		Bus:                 b,
		ChatumClientDetails: d,
		pinger:              time.NewTicker(PingPongInterval),
		ponger:              make(chan bool, 1),
		closer:              make(chan error, 1),
	}
}
