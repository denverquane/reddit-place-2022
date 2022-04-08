package file

import (
	"io/ioutil"
	"log"
)

func DirectoryContains(directory, filename string) bool {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println(err)
		return false
	}
	for _, v := range files {
		if v.Name() == filename {
			return true
		}
	}
	return false
}
