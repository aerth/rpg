package rpg

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func UI() *text.Text {
	var atlas *text.Atlas
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textbox := text.New(pixel.V(0, 0), atlas)
	textbox.Color = colornames.Grey
	fmt.Fprintf(textbox, "Hello\n")
	return textbox

}
func textFrame(win *pixelgl.Window) {
}
