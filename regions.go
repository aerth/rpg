package rpg

import "log"

type Region struct {
	Name     string
	ID       int
	Portals  map[string]int // region id
	Map      []Object
	DObjects []*DObject
}

func (r Region) String() string {
	if r.Name == "" {
		return "unnamed region"
	}
	return "region " + r.Name
}

func (w *World) NewRegion(name string, mapobjects []Object) *Region {
	r := new(Region)
	r.ID = len(w.Regions)
	r.Map = mapobjects
	r.Name = name
	w.Regions = append(w.Regions, r)
	log.Println("added new region:", r)
	return r
}
