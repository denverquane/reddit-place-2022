package reddit

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

const (
	PercentSnapshot float64 = 1
)

func ProcessWorker(dataQueue <-chan Record, totalLines float64) {
	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))

	var line float64
	var percentThreshold = PercentSnapshot
	for r := range dataQueue {
		img.Set(r.rect.Min.X, r.rect.Min.Y, r.color)
		line++
		if line/totalLines*100.0 > percentThreshold {
			filename := fmt.Sprintf("images/place_%d.png", int(percentThreshold))
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
			percentThreshold += PercentSnapshot
		}
	}
	log.Println("Finished processing data")
}
