package client

import (
	"time"
	"ytail/api"
)

// todo: we want to get size of the batch,

type batch struct {
	entries   []api.Entry
	bytes     int
	createdAt time.Time
}

func newBatch(maxEntry int, entries ...api.Entry) *batch {
	b := &batch{
		createdAt: time.Now(),
		bytes:     0,
	}
	for _, e := range entries {
		_ = b.add(e)
	}
	return b
}

func (b *batch) add(e api.Entry) error {
	b.bytes += 1 // todo: calculate entry-size
	b.entries = append(b.entries, e)
	return nil
}

func (b *batch) sizeEntry() int {
	return b.bytes
}

func (b *batch) age() time.Duration {
	return time.Since(b.createdAt)
}
