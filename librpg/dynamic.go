package rpg

import (
	"time"

	"github.com/aerth/rpc/librpg/common"
)

// doors
// loot

type DObjectType int

type DObject struct {
	Object   common.Object
	Contains []Item
	Until    time.Time `json:"-"`
	Type     DObjectType
}

const (
	D_NIL DObjectType = iota
	D_LOOT
	D_DOOR
	D_PORTAL
)
