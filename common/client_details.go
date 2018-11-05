package common

import (
	"context"
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
)

type ChatumClientDetails struct {
	chatum.Chatum_CommunicateServer
	Username string
	Id       uuid.UUID
}

func ExtractClientDetails(srv chatum.Chatum_CommunicateServer) (*ChatumClientDetails, error) {
	ctx := srv.Context()
	username, err := extractUsernameFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &ChatumClientDetails{
		Chatum_CommunicateServer: srv,
		Username:                 username,
		Id:                       uuid.NewV4(),
	}, nil
}

func extractUsernameFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md[UsernameField]) == 0 {
		return "", UsernameMissingErr
	}
	return md[UsernameField][0], nil
}
