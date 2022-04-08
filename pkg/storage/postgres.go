package storage

import (
	"context"
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type PostgresWorker struct {
	Conn *pgx.Conn
}

func (w *PostgresWorker) Init(filepath, url, user, pass string) error {
	log.Println("Connecting to Postgres at", url)
	connString := fmt.Sprintf("postgres://%s?user=%s&password=%s", url, user, pass)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	} else {
		log.Println("Connection successful")
	}
	w.Conn = conn

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tag, err := w.Conn.Exec(context.Background(), string(bytes))
	if err != nil {
		return err
	}
	log.Println(tag.String())
	return nil
}

func (w *PostgresWorker) Close() error {
	if w.Conn != nil {
		return w.Conn.Close(context.Background())
	}
	return nil
}

func (w *PostgresWorker) GetPixelsUpToTime(time time.Time) ([]*reddit.PixelEdit, error) {
	return w.GetPixelsUpToTimeInRegion(time, image.Rect(0, 0, 2000, 2000))
}

func (w *PostgresWorker) GetPixelsUpToTimeInRegion(time time.Time, region image.Rectangle) ([]*reddit.PixelEdit, error) {
	var pixels []*reddit.PixelEdit
	err := pgxscan.Select(context.Background(), w.Conn, &pixels, "select pixel_color, x, y from events "+
		"where timestamp < $1 and x > $2 and y > $3 and x < $4 and y < $5",
		time, region.Min.X, region.Min.Y, region.Max.X, region.Max.Y)
	if err != nil {
		return pixels, err
	}
	return pixels, nil
}

// TODO this might be broken, hard to say without having a db with 100% of the data yet
func (w *PostgresWorker) GetLastEditForPixelsInRegion(region image.Rectangle) ([]*reddit.PixelEdit, error) {
	var pixels []*reddit.PixelEdit
	err := pgxscan.Select(context.Background(), w.Conn, &pixels, "select pixel_color, x, y from "+
		"(select distinct on (x, y) timestamp, pixel_color, x, y from events where x > $1 and y > $2 and x < $3 and y < $4) "+
		"as event order by timestamp desc;", region.Min.X, region.Min.Y, region.Max.X, region.Max.Y)
	if err != nil {
		return pixels, err
	}
	return pixels, nil
}

func (w *PostgresWorker) GetLastEditForEveryPixel() ([]*reddit.PixelEdit, error) {
	return w.GetLastEditForPixelsInRegion(image.Rect(0, 0, 2000, 2000))
}
