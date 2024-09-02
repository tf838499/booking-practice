package testdata

import (
	"path/filepath"
	"runtime"
)

var basepath string

const (
	TestDataConsumer = "consumers.yml"
	TestDataEvent    = "events.yml"
	TestDataSeat     = "seats.yml"
	TestDataTicket   = "tickets.yml"
)

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

func Path(rel string) string {
	return filepath.Join(basepath, rel)
}
