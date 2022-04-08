package storage

import (
	"context"
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"os"
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

func (w *PostgresWorker) GetLastEditForEveryPixel() ([]*reddit.PixelEdit, error) {
	var pixels []*reddit.PixelEdit
	err := pgxscan.Select(context.Background(), w.Conn, &pixels, "select pixel_color, x, y from (select distinct on (x, y) timestamp, pixel_color, x, y from events) as event order by timestamp desc;")
	if err != nil {
		return pixels, err
	}
	return pixels, nil
}
