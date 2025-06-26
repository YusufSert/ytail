package client

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestClientPush(t *testing.T) {
	l := slog.New(slog.NewTextHandler(os.Stdin, nil))

	c := NewWithOptions(Config{
		PushURL:      "http://172.26.98.110:3100/loki/api/v1/push",
		Retry:        10,
		Backoff:      time.Second,
		MaxBackoff:   time.Minute,
		BatchMaxSize: 100,
		BatchMaxWait: time.Second,
		Labels:       map[string]string{"foo": "bar", "service_name": "test"},
	}, WithLogger(l))

	var err error
	for i := 0; i < 1000; i++ {
		err = c.Send([]byte("'log-line':'bok'"))
		time.Sleep(time.Millisecond * 100)
	}
	if err != nil {
		t.Fatal(err)
	}
}
