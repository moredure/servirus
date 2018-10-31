package main

import "github.com/mikefaraponov/chatum"

type Client struct {
	chatum.Chatum_CommunicateServer
	Username string
}
