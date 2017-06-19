// Copyright 2017 aerth <aerth@riseup.net>

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/aerth/rpg"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var LEVEL string

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("> ")
	if len(os.Args) != 2 {
		fmt.Println("Which map name?")
		os.Exit(111)
	}
	LEVEL = os.Args[1]

}

var (
	IM = pixel.IM
	ZV = pixel.ZV
)

func loadSpriteSheet() (pixel.Picture, []*pixel.Sprite) {
	spritesheet, err := rpg.LoadPicture("sprites/tileset.png")
	/* 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
	 * 1
	 * 2
	 * ...
	 * 16
	 */
	if err != nil {

		panic(err)

	}
	var sheetFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			sheetFrames = append(sheetFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var spritemap = []*pixel.Sprite{}
	for i := 0; i < len(sheetFrames); i++ {
		x := i
		spritemap = append(spritemap, pixel.NewSprite(spritesheet, sheetFrames[x]))
	}
	log.Println(len(spritemap), "sprites loaded")
	return spritesheet, spritemap
}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:     "AERPG mapedit",
		Bounds:    pixel.R(0, 0, 800, 600),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	imd := imdraw.New(nil)
	canvas := pixelgl.NewCanvas(pixel.R(-1000.00, -1000.00, 1000.00, 1000.00))
	var oldthings = []rpg.Object{}

	if b, err := ioutil.ReadFile(LEVEL); err == nil {
		err = json.Unmarshal(b, &oldthings)
		if err != nil {
			panic(err)
		}

	}

	// convert old map to new map style (object types)
	var things = []rpg.Object{}
	for _, v := range oldthings {
		if v.P.Tile || v.Type == rpg.O_TILE {
			v.P.Tile = true
			v.P.Block = false
			v.Type = rpg.O_TILE
		}
		if v.P.Block || v.Type == rpg.O_BLOCK {
			v.P.Block = true
			v.P.Tile = false
			v.Type = rpg.O_BLOCK

		}
		things = append(things, v)
	}

	spritesheet, spritemap := loadSpriteSheet()

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	start := time.Now()
	second := time.Tick(time.Second)
	tick := time.Tick(time.Millisecond * 200)
	last := start
	frames := 0

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)
	currentThing := 20 // 20 is grass,  0 should be transparent sprite
	text := rpg.NewTextSmooth(14)
	rpg.DrawText(win.Bounds(), text, win, "ENTER=save LEFT=block RIGHT=tile SHIFT=batch CAPS=highlight U=undo R=redo 4=turbo B=dontreplace")
	cursor := rpg.GetCursor(0)
	undobuffer := []rpg.Object{}
	var turbo = false
	var highlight = true
	var box pixel.Rect
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		_ = dt
		last = time.Now()
		frames++

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		//cam := pixel.IM.Moved(win.Bounds().Center()).Scaled(camPos, camZoom)
		//		cam := pixel.IM.Moved(win.Bounds().Center()).Moved(camPos.Scaled(-camZoom))
		//.Scaled(pixel.ZV, camZoom)
		//		cam := pixel.IM.Scaled(pixel.ZV, camZoom).Moved(win.Bounds().Center()).Moved(camPos.Scaled(-1))
		win.SetMatrix(cam)
		mouse := cam.Unproject(win.MousePosition())
		//.Add(pixel.V(16, 16))
		snap := 32.00
		mouse.X = float64(int(mouse.X/snap)) * snap
		mouse.Y = float64(int(mouse.Y/snap)) * snap // 'snap to grid'

		if win.JustPressed(pixelgl.Key4) {
			turbo = !turbo
		}
		if win.JustPressed(pixelgl.KeyCapsLock) {
			highlight = !highlight
		}

		if turbo {
			dt *= 8
		}

		if win.JustPressed(pixelgl.KeyU) {
			undobuffer = append(undobuffer, things[len(things)-1])
			things = things[:len(things)-1]
		}
		if win.JustPressed(pixelgl.KeyR) {
			if len(undobuffer) > 0 {
				things = append(things, undobuffer[len(undobuffer)-1])
				undobuffer = undobuffer[:len(undobuffer)-1]
			} else {
				log.Println("no undo buffer")
			}
		}

		deleteThing := func(loc pixel.Vec) []rpg.Object {
			var newthings []rpg.Object
			for _, thing := range things {
				if thing.Rect.Contains(mouse) {
					log.Println("deleting:", thing)
				} else {

					newthings = append(newthings, thing)
				}
			}
			return newthings
		}
		if win.Pressed(pixelgl.KeySpace) {
			things = deleteThing(mouse)
		}
		var replace bool
		if !win.Pressed(pixelgl.KeyB) {
			replace = true
		}
		// draw big patch of grass
		if win.JustPressed(pixelgl.MouseButtonMiddle) {
			box.Min.X = mouse.X
			box.Max.Y = mouse.Y
		}
		if win.Pressed(pixelgl.MouseButtonMiddle) {
			box.Min.Y = mouse.Y
			box.Max.X = mouse.X
		}
		// need real tiles though

		if win.Pressed(pixelgl.KeyLeftShift) && win.Pressed(pixelgl.MouseButtonRight) ||
			win.JustPressed(pixelgl.MouseButtonRight) {
			thing := rpg.NewTile(mouse)
			thing.SpriteNum = currentThing
			log.Println("Stamping Tile", mouse, thing.SpriteNum)
			if replace {
				things = append(deleteThing(mouse), thing)
			} else {
				things = append(things, thing)

			}
		}
		if win.Pressed(pixelgl.KeyLeftShift) && win.Pressed(pixelgl.MouseButtonLeft) ||
			win.JustPressed(pixelgl.MouseButtonLeft) {
			thing := rpg.NewBlock(mouse)
			thing.SpriteNum = currentThing
			log.Println("Stamping Block", mouse, thing.SpriteNum)
			if replace {
				things = append(deleteThing(mouse), thing)
			} else {
				things = append(things, thing)
			}
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			b, _ := json.Marshal(things)
			os.Rename(LEVEL, LEVEL+".old")
			ioutil.WriteFile(LEVEL, b, 0644)
			log.Println(LEVEL + " map saved")
		}
		if win.JustPressed(pixelgl.KeyPageUp) {
			currentThing++
			if currentThing > len(spritemap)-1 {
				currentThing = 0
			}
		}
		if win.JustPressed(pixelgl.KeyPageDown) {
			currentThing--
			if currentThing <= 0 {
				currentThing = len(spritemap) - 1
			}
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			camPos.Y += camSpeed * dt
		}

		//	canvas.Clear(pixel.Alpha(0))
		win.Clear(colornames.Green)
		canvas.Clear(colornames.Green)
		batch.Clear()
		for i := range things {
			if things[i].SpriteNum == 0 {
				things[i].Sprite = spritemap[2]
			}
			things[i].Draw(batch, spritesheet, spritemap)
			if highlight {
				things[i].Highlight(batch)
			}

		}

		batch.Draw(canvas)
		rpg.DrawPattern(canvas, spritemap[20], box, 100)
		if b := box.Size(); b.Len() != 0 {
			imd.Clear()
			imd.Color = pixel.RGB(0, 1, 0)
			imd.Push(box.Min, box.Max)
			imd.Rectangle(1)
			imd.Draw(win)
		}
		canvas.Draw(win, IM.Scaled(ZV, camZoom))
		spritemap[182].Draw(win, IM.Scaled(ZV, 2).Moved(pixel.V(16, 16)))
		win.SetMatrix(IM)
		spritemap[currentThing].Draw(win, IM.Scaled(ZV, 2).Moved(pixel.V(64, 64)).Moved(spritemap[0].Frame().Center()))
		text.Draw(win, IM.Moved(pixel.V(80, 10)))
		cursor.Draw(win, IM.Moved(win.MousePosition()).Moved(pixel.V(32, -32)))

		win.Update()

		select {
		default:
		case <-tick:
			//	things = append(things, NewThing(spritemap))

		}

		select {
		default: //
		case <-second:
			//	log.Println("Offset:", offset)
			log.Println("Last DT", dt)
			log.Println("FPS:", frames)
			log.Printf("things: %v", len(things))
			//log.Printf("dynamic things: %v", len(world.DObjects))
			frames = 0
		}
	}
}
func OldNewBlock(spritemap []*pixel.Sprite, r int) rpg.Object {
	return rpg.Object{
		Loc:       randomLocation(),
		SpriteNum: r,
		P: rpg.ObjectProperties{
			Block: true,
		},
	}
}

func NewThing(spritemap []*pixel.Sprite) rpg.Object {
	r := rand.Intn(len(spritemap))

	return rpg.Object{
		Loc:       randomLocation(),
		SpriteNum: r,
	}
}

func main() {
	pixelgl.Run(run)
}
func randomLocation() pixel.Vec {
	return pixel.Vec{float64(rand.Intn(1000)), float64(rand.Intn(1000))}
}
