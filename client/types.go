package client

import "time"

type Stream struct {
    Labels string `json:"labels"`
    Entries[]
}

// Entry is a log entry with a timestamp.
type Entry struct {
    TimeStamp          time.Time `json:"ts"`
    Line               string    `json:"line"`
    StructuredMetadata map[string]string
}
