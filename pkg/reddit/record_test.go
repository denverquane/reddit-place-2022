package reddit

import (
	"image"
	"image/color"
	"testing"
)

const (
	RedditExampleTimeUnixMilli = 1649007502252
)

func TestParseTime(t *testing.T) {
	i := "2022-04-03 17:38:22.252 UTC"
	ti, err := ParseTime(i)
	if err != nil {
		t.Error(err)
	}
	if ti.UnixMilli() != RedditExampleTimeUnixMilli {
		t.Error()
	}
}

func TestToRecord(t *testing.T) {
	line := []string{
		"2022-04-03 17:38:22.252 UTC",
		"yTrYCd4LUpBn4rIyNXkkW2+Fac5cQHK2lsDpNghkq0oPu9o//8oPZPlLM4CXQeEIId7l011MbHcAaLyqfhSRoA==",
		"#FF3881",
		"0,0",
	}
	r, err := ToRecord(line)
	if err != nil {
		t.Error(err)
	}
	if r.Time.UnixMilli() != RedditExampleTimeUnixMilli {
		t.Error("Time not parsed correctly")
	}
	if r.userID != "yTrYCd4LUpBn4rIyNXkkW2+Fac5cQHK2lsDpNghkq0oPu9o//8oPZPlLM4CXQeEIId7l011MbHcAaLyqfhSRoA==" {
		t.Error("UserID not parsed correctly")
	}
	p := color.RGBA{
		R: 255,
		G: 56,
		B: 129,
		A: 255,
	}
	if r.color != p {
		t.Error("Color not parsed correctly")
	}
	if r.rect != image.Rect(0, 0, 0, 0) {
		t.Error("Rect not parsed correctly")
	}

	line[3] = "1,2,3,4"
	r, err = ToRecord(line)
	if err != nil {
		t.Error(err)
	}
	if r.rect != image.Rect(1, 2, 3, 4) {
		t.Error("Moderator rect not parsed correctly")
	}
}
