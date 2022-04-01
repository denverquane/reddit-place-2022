package pkg

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"sync"
)

func HandleRequestWrapper(image *image.Image, imageLock *sync.RWMutex) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		imageLock.RLock()
		err := png.Encode(buf, *image)
		imageLock.RUnlock()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(buf.Bytes())
		return
	}
}
