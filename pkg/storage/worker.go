package storage

import (
	"errors"
	"github.com/denverquane/reddit-place-2022/pkg/reddit"
)

var (
	noInitError = errors.New("worker has not been initialized")
)

type Worker interface {
	Init(filepath, url, user, pass string) error
	Add(record reddit.Record) error
	Start(<-chan reddit.Record) error
}
