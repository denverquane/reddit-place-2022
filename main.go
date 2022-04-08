package main

import (
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/denverquane/reddit-place-2022/pkg/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	Cleanup = false
)

func main() {
	worker := storage.PostgresWorker{}
	err := worker.Init("internal/postgres.sql", os.Getenv("POSTGRES_URL"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	filenames := file.GenerateFileNames()
	workQueue := make(chan reddit.Record)

	// TODO parallelize across more files? Probably an I/O bound application, but probably not across separate files
	go func() {
		for _, name := range filenames {
			if !file.DirectoryContains(".", name+".csv") {
				if !file.DirectoryContains(".", name+".csv.gzip") {
					log.Printf("Missing %s.csv.gzip, downloading now\n", name)
					file.DownloadGzip(name+".csv", file.DataBaseURL+name+".csv.gzip")
					file.Parse(name+".csv", workQueue)
				}
			}
		}
	}()

	go worker.Start(workQueue)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	if Cleanup {
		log.Println("Cleanup")
		for _, name := range filenames {
			os.Remove(name + ".csv")
			os.Remove(name + ".csv.gzip")
		}
	}
}
