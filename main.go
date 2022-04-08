package main

import (
	"github.com/denverquane/reddit-place-2022/pkg/file"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	Cleanup = false
)

func main() {
	filenames := file.GenerateFileNames()
	go func() {
		for _, name := range filenames {
			if !file.DirectoryContains(".", name+".csv") {
				if !file.DirectoryContains(".", name+".csv.gzip") {
					log.Printf("Missing %s.csv.gzip, downloading now\n", name)
					file.DownloadGzip(name+".csv", file.DataBaseURL+name+".csv.gzip")
					file.Parse(name + ".csv")
				}
			}
		}
	}()

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
