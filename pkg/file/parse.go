package file

import (
	"bufio"
	"context"
	"github.com/denverquane/reddit-place-2022/pkg/storage"
	"log"
	"os"
	"strings"
)

func ParseAndAdd(filename string, worker storage.PostgresWorker) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	newFile, err := os.Create(filename + "_tmp")
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(f)
	writer := bufio.NewWriter(newFile)
	if sc.Scan() {
		// header
		_, err := writer.WriteString("timestamp,user_id,pixel_color,x,y,x1,y1\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	for sc.Scan() {
		line := strings.ReplaceAll(sc.Text(), "\"", "")
		if strings.Count(line, ",") == 4 {
			line += ",,"
		}
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
	newFile.Close()
	newFile, err = os.Open(filename + "_tmp")
	if err != nil {
		log.Fatal(err)
	}

	// use a temp table to prevent errors when adding duplicate primary keys (thanks Reddit)
	_, err = worker.Conn.PgConn().CopyFrom(context.Background(), newFile, "COPY tmp_table FROM STDIN CSV HEADER")
	if err != nil {
		newFile.Close()
		return err
	}
	r := worker.Conn.PgConn().Exec(context.Background(), "INSERT INTO events SELECT * FROM tmp_table ON CONFLICT DO NOTHING")
	res, err := r.ReadAll()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res)
	}

	newFile.Close()
	return os.Remove(filename + "_tmp")
}
