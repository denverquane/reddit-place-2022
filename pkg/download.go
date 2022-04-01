package pkg

import (
	"errors"
	"image"
	"image/png"
	"net/http"
)

func DownloadImage(URL string) (image.Image, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}

	return png.Decode(response.Body)
}
