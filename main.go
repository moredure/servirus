//go:generate protoc -I ../chatum --go_out=plugins=grpc:../chatum ../chatum/chat.proto

package main

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		NewServerPort(),
		fx.Provide(NewChatumServer),
		fx.Provide(NewListener),
		fx.Provide(grpc.NewServer),
		fx.Provide(NewBus),
		fx.Invoke(Register),
	).Run()
}
