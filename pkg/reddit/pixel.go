package reddit

type PixelEdit struct {
	HexColor string `db:"pixel_color"`
	X        uint16 `db:"x"`
	Y        uint16 `db:"y"`
}
