package rpg

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
			case "manastorm":
				return button, func(win pixel.Target, world *World) {
					world.Action(w.Char, w.Char.Rect.Center(), ManaStorm)
				}, true
			default:
				return button, func(win pixel.Target, world *World) {
					world.Message(fmt.Sprintf("Bad button %s", point))
				}, true
			case "reset":
				return button, func(win pixel.Target, world *World) {
					world.Reset()
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

func (c *Character) DrawBars(target pixel.Target) {
	imd := imdraw.New(nil)
	xp := float64(c.Stats.XP)
	next := float64(c.NextLevel())
	percent := xp / next

	// XP
	imd.Color = colornames.Purple
	imd.Push(pixel.V(10, 100))
	imd.Push(pixel.V(110, 120))
	imd.Rectangle(4)
	if xp > 0 {
		imd.Color = colornames.Purple
		imd.Push(pixel.V(10, 100))
		imd.Push(pixel.V(115*percent, 120))
		imd.Rectangle(0)
	}

	// HP
	pt := (110 * float64(c.Health) / 255)
	if c.Health == 0 {
		pt = 10
	}
	imd.Color = colornames.Red
	imd.Push(pixel.V(10, 130))
	imd.Push(pixel.V(110, 150))
	imd.Rectangle(4)

	imd.Color = colornames.Red
	imd.Push(pixel.V(10, 130))
	imd.Push(pixel.V(pt, 150))
	imd.Rectangle(0)

	// MP
	pt = (110 * float64(c.Mana) / 255)
	if c.Mana == 0 {
		pt = 10
	}
	imd.Color = colornames.Blue
	imd.Push(pixel.V(10, 160))
	imd.Push(pixel.V(110, 180))
	imd.Rectangle(4)
	imd.Color = colornames.Blue

	if pt > 10 {
		imd.Push(pixel.V(10, 160))
		imd.Push(pixel.V(pt, 180))
		imd.Rectangle(0)
	}

	imd.Draw(target)

}
