package rpg

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func InventoryLoop(win *pixelgl.Window, world *World) {
	text := NewText(28)
	text.WriteString("\tGAME PAUSED\n")
	text.WriteString("\tESC or any key to return\n\tUse PGUP/PGDOWN to scroll\n\n")
	text.WriteString(world.Char.Stats.String())
	text.WriteString(fmt.Sprintf("Level %v\nHealth: %v\nMana: %v\nXP: %v/%v", world.Char.Level, world.Char.Health, world.Char.Mana, world.Char.Stats.XP, world.Char.NextLevel()))
	for _, item := range world.Char.Inventory {
		if item.Effect != nil {
			text.WriteString(fmt.Sprintf("\n Effects: %q", item.Name))

		}
	}

	text.WriteString("\n\n===INVENTORY===\n" + FormatItemList(world.Char.Inventory))
	var page = 500.00
	for !win.Closed() {

		// controls
		switch {
		default:
		case win.JustPressed(pixelgl.KeyPageUp) || win.JustPressed(pixelgl.KeyUp):
			page -= 100
			log.Println(page)
			if page < 500 {
				page = 500
			}
		case win.JustPressed(pixelgl.KeyPageDown) || win.JustPressed(pixelgl.KeyDown):
			page += 100

			log.Println(page)
		}

		win.Clear(colornames.Black)
		imd := imdraw.New(nil)
		imd.Color = colornames.Green
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(30)
		imd.Color = colornames.Yellow
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(20)
		imd.Color = colornames.Red
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(10)
		imd.Draw(win)

		// break loop
		if win.Typed() != "" || win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyEnter) {
			break
		}

		// draw text
		text.Draw(win, pixel.IM.Moved(pixel.V(30, page)))

		// update window
		win.Update()

		// set title
		win.SetTitle("AERPG inventory")
	}
}
