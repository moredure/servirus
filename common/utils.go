package common

import (
	"github.com/mikefaraponov/chatum"
)

func NewMessage(username, message string) *chatum.ServerSideEvent {
	return &chatum.ServerSideEvent{
		Username: username,
		Message:  message,
	}
}

func NewPingMessage() *chatum.ServerSideEvent {
	return &chatum.ServerSideEvent{
		Type: chatum.EventType_PING,
	}
}
