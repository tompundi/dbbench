package testdb

import (
	"errors"
)

var errMemoryNotAllowed = errors.New(":memory: path not available")

type Options struct {
	Port  int
	Which string
	Fsync bool
	Path  string
}

type Store interface {
	Close() error
	Set(key, value []byte) error
	PSet(keys, values [][]byte) error
	Get(key []byte) ([]byte, bool, error)
	PGet(keys [][]byte) ([][]byte, []bool, error)
	Del(key []byte) (bool, error)
	Keys(pattern []byte, limit int, withvalues bool) ([][]byte, [][]byte, error)
	FlushDB() error
}

func bcopy(b []byte) []byte {
	r := make([]byte, len(b))
	copy(r, b)
	return r
}
