package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"ytail/api"
	"ytail/common"
)

const (
	chanSize = 1 << 12
	timeout  = 3 * time.Second // client timeout
	bachSize = 1 << 8
)

type Client struct {
	in     chan api.Entry
	c      http.Client
	config Config
	err    error
	stop   func()
	logger *slog.Logger
}

// New creates new client and runs goroutine to send entries
func New(config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		in:     make(chan api.Entry, chanSize),
		c:      http.Client{Timeout: timeout},
		config: config,
		err:    nil,
		stop:   cancel,
		logger: slog.New(&common.DiscardHandler{}),
	}

	go client.run(ctx)

	return client
}

// Option configures Client using the functional options paradigm popularized by Rob Pike and Dave Cheney.
// If you're unfamiliar with this style,
// see https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html and
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
type Option interface {
	apply(v *Client)
}

type optionFunc func(v *Client)

func (fn optionFunc) apply(v *Client) {
	fn(v)
}

// todo: implement default values for config for easy testing
func NewWithOptions(config Config, opts ...Option) *Client {
	v := New(config)
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

// todo: fix this error
func (c *Client) Send(b []byte) error {
	c.in <- api.Entry{
		TimeStamp: time.Now(),
		Line:      string(b),
	}
	c.logger.Debug("loki-clint: adding new entry to batch")
	return c.Err()
}

func (c *Client) run(ctx context.Context) {
	d := c.config.BatchMaxWait
	batch := newBatch(0) // size not implemented you can give ever you want now

	t := time.NewTimer(d)
	for {
		select {
		case <-ctx.Done():
			return

		case e, ok := <-c.in:
			if !ok {
				// todo: maybe channel closed error
				return
			}
			if err := batch.add(e); err != nil {
				c.err = err
				break
			}
			c.logger.Info("loki-client: new entries added to batch", "batchSize", batch.sizeEntry())
			if batch.sizeEntry() <= c.config.BatchMaxSize {
				continue
			}
			c.send(batch)
			batch = newBatch(0)

		case <-t.C:
			if batch.sizeEntry() > 0 {
				c.send(batch)
				batch = newBatch(0)
			}
		}
		if !t.Stop() {
			select {
			case <-t.C:
			default:
			}
		}
		t.Reset(d)
	}
}

func (c *Client) send(b *batch) {
	s := stream{
		Labels: c.config.Labels,
		Values: b.entries,
	}
	body := postLoki{
		Streams: []stream{s},
	}
	data, err := json.Marshal(body)
	if err != nil {
		c.err = err
		return
	}

	err = c.retry(func() error {
		_, err = c.c.Post(c.config.PushURL, "application/json", bytes.NewReader(data))
		return err
	})

	if err != nil {
		c.err = err
		c.logger.Error("loki-client: error sending logs", "url", c.config.PushURL, "err", c.err)
		return
	}
	c.logger.Info("loki-client: logs are pushed", "url", c.config.PushURL, "size", len(data))
}

func (c *Client) retry(fn func() error) error {
	var err error
	backoff, maxBackOff := c.config.Backoff, c.config.MaxBackoff

	for i := 0; i < c.config.Retry; i++ {
		if backoff > maxBackOff {
			backoff = maxBackOff
		}

		err = fn()
		if err == nil {
			return err
		}
		slog.Warn("loki-client: error sending http request, retrying http request after "+backoff.String(), "err", err)
		time.Sleep(backoff)
		backoff = backoff << 1
	}
	return err
}

func (c *Client) Err() error {
	return c.err
}

// Stop stops sender goroutine from sending logs to [loki]
func (c *Client) Stop() {
	c.c.CloseIdleConnections()
	c.stop()
}

type stream struct {
	Labels map[string]string `json:"stream"`
	Values []api.Entry       `json:"values"`
}

type postLoki struct {
	Streams []stream `json:"streams"`
}

// test data for loki ingest
var testData = `
{"streams": [{ "stream": { "label": "bar2" }, "values": [ [ "1742124537511023448", "fizzbuzz info" ] ] }]}
`

// loki push api json format
// endpoint http://localhost:3100/loki/api/v1/push
// method POST
const format = `
{
  "streams": [
    {
      "stream": {
        "label": "value"
      },
      "values": [
          [ "<unix epoch in nanoseconds>", "<log line>", {key=val}], //You can optionally attach structured metadata
          [ "<unix epoch in nanoseconds>", "<log line>", {key=val}]
      ]
    }
  ]
}`

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return optionFunc(func(v *Client) {
		v.logger = l
	})
}
