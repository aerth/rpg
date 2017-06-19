package rpg

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Button struct {
	Name  string
	Frame pixel.Rect
}

var version = "0.0.91"

func TitleMenu(w *World, win *pixelgl.Window) {
	text := NewText(40)

	fmt.Fprintf(text, "AERPG v%s\nPRESS ENTER", version)
	dot := pixel.V(30, 400)
	text.Color = colornames.White
	text.Dot = dot
	text.Orig = text.Dot
	win.Clear(colornames.Black)
	text.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	win.SetTitle("AERPG (https://github.com/aerth/rpg)")
	log.Println("AERPG (https://github.com/aerth/rpg)")
	tick := time.Tick(time.Second)
	var frames = 0
	for !win.Closed() {

		frames++
		/*imd := imdraw.New(nil)
		imd.Color = colornames.Black
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(0)
		imd.Color = colornames.White
		imd.Push(pixel.V(30, 30), pixel.V(130, 130))
		imd.Rectangle(0)
		imd.Draw(win)*/

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			w.Exit(0)
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) || win.JustPressed(pixelgl.KeyEnter) {
			break
		}
		win.Update()
		select {
		case <-tick:
			win.SetTitle(fmt.Sprintf("AERPG (https://github.com/aerth/rpg) %v fps", frames))
			frames = 0
		default:

		}
	}
	log.Println("thanks for playing")
}

var textmatrix = pixel.IM.Moved(pixel.V(10, 580))

func (w *World) IsButton(buttons []Button, point pixel.Vec) (Button, func(win pixel.Target, world *World), bool) {

	for _, button := range buttons {
		if button.Frame.Contains(point) {
			switch button.Name {
			case "manastorm", "magic":
				return button, func(win pixel.Target, world *World) {
					world.Action(w.Char, w.Char.Rect.Center(), Magic)
				}, true
			default:
				return button, func(win pixel.Target, world *World) {
					world.Message(fmt.Sprintf("Bad button %s", point))
				}, true
			case "reset":
				return button, func(win pixel.Target, world *World) {
					world.Char.ResetLocation()
					//					world.Char.Inventory = []Item{}
				}, true

			}
		}
	}

	return Button{}, nil, false
}

func (w *World) Exit(code int) {
	os.Exit(code)
}
