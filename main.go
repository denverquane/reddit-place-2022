package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/gin-gonic/gin"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"image"
	"image/color/palette"
	"image/gif"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ParquetFileName = "2022_place_deephaven.parquet"
	ParquetBaseURL  = "https://deephaven.io/wp-content/"
	RowBufferSize   = 10000
	FirstTimeMillis = 1648817050315
	LastTimeMillis  = 1649117640195

	// The total length of r/place 2022, as indicated by this dataset
	PlaceTotalMillisElapsed = 300_589_880
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

	r := gin.Default()
	r.GET("/gif", func(c *gin.Context) {
		rect, err := getRectFromQuery(c.Query("rect"))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
			return
		}
		correctedRect := image.Rect(0, 0, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
		// TODO determine better palette?
		frame := image.NewPaletted(correctedRect, palette.WebSafe)
		pr, err := reader.NewParquetReader(fr, new(reddit.Edit), 4)
		if err != nil {
			log.Println("Can't create parquet reader", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}
		outGif := &gif.GIF{}
		totalRows := pr.GetNumRows()
		var row int64
		var prevTime int64
		for row < totalRows {
			pixels := make([]reddit.Edit, RowBufferSize)
			if err = pr.Read(&pixels); err != nil {
				log.Println("Read error", err)
				c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
				return
			} else {
				for _, pix := range pixels {
					if pix.Overlaps(*rect) {
						// use the provided rect as an offset for the smaller image
						frame.Set(int(pix.X1)-rect.Min.X, int(pix.Y1)-rect.Min.Y, pix.GetColor())

						// capture a gif frame every 60 seconds, but only if a pixel changed in our rectangle bounds
						// divide by 1 mil b/c every datetime seems to be that inflated in the raw data...
						if (pix.Timestamp-prevTime)/1_000_000 > (time.Second.Milliseconds() * 60) {
							clone := *frame
							clone.Pix = make([]uint8, len(frame.Pix))
							copy(clone.Pix, frame.Pix)

							outGif.Image = append(outGif.Image, &clone)
							outGif.Delay = append(outGif.Delay, 2)
							prevTime = pix.Timestamp
						}
					}
					row++
				}
			}
		}
		pr.ReadStop()
		buf := bytes.NewBuffer([]byte{})
		err = gif.EncodeAll(buf, outGif)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}
		c.Data(http.StatusOK, "image/gif", buf.Bytes())
	})
	r.Run()
	fr.Close()
}

func getRectFromQuery(rectString string) (*image.Rectangle, error) {
	tokens := strings.Split(rectString, ",")
	if len(tokens) < 4 {
		return nil, errors.New("invalid rect string supplied, expected of the form \"rect=<x1>,<y1>,<x2>,<y2>\"")
	}
	x1, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return nil, err
	}
	y1, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return nil, err
	}
	x2, err := strconv.ParseUint(tokens[2], 10, 64)
	if err != nil {
		return nil, err
	}
	y2, err := strconv.ParseUint(tokens[3], 10, 64)
	if err != nil {
		return nil, err
	}
	// invalid negatives handled by the parseuint above
	if x1 > 2000 || y1 > 2000 || x2 > 2000 || y2 > 2000 {
		return nil, errors.New("expected parameters in the range [0-2000]")
	}
	if x1 > x2 || y1 > y2 {
		return nil, errors.New("x2 and y2 must be strictly greater than x1 and y2")
	}
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	return &r, nil
}
