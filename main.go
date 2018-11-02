//go:generate protoc -I ../chatum --go_out=plugins=grpc:../chatum ../chatum/chat.proto

package main

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"github.com/mikefaraponov/servirus/common"
	"github.com/mikefaraponov/servirus/server"
	"github.com/mikefaraponov/servirus/bus"
)

func main() {
	fx.New(
		common.NewServerPort(),
		fx.Provide(server.NewChatumServer),
		fx.Provide(server.NewListener),
		fx.Provide(grpc.NewServer),
		fx.Provide(bus.NewBus),
		fx.Invoke(server.Bootstrap),
	).Run()
}
