package api

import (
	"encoding/json"
	"strconv"
	"time"
)

// if two modules will share the tpyes we may use this module to share module apis

// Entry is a log entry with a timestamp.
type Entry struct {
	TimeStamp          time.Time         `json:"ts"`
	Line               string            `json:"line"`
	StructuredMetadata map[string]string `json:"meta"`
}

func (e Entry) list() []string {
	a := make([]string, 0, 2)
	a = append(a, strconv.FormatInt(e.TimeStamp.UnixNano(), 10), e.Line)
	return a
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.list())
}

// EntryHandler is something that can "handle" entries via channel.
// Stop must be called to gracefully shutdown the EntryHandler
type EntryHandler interface {
	Chan() chan<- Entry
	Stop()
}

// Client we don't care how Client sends the data
type Client interface {
	Send(b []byte) error
	// add Stop() functionality to the interface
}
