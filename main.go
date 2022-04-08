package main

import (
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"github.com/denverquane/reddit-place-2022/pkg/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO configure download vs postgres vs image generation
func main() {
	worker := storage.PostgresWorker{}
	err := worker.Init("internal/postgres.sql", os.Getenv("POSTGRES_URL"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	filenames := file.GenerateFileNames()
	fileQueue := make(chan string, len(filenames))

	// one background task for downloading files (probably won't benefit from parallelism)
	go func() {
		for _, name := range filenames {
			if !file.DirectoryContains("data", name+".csv") {
				log.Printf("Missing data/%s.csv, downloading now\n", name)
				err := file.DownloadGzip("data/"+name+".csv", file.DataBaseURL+name+".csv.gzip")
				if err != nil {
					log.Println(err)
					continue
				}
			}
			fileQueue <- "data/" + name + ".csv"
		}
		log.Println("Download worker is done")
		close(fileQueue)
	}()

	go fileWorker(fileQueue, worker)

	//px, err := worker.GetLastEditForEveryPixel()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//reddit.MakeImage(px)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// TODO more gracefully determine what files were completely processed
}

func fileWorker(fileQueue <-chan string, worker storage.PostgresWorker) {
	for name := range fileQueue {
		log.Printf("Worker picked up %s\n", name)
		err := file.ParseAndAdd(name, worker)
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(name)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Worker finished %s\n", name)
	}
	log.Printf("File worker is done\n")
}
