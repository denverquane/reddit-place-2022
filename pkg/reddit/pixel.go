package reddit

import "image/color"

const (
	// See https://deephaven.io/core/docs/reference/query-language/types/nulls/
	NullShort int32 = -32768
)

type Edit struct {
	// unused fields right now, faster to unmarshal without them specified
	// Timestamp int64 `parquet:"name=timestamp, type=INT64, convertedtype=TIMESTAMP"`
	// UserID    int   `parquet:"name=user_id, type=INT64"`
	Color int32 `parquet:"name=rgb, type=INT32"`
	X1    int32 `parquet:"name=x1, type=INT32, convertedtype=INT_16"`
	Y1    int32 `parquet:"name=y1, type=INT32, convertedtype=INT_16"`
	X2    int32 `parquet:"name=x2, type=INT32, convertedtype=INT_16"`
	Y2    int32 `parquet:"name=y2, type=INT32, convertedtype=INT_16"`
}

func (e Edit) GetColor() color.Color {
	return color.RGBA{
		R: uint8((e.Color >> 16) & 255),
		G: uint8((e.Color >> 8) & 255),
		B: uint8(e.Color & 255),
		A: 255,
	}
}

func (e Edit) IsMod() bool {
	// only if the X2/Y2 are set, and actually differ from x1/y1
	return (e.X2 != NullShort || e.Y2 != NullShort) && (e.X2 != e.X1 || e.Y2 != e.Y1)
}
