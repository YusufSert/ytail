package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"ytail/common"
)

type Client struct {
	in     chan *value
	c      http.Client
	config Config
	err    error
	stop   func()
	logger *slog.Logger
}

func New(config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		in:     make(chan *value, 1000),
		c:      http.Client{Timeout: 3 * time.Second},
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

func NewWithOptions(config Config, opts ...Option) *Client {
	v := New(config)
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

func (c *Client) Send(b []byte) error {
	c.in <- &value{
		Ts:   time.Now(),
		Line: string(b),
	}
	return c.Err()
}

func (c *Client) run(ctx context.Context) {
	batch := make([][]any, 0)
	d := c.config.BatchMaxWait
	maxSize := c.config.BatchMaxSize

	t := time.NewTimer(d)
	for {
		select {
		case <-ctx.Done():
			return

		case v := <-c.in:
			batch = append(batch, v.List())
			if len(batch) <= maxSize {
				continue
			}
			c.send(batch)
			batch = make([][]any, 0)

		case <-t.C:
			if len(batch) > 0 {
				c.send(batch)
				batch = make([][]any, 0)
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

func (c *Client) send(b [][]any) {
	s := stream{
		Labels: c.config.Labels,
		Values: b,
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
	c.logger.Info("logs are pushed", "url", c.config.PushURL)
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
		backoff = backoff << 2
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
	Values [][]any           `json:"values"`
}

type value struct {
	Ts   time.Time      `json:"ts"`
	Line string         `json:"line"`
	Meta map[string]any `json:"meta"`
}

func (v *value) List() []any {
	var a []any
	a = append(a, v.nano(), v.Line)
	return a
}

func (v *value) nano() string {
	return strconv.FormatInt(v.Ts.UnixNano(), 10)
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
