package rpg

import (
	"log"
	"testing"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func init() {
	log.SetFlags(0)
}
func TestAlpha(t *testing.T) {
	color := pixel.ToRGBA(colornames.Red)
	log.Println(color)
	color = color.Mul(pixel.Alpha(0.2))
	log.Println(color)
	color = color.Scaled((0.2))
	log.Println(color)
	color = color.Mul(pixel.Alpha(0))
	log.Println(color)
}
