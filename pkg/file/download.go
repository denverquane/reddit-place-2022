package file

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	TotalFiles   = 79
	DataBaseURL  = "https://placedata.reddit.com/data/canvas-history/"
	DataFileBase = "2022_place_canvas_history-0000000000" // removed two 0s
)

func GenerateFileNames() [TotalFiles]string {
	var arr [TotalFiles]string
	for i := 0; i < TotalFiles; i++ {
		if i < 10 {
			arr[i] = fmt.Sprintf("%s0%d", DataFileBase, i)
		} else {
			arr[i] = fmt.Sprintf("%s%d", DataFileBase, i)
		}
	}
	return arr
}

func DownloadGzip(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	start := time.Now()

	client := new(http.Client)
	request, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	reader, err = gzip.NewReader(resp.Body)

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)

	return nil
}
