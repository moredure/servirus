package main

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"google.golang.org/grpc/metadata"
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

func ExtractUsernameFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md[UsernameField]) == 0 {
		return "", UsernameMissingErr
	}
	return md[UsernameField][0], nil
}
