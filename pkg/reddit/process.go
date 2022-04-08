package reddit

import (
	"image"
	"image/png"
	"log"
	"os"
)

func MakeImage(filename string, size image.Rectangle, pixels []*PixelEdit) {
	img := image.NewRGBA(size)

	for _, r := range pixels {
		c, err := ParseHexColor(r.HexColor)
		if err != nil {
			log.Println(err)
		}
		img.Set(int(r.X), int(r.Y), c)
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	if err = png.Encode(f, img); err != nil {
		log.Println(err)
	}
	f.Close()
}
