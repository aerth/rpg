package rpg

import (
	"fmt"

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
	imd := imdraw.New(nil)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)
	text.WriteString("\n\n===INVENTORY===\n" + FormatItemList(world.Char.Inventory))
	var page = 600.00
	if win.Bounds().Max.Y < 800 {
		page = 500
	}
	xpage := 30.00
	for !win.Closed() {

		// controls
		switch {
		default:
		case win.JustPressed(pixelgl.KeyPageUp) || win.JustPressed(pixelgl.KeyUp):
			page -= 100
			if page < 500 {
				page = 500
			}
		case win.JustPressed(pixelgl.KeyPageDown) || win.JustPressed(pixelgl.KeyDown):
			page += 100
		case win.JustPressed(pixelgl.KeyLeft):
			xpage += 100
		case win.JustPressed(pixelgl.KeyRight):
			xpage -= 100
		}

		win.Clear(colornames.Black)
		imd.Color = pixel.ToRGBA(colornames.Darkslategrey).Scaled(0.01)

		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(0)
		imd.Color = colornames.Green
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(30)
		imd.Color = colornames.Yellow
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(20)
		imd.Color = colornames.Red
		imd.Push(pixel.V(0, 0), win.Bounds().Max)
		imd.Rectangle(10)

		// break loop
		if win.Typed() != "" || win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyEnter) {
			break
		}

		// draw text
		imd.Draw(win)
		batch.Draw(win)
		text.Draw(win, pixel.IM.Moved(pixel.V(xpage, page)))

		// update window
		win.Update()

		// set title
		win.SetTitle("AERPG inventory")
	}
}
