package tailer

import (
	"fmt"
	"ytail/api"
)

// WithClient sets a custom logger.
func WithClient(c api.Client) Option {
	return optionFunc(func(t *Tail) {
		t.c = c
	})
}

type Kudim struct {
}

func (*Kudim) Send(b []byte) error {
	fmt.Println(b)
	return nil
}
