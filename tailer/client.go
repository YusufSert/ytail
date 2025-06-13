package tailer

import (
	"fmt"
)

// WithClient sets a custom logger.
func WithClient(c Client) Option {
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
