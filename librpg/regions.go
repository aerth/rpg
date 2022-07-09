package rpg

import (
	"log"

	"github.com/aerth/rpc/librpg/common"
)

type Region struct {
	Name     string
	ID       int
	Portals  map[string]int // region id
	Map      []common.Object
	DObjects []*DObject
}

func (r Region) String() string {
	if r.Name == "" {
		return "unnamed region"
	}
	return "region " + r.Name
}

func (w *World) NewRegion(name string, mapobjects []common.Object) *Region {
	r := new(Region)
	r.ID = len(w.Regions)
	r.Map = mapobjects
	r.Name = name
	w.Regions = append(w.Regions, r)
	log.Println("added new region:", r)
	return r
}
