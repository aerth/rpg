package rpg

import (
	"fmt"
	"unicode"

	"github.com/aerth/rpg/assets"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func getrunes() []rune {
	RUNES := make([]rune, unicode.MaxRune-32)
	for i := range RUNES {
		RUNES[i] = rune(32 + i)
	}
	return RUNES
}

var GraphicRanges = []*unicode.RangeTable{
	unicode.L, unicode.M, unicode.N, unicode.P, unicode.S, unicode.Zs,
}

func NewTextSmooth(size float64) *text.Text {
	font := ttfFromBytesMust(goregular.TTF, size)
	basicAtlas := text.NewAtlas(font, text.ASCII, text.RangeTable(unicode.Common))
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)
	return basicTxt
}

func NewText(size float64) *text.Text {
	//	font := ttfFromBytesMust(goregular.TTF, size)
	b, err := assets.Asset("font/admtas.ttf")
	if err != nil {
		panic(err)
	}
	font := ttfFromBytesMust(b, size)
	// font := basicfont.Face7x13
	basicAtlas := text.NewAtlas(font, text.ASCII, text.RangeTable(unicode.Common))
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)
	return basicTxt
}

func DrawText(winbounds pixel.Rect, t *text.Text, canvas pixel.Target, format string, i ...interface{}) {
	imd := imdraw.New(nil)
	color := pixel.ToRGBA(colornames.Darkslategrey)
	imd.Color = color.Scaled(0.5)
	imd.Push(pixel.V(0, winbounds.Max.Y-50), pixel.V(winbounds.Max.XY()))
	imd.Rectangle(0)
	imd.Push(pixel.V(0, 0), pixel.V(winbounds.Max.X, 80))
	imd.Rectangle(0)
	imd.Draw(canvas)
	t.Dot = pixel.V(10, 10)
	t.Orig = pixel.V(10, 10)
	fmt.Fprintf(t, format, i...)
	t.Draw(canvas, pixel.IM)
}

func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}
func DrawScore(winbounds pixel.Rect, t *text.Text, canvas pixel.Target, format string, i ...interface{}) {
	imd := imdraw.New(nil)
	color := pixel.ToRGBA(colornames.Darkslategrey)
	imd.Color = color.Scaled(0.5)
	imd.Push(pixel.V(0, winbounds.Max.Y-50), pixel.V(winbounds.Max.XY()))
	imd.Rectangle(0)
	imd.Push(pixel.V(0, 0), pixel.V(winbounds.Max.X, 80))
	imd.Rectangle(0)
	imd.Draw(canvas)
	t.Dot = pixel.V(10, winbounds.Max.Y-40)
	t.Orig = t.Dot
	fmt.Fprintf(t, format, i...)
	t.Draw(canvas, pixel.IM)
}

func (w *World) Message(s string) {
	w.Messages = append(w.Messages, s)
}
