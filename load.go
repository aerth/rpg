// aerth game
// copyright 2017 aerth <aerth@riseup.net>
package rpg

import (
	"bytes"
	"image"
	"math/rand"
	"time"

	_ "image/png"

	"github.com/aerth/rpg/assets"
	"github.com/faiface/pixel"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// loadPicture from assets
func LoadPicture(path string) (pixel.Picture, error) {
	b, err := assets.Asset(path)
	if err != nil {
		return nil, err
	}
	file := bytes.NewReader(b)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// loadCharacterSheet returns an animated spritesheet
func LoadCharacterSheet(sheetPath string, numframes uint8) (sheet pixel.Picture, anims map[Direction][]pixel.Rect, err error) {
	sheet, err = LoadPicture("sprites/char.png")
	frameWidth := int(sheet.Bounds().Max.X/float64(numframes)) * 2
	//log.Println(frameWidth, "width")
	// create a array of frames inside the spritesheet
	var frames = new([16]pixel.Rect)

	for i, x := 0, 0.0; x+float64(frameWidth) <= sheet.Bounds().Max.X; i, x = i+1, x+float64(frameWidth) {
		if i > 15 {
			break
		}
		frames[i] = pixel.R(
			x,
			0,
			x+float64(frameWidth),
			sheet.Bounds().H(),
		)
	}

	anims = make(map[Direction][]pixel.Rect)
	anims[LEFT] = frames[:4]
	anims[RIGHT] = frames[4:8]
	anims[DOWN] = frames[8:12]
	anims[UP] = frames[12:]
	return sheet, anims, nil
}
