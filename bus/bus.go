package bus

import (
	"github.com/mikefaraponov/chatum"
	"github.com/satori/go.uuid"
)

type Bus interface {
	BroadcastExceptUUID(uuid.UUID, *chatum.ServerSideEvent)
	BroadcastExceptUsername(*chatum.ServerSideEvent)
	Add(chatum.Chatum_CommunicateServer) (*Connection, error)
	Remove(*Connection)
}
