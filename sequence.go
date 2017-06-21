package rpg

import (
	"log"
	"os"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func NewGame(win *pixelgl.Window, difficulty int, leveltest ...string) {
	if leveltest == nil || len(leveltest) == 0 {
		leveltest = []string{""}
	}
	var (
		level = "1"
	)
	if leveltest[0] != "" {
		level = leveltest[0]
	}
	world := NewWorld(level, difficulty)
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
