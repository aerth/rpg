package rpg

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/aerth/rpc/librpg/common"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Button struct {
	Name  string
	Frame pixel.Rect
}

var version = "0.0.95"

func Version() string {
	return "AERPG " + version
}

func TitleMenu(win *pixelgl.Window) (breakloop bool) {
	title := NewTitleText(64)
	text := NewText(36)
	dot := pixel.V(30, 400)
	title.Dot = dot
	title.Orig = title.Dot

	dot = pixel.V(30, 200)
	text.Dot = dot
	text.Orig = text.Dot

	fmt.Fprintf(title, "AERPG v%s\nPRESS ENTER", version)
	fmt.Fprintln(text, "https://github.com/aerth/rpg\n\nCTRL-Q to QUIT")

	win.Clear(colornames.Black)

	win.SetTitle("AERPG (https://github.com/aerth/rpg)")
	log.Println("AERPG (https://github.com/aerth/rpg)")
	log.Println("update often")

	tick := time.Tick(time.Second)
	var frames = 0
	for !win.Closed() {

		win.Clear(colornames.Black)
		title.Color = RandomColor()
		text.Draw(win, pixel.IM)
		title.Draw(win, pixel.IM)
		frames++
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return true
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			return false
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
	return true
}

func (w *World) IsButton(buttons []Button, point pixel.Vec) (Button, func(win pixel.Target, world *World), bool) {

	for _, button := range buttons {
		if button.Frame.Contains(point) {
			switch button.Name {
			case "manastorm":
				return button, func(win pixel.Target, world *World) {
					world.Action(w.Char, w.Char.Rect.Center(), ManaStorm, OUT)
				}, true
			case "magicbullet":
				return button, func(win pixel.Target, world *World) {
					world.Action(w.Char, w.Char.Rect.Center(), MagicBullet, world.Char.Dir)
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

func (c *Character) DrawBars(target pixel.Target, bounds pixel.Rect) {
	var barheight = 10.00
	var startY = 50.00
	imd := imdraw.New(nil)
	xp := float64(c.Stats.XP)
	next := float64(c.NextLevel())
	rect := bounds
	rect.Min.Y = startY
	rect.Max.Y = rect.Min.Y + barheight
	common.DrawBar(imd, colornames.Red, float64(c.Health), float64(255), rect)
	common.DrawBar(imd, colornames.Blue, float64(c.Mana), 255.00, rect.Moved(pixel.V(0, barheight+1)))
	common.DrawBar(imd, colornames.Purple, xp, next, rect.Moved(pixel.V(0, (barheight*2)+1)))
	imd.Draw(target)
}
