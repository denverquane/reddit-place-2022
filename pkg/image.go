package pkg

import (
	"image"
	"image/draw"
)

func CombineDiffToBase(base, diff image.Image) (image.Image, error) {
	draw.Draw(base.(*image.Paletted), base.Bounds(), diff, image.Point{}, draw.Over)
	return base, nil
}
