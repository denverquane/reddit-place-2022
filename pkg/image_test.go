package pkg

import (
	"image/png"
	"log"
	"os"
	"testing"
)

func TestCombineDiffToBase(t *testing.T) {
	baseFile, err := os.Open("testdata/place.png")
	if err != nil {
		t.Fatal(err)
	}
	defer baseFile.Close()
	diffFile, err := os.Open("testdata/diff.png")
	if err != nil {
		t.Fatal(err)
	}
	defer diffFile.Close()

	baseImg, err := png.Decode(baseFile)
	if err != nil {
		t.Error(err)
	}
	diffImg, err := png.Decode(diffFile)
	if err != nil {
		t.Error(err)
	}

	out, err := CombineDiffToBase(baseImg, diffImg)
	if err != nil {
		t.Fatal(err)
	}

	for y := 0; y < 1000; y++ {
		for x := 0; x < 1000; x++ {
			_, _, _, a := diffImg.At(x, y).RGBA()
			if a == 65535 {
				if out.At(x, y) != diffImg.At(x, y) {
					t.Error("pixels differ between outfile and diff file at x,y: ", x, y)
					log.Println(out.At(x, y))
					log.Println(diffImg.At(x, y))
				}
			}
		}
	}
}
