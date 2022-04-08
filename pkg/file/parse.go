package file

import (
	"encoding/csv"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"io"
	"log"
	"os"
)

func Parse(filename string, workQueue chan<- reddit.Record) {
	log.Println("Parsing", filename)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	// skip the header
	r.Read()
	for {
		tokens, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		recs, err := reddit.ToRecords(tokens)
		if err != nil {
			log.Println(err)
		} else {
			for _, rec := range recs {
				workQueue <- rec
			}
		}
	}
	log.Println(filename, "complete")
}
