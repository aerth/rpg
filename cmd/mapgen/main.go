package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/aerth/rpc/librpg/common"
	"github.com/aerth/rpc/librpg/maps"
)

func SaveMap(olist []common.Object) {
	os.MkdirAll("maps", 0755)
	b, err := json.Marshal(olist)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := ioutil.TempFile("maps", "map")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("saved file:", f.Name())
}
func main() {
	olist := maps.GenerateMap("")
	SaveMap(olist)
}
