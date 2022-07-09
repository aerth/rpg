package main

import (
	"github.com/aerth/rpc/librpg/maps"
)

func main() {
	olist := maps.GenerateMap("")
	maps.SaveMap(olist)
}
