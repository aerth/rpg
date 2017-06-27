package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aerth/rpg"
	astar "github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

var BOUNDS float64 = 700
var numbers = "0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
	// seed or random
	if len(os.Args) == 2 {
		if strings.HasPrefix(os.Args[1], "-h") {
			fmt.Println("Usage:")
			fmt.Println("\tmapgen [seed]")
			fmt.Println("Example:")
			fmt.Println("\tmapgen mycoolseed")

			os.Exit(111)
		}
		hashb := md5.Sum([]byte(os.Args[1]))
		hash := []byte(fmt.Sprintf("%x", hashb))
		var seed []byte
		for _, b := range hash {
			if bytes.IndexAny([]byte{b}, numbers) != -1 {
				log.Println(string(b), "is a number")
				seed = append(seed, b)
			} else {
				log.Println(string(b), "is a letter")
			}

		}
		worldseed, err := strconv.ParseInt(string(seed), 10, 64)
		if err != nil {
			log.Println(err)
		}
		rand.Seed(worldseed)
		log.Printf("Using world seed: %q -> %v", os.Args[1], worldseed)
		log.Printf("Hash: %q", string(hash))
	}
	// create maps dir if not exist
	os.Mkdir("maps", 0755)
}

func main() {
	olist := GenerateMap()
	SaveMap(olist)
}

func GenerateMap() []rpg.Object {
	var olist []rpg.Object
	t := rpg.O_TILE
	for i := 0; i < 100; i++ {

		currentThing := 20 // grass
		t = rpg.O_TILE
		if i%3 == 0 {
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
			if rpg.GetObjects(olist, obj.Loc) == nil {
				olist = append(olist, obj)
			}
		}

	}

	// make dummy world for path finding
	world := new(rpg.World)
	world.Tiles = rpg.GetTiles(olist)

	// detect islands, make bridges
	oldlist := olist
	olist = nil
	spot := world.Tile(rpg.FindRandomTile(oldlist))
	for _, o := range oldlist {
		o.W = world
		_, _, found := astar.Path(o, spot)
		if o.Type == rpg.O_TILE && !found {
			log.Println("found island tile", o)
		} else {
			olist = append(olist, o)
		}
	}

	// fill in with water blocks
	waterworld := rpg.DrawPatternObject(53, rpg.O_BLOCK, pixel.R(-BOUNDS, -BOUNDS, BOUNDS, BOUNDS), 0)
	for _, water := range waterworld {
		if rpg.GetObjects(olist, water.Loc) == nil {
			olist = append(olist, water)
		}
	}

	return olist
}

func SaveMap(olist []rpg.Object) {

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
	step := 32.00
	f := float64(rand.Intn(int(BOUNDS)))
	f = math.Floor(f)
	f = float64(int(f/step)) * step
	switch rand.Intn(2) {
	case 0:
		f = -f
	default:
	}

	return f

}
