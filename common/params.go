package common

import (
	"go.uber.org/fx"
	"os"
)

type ServerPort string

func NewServerPort() fx.Option {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return fx.Error(PortEnvErr)
	}
	return fx.Provide(func() ServerPort {
		return ServerPort(":" + port)
	})
}
