package rpg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/colornames"

	"github.com/aerth/rpc/librpg/common"
	"github.com/aerth/rpg/assets"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func NewGame(win *pixelgl.Window, difficulty int, leveltest string, worldseed string) {
	world := NewWorld(leveltest, difficulty, worldseed)
	if world == nil {
		log.Println("bad world")
		os.Exit(111)
	}
	// have window, have world
	if err := world.drawTiles("tileset.png"); err != nil {
		log.Println("bad tiles", err)
		os.Exit(111)
	}

	if TitleMenu(win) {
		os.Exit(0)
	}
	var camPos = pixel.V(0, 0)
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X--
			log.Println(camPos)
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X++
			log.Println(camPos)
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y++
			log.Println(camPos)
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y--
			log.Println(camPos)
		}
		win.SetMatrix(pixel.IM.Moved(camPos))
		win.Clear(colornames.Green)
		world.Draw(win)
		win.Update()
	}
	os.Exit(111)
}

func (w *World) LoadMapFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return w.loadmap(b)
}
func (w *World) LoadMap(path string) error {
	b, err := assets.Asset(path)
	if err != nil {
		return err
	}
	return w.loadmap(b)
}
func (w *World) loadmap(b []byte) error {
	var things = []common.Object{}
	err := json.Unmarshal(b, &things)
	if err != nil {
		return fmt.Errorf("invalid map: %v", err)
	}
	return w.injectMap(things)
}

func (w *World) InjectMap(things []common.Object) error {
	return w.injectMap(things)
}

func (w *World) injectMap(things []common.Object) error {
	total := len(things)
	for i, t := range things {
		t.W = w
		t.Rect = common.DefaultSpriteRectangle.Moved(t.Loc)
		switch t.SpriteNum {
		case 53: // water
			t.Type = common.O_BLOCK
		default:
		}

		switch t.Type {
		case common.O_BLOCK:
			//log.Printf("%v/%v block object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
			w.Blocks = append(w.Blocks, t)
		case common.O_TILE:
			//log.Printf("%v/%v tile object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
			w.Tiles = append(w.Tiles, t)

		default: //
			log.Printf("%v/%v skipping bad object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
		}
	}
	log.Printf("map has %v blocks, %v tiles", len(w.Blocks), len(w.Tiles))
	if len(w.Blocks) == 0 && len(w.Tiles) == 0 {
		return fmt.Errorf("invalid map")
	}
	return nil
}
