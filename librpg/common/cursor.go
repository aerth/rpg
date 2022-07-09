package common

import (
	"strconv"

	"github.com/faiface/pixel"
)

func GetCursor(num int) *pixel.Sprite {
	pic, err := LoadPicture("sprites/cursor" + strconv.Itoa(num) + ".png")
	if err != nil {

		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return sprite
}
