package common

import (
	"context"
	"go.uber.org/fx"
	"os"
	"testing"
)

const ExpectedValue = "1337"

func TestNewServerPort(t *testing.T) {
	t.Run("test when port env is set", func(t *testing.T) {
		os.Setenv("PORT", ExpectedValue)
		defer os.Unsetenv("PORT")
		var port ServerPort
		app := fx.New(
			NewServerPort(),
			fx.Populate(&port),
		)
		if err := app.Start(context.Background()); err != nil {
			t.Fatal(err)
		}
		defer app.Stop(context.Background())

		if string(port) != ":"+ExpectedValue {
			t.Fail()
		}
	})
	t.Run("test when port env is not set", func(t *testing.T) {
		os.Unsetenv("PORT")
		var port ServerPort
		app := fx.New(
			NewServerPort(),
			fx.Populate(&port),
		)
		if err := app.Start(context.Background()); err != PortEnvErr {
			t.Fatal(err)
		}
		defer app.Stop(context.Background())
	})
}
