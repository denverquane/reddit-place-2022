package reddit

import (
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Time   time.Time
	userID string
	color  color.RGBA
	rect   image.Rectangle
}

func ToRecord(line []string) (Record, error) {
	t, err := ParseTime(line[0])
	if err != nil {
		return Record{}, err
	}
	coords := strings.Split(line[3], ",")

	x, err := strconv.ParseUint(coords[0], 10, 64)
	if err != nil {
		return Record{}, err
	}
	y, err := strconv.ParseUint(coords[1], 10, 64)
	if err != nil {
		return Record{}, err
	}
	rect := image.Rect(int(x), int(y), int(x), int(y))

	// exceptional overwrite case for moderator actions
	if len(coords) == 4 {
		x1, err := strconv.ParseInt(coords[2], 10, 16)
		if err != nil {
			return Record{}, err
		}
		y1, err := strconv.ParseInt(coords[3], 10, 16)
		if err != nil {
			return Record{}, err
		}
		rect.Max = image.Point{
			X: int(x1),
			Y: int(y1),
		}
	}
	c, err := ParseHexColor(line[2])
	if err != nil {
		return Record{}, err
	}
	return Record{
		Time:   t,
		userID: line[1],
		color:  c,
		rect:   rect,
	}, nil
}

func ParseTime(input string) (time.Time, error) {
	return time.Parse(`2006-01-02 15:04:05.999 MST`, input)
}
