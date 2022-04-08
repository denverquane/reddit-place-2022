package storage

import (
	"context"
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"os"
)

type PostgresWorker struct {
	pool *pgxpool.Pool
}

func (w *PostgresWorker) Init(filepath, url, user, pass string) error {
	log.Println("Connecting to Postgres at", url)
	connString := fmt.Sprintf("postgres://%s?user=%s&password=%s", url, user, pass)
	dbpool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return err
	} else {
		log.Println("Connection successful")
	}
	w.pool = dbpool

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tag, err := w.pool.Exec(context.Background(), string(bytes))
	if err != nil {
		return err
	}
	log.Println(tag.String())
	return nil
}

func (w *PostgresWorker) Add(record reddit.Record) error {
	if w.pool == nil {
		return noInitError
	}
	_, err := w.pool.Exec(context.Background(), "INSERT INTO events VALUES ($1, $2, $3, $4, $5);", record.Time, record.UserID, record.Color, record.Pixel.X, record.Pixel.Y)
	return err
}

func (w *PostgresWorker) Start(tasks <-chan reddit.Record) error {
	if w.pool == nil {
		return noInitError
	}

	for record := range tasks {
		err := w.Add(record)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
