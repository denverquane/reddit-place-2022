package file

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
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

var snapshotRegex = regexp.MustCompile(`^checkpoint_(?P<timestamp>[0-9]+)_(?P<row>[0-9]+).png$`)

func GetCheckpointsFromDirectory(directory string) map[int64]int64 {
	timestampsToRows := make(map[int64]int64)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println(err)
		return timestampsToRows
	}
	for _, file := range files {
		match := snapshotRegex.FindStringSubmatch(file.Name())
		if len(match) == 3 {
			time, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				continue
			}
			row, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				continue
			}
			timestampsToRows[time] = row
		}
	}
	return timestampsToRows
}
