package reddit

import (
	"image"
	"image/png"
	"log"
	"os"
)

func DrawSubregionToFile(img *image.RGBA, region image.Rectangle, filename string) error {
	log.Println("Writing to", filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	sub := img.SubImage(region)
	if err = png.Encode(f, sub); err != nil {
		return err
	}
	return f.Close()
}

func DrawToFile(img *image.RGBA, filename string) error {
	log.Println("Writing to", filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err = png.Encode(f, img); err != nil {
		return err
	}
	return f.Close()
}
