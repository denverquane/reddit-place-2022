package main

import (
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"image"
	"log"
	"os"
	"time"
)

const (
	ParquetFileName  = "2022_place_deephaven.parquet"
	ParquetBaseURL   = "https://deephaven.io/wp-content/"
	RowBufferSize    = 10000
	DrawPreModImages = true
	// set to 0 to disable
	DrawEveryVarPercent = 5
)

func main() {
	// TODO can't end in /, fix and use proper paths
	dataDir := "data"
	if os.Getenv("PLACE_DATA_DIR") != "" {
		dataDir = os.Getenv("PLACE_DATA_DIR")
	}
	if !file.DirectoryContains(dataDir, ParquetFileName) {
		log.Println("Looks like you're missing", dataDir+"/"+ParquetFileName)
		err := file.DownloadFile(ParquetBaseURL+ParquetFileName, dataDir+"/"+ParquetFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	fr, err := local.NewLocalFileReader(dataDir + "/" + ParquetFileName)
	if err != nil {
		log.Println("Can't open file")
		return
	}
	pr, err := reader.NewParquetReader(fr, new(reddit.Edit), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	var row int64 = 0
	totalRows := pr.GetNumRows()
	percentThreshold := DrawEveryVarPercent
	start := time.Now()
	for row < totalRows {
		pixels := make([]reddit.Edit, RowBufferSize)
		if err = pr.Read(&pixels); err != nil {
			log.Println("Read error", err)
		}
		for _, r := range pixels {
			if DrawPreModImages && r.IsMod() {
				err := reddit.DrawSubregionToFile(img, image.Rect(int(r.X1), int(r.Y1), int(r.X2), int(r.Y2)),
					fmt.Sprintf("images/place_mod_[%d,%d]-[%d,%d].png", r.X1, r.Y1, r.Y1, r.Y2))
				if err != nil {
					log.Println(err)
				}
			} else if r.X1 < 0 || r.Y1 < 0 || r.X1 > 2000 || r.Y1 > 2000 {
				// TODO determine these pixels' existence. Looks like an unmarshalling error from Parquet(?)
			}
			img.Set(int(r.X1), int(r.Y1), r.GetColor())
			row++
			if percentThreshold > 0 && (float64(row)/float64(totalRows)*100.0) > float64(percentThreshold) {
				err := reddit.DrawToFile(img, fmt.Sprintf("images/place_%d.png", percentThreshold))
				if err != nil {
					log.Println(err)
				}
				percentThreshold += DrawEveryVarPercent
			}
		}
	}
	log.Println("Took", time.Since(start), "to generate final image")
	err = reddit.DrawToFile(img, "images/place.png")
	if err != nil {
		log.Println(err)
	}
	pr.ReadStop()
	fr.Close()
}
