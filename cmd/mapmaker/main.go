// Copyright 2017 aerth <aerth@riseup.net>

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/image/colornames"

	rpg "github.com/aerth/rpc/librpg"
	"github.com/aerth/rpc/librpg/common"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var LEVEL string

func FlagInit() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("> ")
	if flag.NArg() != 1 {
		fmt.Println("Which map name?")
		os.Exit(111)
	}
	LEVEL = flag.Arg(0)

}

var convert = flag.Bool("danger", false, "convert old to new (experimental)")
var (
	IM = pixel.IM
	ZV = pixel.ZV
)

var helpText = "ENTER=save LEFT=block RIGHT=tile SHIFT=batch SPACE=del CAPS=highlight U=undo R=redo 4=turbo B=dontreplace"

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
	log.Println(spritemap[0].Frame())
	return spritesheet, spritemap
}

func run() {
	flag.Parse()
	FlagInit()
	cfg := pixelgl.WindowConfig{
		Title:     "AERPG mapedit",
		Bounds:    pixel.R(0, 0, 800, 600),
		Resizable: true,
		VSync:     false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	var oldthings = []common.Object{}
	if b, err := ioutil.ReadFile(LEVEL); err == nil {
		err = json.Unmarshal(b, &oldthings)
		if err != nil {
			panic(err)
		}
	}
	var things []common.Object
	for _, v := range oldthings {
		if *convert {
			log.Println("Converting")
			v.Type = common.O_TILE
			if v.SpriteNum == 53 && v.Type == common.O_TILE {
				v.Type = common.O_BLOCK
			}
		}

		v.Rect = common.DefaultSpriteRectangle.Moved(v.Loc)

		things = append(things, v)

	}

	spritesheet, spritemap := loadSpriteSheet()

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	start := time.Now()
	second := time.Tick(time.Second)
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
	fmt.Fprint(text, helpText)
	cursor := rpg.GetCursor(2)
	undobuffer := []common.Object{}
	var turbo = false
	var highlight = true
	var box pixel.Rect
	var replace = true
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		_ = dt
		last = time.Now()
		frames++

		//		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		// camera
		cam := pixel.IM.Scaled(pixel.ZV, camZoom).Moved(win.Bounds().Center()).Moved(camPos.Scaled(-camZoom))
		win.SetMatrix(cam)

		// snap to grid
		snap := 32.00 // 16 for half grid ?
		mouse := cam.Unproject(win.MousePosition())
		mouse.X = float64(int(mouse.X/snap)) * snap
		mouse.Y = float64(int(mouse.Y/snap)) * snap
		//		mouse.X = mouse.X - 16
		//		mouse.Y = mouse.Y - 16
		if win.JustPressed(pixelgl.Key4) {
			turbo = !turbo
			log.Println("turbo:", turbo)
		}
		if win.JustPressed(pixelgl.KeyCapsLock) {
			highlight = !highlight
			log.Println("highlight:", highlight)
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
				if !win.Pressed(pixelgl.KeyLeftShift) {
					undobuffer = undobuffer[:len(undobuffer)-1]
				}
			} else {
				log.Println("no undo buffer")
			}
		}

		deleteThing := func(loc pixel.Vec) []common.Object {
			var newthings []common.Object
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
		if win.JustPressed(pixelgl.KeyB) {
			replace = !replace
			log.Println("replace:", replace)
		}
		// draw big patch of grass
		if win.Pressed(pixelgl.KeyLeftControl) && (win.JustPressed(pixelgl.MouseButtonLeft) || win.JustPressed(pixelgl.MouseButtonRight)) {
			box.Min.Y = mouse.Y
			box.Min.X = mouse.X
		} else {
			if win.Pressed(pixelgl.KeyLeftShift) && win.Pressed(pixelgl.MouseButtonRight) ||
				win.JustPressed(pixelgl.MouseButtonRight) {
				thing := common.NewBlock(mouse)
				thing.SpriteNum = currentThing
				log.Println("Stamping Block", mouse, thing.SpriteNum)
				if replace {
					undobuffer = append(undobuffer, thing)
					things = append(deleteThing(mouse), thing)
				} else {
					things = append(things, thing)

				}
			}
			if win.Pressed(pixelgl.KeyLeftShift) && win.Pressed(pixelgl.MouseButtonLeft) ||
				win.JustPressed(pixelgl.MouseButtonLeft) {
				thing := common.NewTile(mouse)
				thing.SpriteNum = currentThing
				log.Println("Stamping Tile", mouse, thing.SpriteNum)
				if replace {
					undobuffer = append(undobuffer, thing)
					things = append(deleteThing(mouse), thing)
				} else {
					things = append(things, thing)
				}
			}
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			b, err := json.Marshal(things)
			if err != nil {
				panic(err)
			}
			os.Rename(LEVEL, LEVEL+".old")
			if err := ioutil.WriteFile(LEVEL, b, 0644); err != nil {
				log.Println(LEVEL + " map saved")
			}
		}
		if win.JustPressed(pixelgl.KeyPageUp) {
			currentThing++
			if currentThing > len(spritemap)-1 {
				currentThing = 0
			}
			log.Println("current sprite:", currentThing)
		}
		if win.JustPressed(pixelgl.KeyPageDown) {
			currentThing--
			if currentThing <= 0 {
				currentThing = len(spritemap) - 1
			}
			log.Println("current sprite:", currentThing)
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
		batch.Clear()

		batch.Draw(win)
		if b := box.Size(); b.Len() != 0 {
			if win.Pressed(pixelgl.KeyLeftControl) {
				if win.JustReleased(pixelgl.MouseButtonLeft) {
					box.Max = mouse
					box = box.Norm()
					log.Println("drawing rectangle:", box, currentThing)
					things = append(DeleteThings(things, box), common.DrawPatternObject(currentThing, common.O_TILE, box, 100)...)
				}
				if win.JustReleased(pixelgl.MouseButtonRight) {
					box.Max = mouse
					box = box.Norm()
					log.Println("drawing rectangle:", box, currentThing)
					things = append(DeleteThings(things, box), common.DrawPatternObject(currentThing, common.O_BLOCK, box, 100)...)

				}
			}
		}

		for i := range things {
			things[i].Draw(batch, spritesheet, spritemap, 0)
			if highlight {
				color := common.TransparentRed
				if things[i].Type == common.O_TILE {
					color = common.TransparentBlue
				}
				things[i].Highlight(batch, color)
			}
			if things[i].Rect.Contains(mouse) {
				things[i].Highlight(batch, common.TransparentPurple)

			}

		}

		batch.Draw(win)

		// draw player spawn
		spritemap[182].Draw(win, IM.Scaled(ZV, 2).Moved(pixel.V(8, 8))) // incorrect offset

		// return cam
		win.SetMatrix(IM)
		spritemap[currentThing].Draw(win, IM.Scaled(ZV, 2).Moved(pixel.V(64, 64)).Moved(spritemap[0].Frame().Center()))
		text.Draw(win, IM.Moved(pixel.V(10, 10)))
		//		cursor.Draw(win, IM.Moved(win.MousePosition()).Moved(pixel.V(32, -32)))

		cursor.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()).Moved(pixel.V(0, -32)))
		win.Update()

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

func main() {
	pixelgl.Run(run)
}

func DeleteThings(from []common.Object, at pixel.Rect) []common.Object {
	var cleaned []common.Object
	for _, o := range from {
		if !at.Contains(o.Loc) {
			cleaned = append(cleaned, o)
		}

	}
	return cleaned
}
