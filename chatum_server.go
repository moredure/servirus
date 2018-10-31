package main

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"io"
	"time"
)

type chatumServer struct {
	Bus
}

func NewChatumServer(bus Bus) chatum.ChatumServer {
	return &chatumServer{bus}
}

func (cs *chatumServer) Communicate(srv chatum.Chatum_CommunicateServer) error {
	ctx := srv.Context()
	username, err := ExtractUsernameFromContext(ctx)
	if err != nil {
		return err
	}

	sid := uuid.NewV4()

	cs.Add(username, sid, srv)
	defer cs.Remove(username, sid)

	ponger := make(chan bool, 1)
	defer close(ponger)

	closer := make(chan error, 1)
	defer close(closer)

	go func() {
		for {
			msg, err := srv.Recv()
			if err == io.EOF {
				closer <- nil
				return
			}
			if err != nil {
				closer <- err
				return
			}
			switch msg.GetType() {
			case chatum.EventType_DEFAULT:
				go cs.BroadcastExceptUUID(sid, NewMessage(username, msg.GetMessage()))
			case chatum.EventType_PONG:
				ponger <- true
			default:
			}
		}
	}()

	pinger := time.NewTicker(PingPongInterval)
	defer pinger.Stop()

	for {
		select {
		case err := <-closer:
			return err
		case <-pinger.C:
		}

		if err := srv.Send(NewPingMessage()); err != nil {
			return err
		}

		select {
		case <-time.After(PingPongTimeout):
			return PingPongTimeoutErr
		case <-ponger:
		}
	}
}
