package file

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func DownloadFile(url, filepath string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Println("Downloading", url, "to", filepath)
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)

	return nil
}
