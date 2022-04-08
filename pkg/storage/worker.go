package storage

type Worker interface {
	Init(filepath, url, user, pass string) error
	Close() error
}
