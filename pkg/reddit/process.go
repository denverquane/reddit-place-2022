package reddit

import (
	"image"
	"image/png"
	"log"
	"os"
)

func MakeImage(pixels []*PixelEdit) {
	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))

	for _, r := range pixels {
		c, err := ParseHexColor(r.HexColor)
		if err != nil {
			log.Println(err)
		}
		img.Set(int(r.X), int(r.Y), c)
	}
	filename := "place.png"
	log.Println("Writing to", filename)
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	if err = png.Encode(f, img); err != nil {
		log.Println(err)
	}
	f.Close()
	log.Println("Finished writing img data")
}
