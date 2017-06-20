package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/aerth/rpg"
	"github.com/faiface/pixel"
)

func main() {
	os.Mkdir("maps", 0755)
	var olist []rpg.Object
	t := rpg.O_TILE
	for i := 0; i < 100; i++ {

		currentThing := 20 // grass
		t = rpg.O_TILE
		if i%2 == 1 {
			currentThing = 53 // water
			t = rpg.O_BLOCK
		}
		xmin := randfloat()
		ymin := randfloat()
		xmax := randfloat()
		ymax := randfloat()
		box := pixel.R(xmin, ymin, xmax, ymax).Norm()
		log.Println(t, box)
		pattern := rpg.DrawPatternObject(currentThing, t, box, 100)
		for _, obj := range pattern {

			if len(rpg.GetTilesAt(olist, obj.Loc)) == 0 {
				olist = append(olist, obj)
			}
		}

	}

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

func randfloat() float64 {

	f := float64(rand.Intn(1500))
	switch rand.Intn(2) {
	case 0:
		f = -f
	default:
	}

	return f

}
