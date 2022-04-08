package reddit

import (
	"image"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Time   time.Time
	UserID string
	Color  string
	Pixel  image.Point
}

func ToRecords(line []string) ([]Record, error) {
	t, err := ParseTime(line[0])
	if err != nil {
		return nil, err
	}

	coords := strings.Split(line[3], ",")
	x1, err := strconv.ParseInt(coords[0], 10, 64)
	if err != nil {
		return nil, err
	}
	y1, err := strconv.ParseInt(coords[1], 10, 64)
	if err != nil {
		return nil, err
	}
	if len(coords) == 2 {
		return []Record{
			{
				Time:   t,
				UserID: line[1],
				Color:  line[2],
				Pixel:  image.Point{X: int(x1), Y: int(y1)},
			},
		}, nil
	}
	x2, err := strconv.ParseInt(coords[2], 10, 16)
	if err != nil {
		return nil, err
	}
	y2, err := strconv.ParseInt(coords[3], 10, 16)
	if err != nil {
		return nil, err
	}

	records := make([]Record, (x2-x1)*(y2-y1))
	i := 0
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			records[i] = Record{
				Time:   t,
				UserID: line[1],
				Color:  line[2],
				Pixel:  image.Point{X: int(x), Y: int(y)},
			}
		}
	}

	return records, nil
}

func ParseTime(input string) (time.Time, error) {
	return time.Parse(`2006-01-02 15:04:05.999 MST`, input)
}
