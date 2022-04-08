package main

import (
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
	"github.com/denverquane/reddit-place-2022/pkg/storage"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	worker := storage.PostgresWorker{}
	// TODO can't end in /, fix and use proper paths
	postgresSchemaDir := "internal"
	if os.Getenv("POSTGRES_SCHEMA_DIR") != "" {
		postgresSchemaDir = os.Getenv("POSTGRES_SCHEMA_DIR")
	}
	// TODO can't end in /, fix and use proper paths
	dataDir := "data"
	if os.Getenv("PLACE_DATA_DIR") != "" {
		dataDir = os.Getenv("PLACE_DATA_DIR")
	}
	err := worker.Init(postgresSchemaDir+"/postgres.sql", os.Getenv("POSTGRES_URL"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	// TODO configure download vs postgres vs image generation
	if os.Getenv("ONLY_DRAW") == "" {
		filenames := file.GenerateFileNames()
		fileQueue := make(chan string, len(filenames))

		//one background task for downloading files (probably won't benefit from parallelism)
		go func() {
			for _, name := range filenames {
				if !file.DirectoryContains(dataDir, name+".csv") && !file.DirectoryContains(dataDir, name+".csv.complete") {
					log.Printf("Missing %s/%s.csv, downloading now\n", dataDir, name)
					err := file.DownloadGzip(dataDir+"/"+name+".csv", file.DataBaseURL+name+".csv.gzip")
					if err != nil {
						log.Println(err)
						continue
					}
				}
				// regardless of if we download or not, only enqueue files that aren't marked as complete
				if !file.DirectoryContains(dataDir, name+".csv.complete") {
					fileQueue <- dataDir + "/" + name + ".csv"
				}
			}
			log.Println("Download worker is done")
			close(fileQueue)
		}()

		go fileWorker(fileQueue, &worker)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	} else {
		drawImage(&worker)
	}

	worker.Close()
}

func drawImage(worker *storage.PostgresWorker) {
	start := time.Now()
	px, err := worker.GetPixelsUpToTimeInRegion(
		time.Date(2022, time.April, 3, 0, 0, 0, 0, time.UTC),
		image.Rect(0, 0, 500, 500))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Took", time.Since(start), "to fetch data from Postgres")
	start = time.Now()
	// TODO the pixels drawn might be offset if we requested an offset region from Postgres
	reddit.MakeImage("place.png", image.Rect(0, 0, 500, 500), px)
	log.Println("Took", time.Since(start), "to draw place.png")
}

func fileWorker(fileQueue <-chan string, worker *storage.PostgresWorker) {
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
		f, err := os.Create(name + ".complete")
		if err != nil {
			log.Println(err)
		}
		f.Close()
		log.Printf("Worker finished %s\n", name)
	}
	log.Printf("File worker is done\n")
}
