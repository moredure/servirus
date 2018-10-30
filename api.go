package main

import (
	"context"
	"errors"
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
	"io"
	"time"
)

type chatumServer struct {
	bus Bus
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
	cs.bus.Add(username, sid, srv)
	defer cs.bus.Remove(username, sid)
	var wg errgroup.Group
	pinger := time.NewTicker(PingPongInterval)
	defer pinger.Stop()
	ponger := make(chan bool, 1)
	defer close(ponger)
	wg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			msg, err := srv.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				continue
			}
			switch msg.GetType() {
			case chatum.EventType_DEFAULT:
				go cs.bus.BroadcastExceptUUID(sid, &chatum.ServerSideEvent{
					Username: username,
					Message:  msg.Message,
				})
			case chatum.EventType_PONG:
				ponger <- true
			}
		}
	})
	wg.Go(func() error {
		for {
			<-pinger.C
			srv.Send(&chatum.ServerSideEvent{
				Type: chatum.EventType_PING,
			})
			select {
			case <-time.After(PingPongTimeout):
				return errors.New("health check failed")
			case <-ponger:
			}
		}
	})
	return wg.Wait()
}

func ExtractUsernameFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md[UsernameField]) == 0 {
		return "", UsernameMissingErr
	}
	return md[UsernameField][0], nil
}
