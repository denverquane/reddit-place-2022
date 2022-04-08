package reddit

import (
	"image"
	"time"
)

type Record struct {
	Time   time.Time
	UserID string
	Color  string
	Pixel  image.Point
}
