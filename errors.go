package main

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	PortEnvErr         = errors.New("$PORT is not set")
	UsernameMissingErr = status.Error(codes.Unauthenticated, "username is missing")
	PingPongTimeoutErr = errors.New("health check failed")
)
