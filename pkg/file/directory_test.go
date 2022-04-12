package file

import "testing"

const (
	TestdataTimestamp = 1648820056223
	TestdataRow       = 223408
)

func TestDirectoryContains(t *testing.T) {
	if !DirectoryContains("testdata", "checkpoint_1648820056223_223408.png") {
		t.Error("Expected testdata to contain checkpoint_1648820056223_223408.png")
	}
}

func TestGetSnapshotsFromDirectory(t *testing.T) {
	c := GetCheckpointsFromDirectory("testdata")
	if len(c) != 1 {
		t.Error("Expected 1 checkpoint element from testdata directory")
	}

	var row int64
	var ok bool
	if row, ok = c[TestdataTimestamp]; !ok {
		t.Error("No element found for timestamp", TestdataTimestamp)
	}
	if row != TestdataRow {
		t.Error("Element at timestamp", TestdataTimestamp, "wasn't", TestdataRow)
	}
}
