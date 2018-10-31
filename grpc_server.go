package main

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func NewListener(port ServerPort) (net.Listener, error) {
	return net.Listen("tcp", string(port))
}

type GRPCParams struct {
	fx.In
	Server       *grpc.Server
	Listener     net.Listener
	ChatumServer chatum.ChatumServer
}

func Register(lc fx.Lifecycle, p GRPCParams) {
	chatum.RegisterChatumServer(p.Server, p.ChatumServer)
	reflection.Register(p.Server)
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := p.Server.Serve(p.Listener); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			p.Server.GracefulStop()
			return nil
		},
	})
}
