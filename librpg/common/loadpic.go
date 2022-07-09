package common

import (
	"bytes"
	"image"

	"github.com/aerth/rpc/assets"
	"github.com/faiface/pixel"
)

// loadPicture from embedded assets
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
