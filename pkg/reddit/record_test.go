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

	// test the variadic number of digits after the decimal
	i = "2022-04-03 17:38:22.25 UTC"
	ti, err = ParseTime(i)
	if err != nil {
		t.Error(err)
	}
	if ti.UnixMilli() != RedditExampleTimeUnixMilli-2 {
		t.Error()
	}

	i = "2022-04-03 17:38:22.2 UTC"
	ti, err = ParseTime(i)
	if err != nil {
		t.Error(err)
	}
	if ti.UnixMilli() != RedditExampleTimeUnixMilli-52 {
		t.Error()
	}

	i = "2022-04-03 17:38:22 UTC"
	ti, err = ParseTime(i)
	if err != nil {
		t.Error(err)
	}
	if ti.UnixMilli() != RedditExampleTimeUnixMilli-252 {
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
	r, err := ToRecords(line)
	if err != nil {
		t.Error(err)
	}
	if r[0].Time.UnixMilli() != RedditExampleTimeUnixMilli {
		t.Error("Time not parsed correctly")
	}
	if r[0].UserID != "yTrYCd4LUpBn4rIyNXkkW2+Fac5cQHK2lsDpNghkq0oPu9o//8oPZPlLM4CXQeEIId7l011MbHcAaLyqfhSRoA==" {
		t.Error("UserID not parsed correctly")
	}
	p := color.RGBA{
		R: 255,
		G: 56,
		B: 129,
		A: 255,
	}
	if r[0].Color != p {
		t.Error("Color not parsed correctly")
	}
	if r[0].Pixel != image.Pt(0, 0) {
		t.Error("Rect not parsed correctly")
	}

	line[3] = "1,2,3,4"
	r, err = ToRecords(line)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 4 {
		t.Error("Did not generate 4 expected pixels from a moderator square")
	}
}
