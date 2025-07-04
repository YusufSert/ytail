package tailer

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"regexp"
	"ytail/api"
	"ytail/common"

	"github.com/fsnotify/fsnotify"
)

type Tail struct {
	r         *bufio.Reader
	rc        io.ReadCloser
	err       error
	totalLine uint64
	w         *fsnotify.Watcher
	c         api.Client
	logger    *slog.Logger
	config    Config
}

func New(config Config) (*Tail, error) {
	t := &Tail{
		c:      &Kudim{},
		logger: slog.New(&common.DiscardHandler{}),
		config: config,
	}
	return t, nil
}

// Option configures Tail using the functional options paradigm popularized by Rob Pike and Dave Cheney.
// If you're unfamiliar with this style,
// see https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html and
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
type Option interface {
	apply(t *Tail)
}

type optionFunc func(t *Tail)

func (fn optionFunc) apply(t *Tail) {
	fn(t)
}

func NewWithOptions(config Config, opts ...Option) (*Tail, error) {
	t, err := New(config)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt.apply(t)
	}
	return t, nil
}

func (t *Tail) Run(ctx context.Context) error {
	err := t.setWatcher(t.config.ScrapePath)
	if err != nil {
		return err
	}

	err = t.trySetFile(t.config.ScrapePath)
	if err != nil {
		return err
	}

	defer t.w.Close()
	defer t.rc.Close()

loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		select {
		case err, ok := <-t.w.Errors:
			if !ok {
				break loop
			}
			return err
		case e, ok := <-t.w.Events:
			if !ok {
				break loop
			}

			if e.Has(fsnotify.Write) {
				t.logger.Debug("tail: notify write")
				line, err := t.readLine()
				if err != nil {
					return err
				}

				if len(line) > 0 {
					err = t.c.Send(line)
					if err != nil {
						return err
					}
				}
			}

			//todo: maybe specify flag for ths process. true|false, use regex to get correct file because editor create extra file(regex for *-log.txt)
			if e.Has(fsnotify.Create) {
				t.logger.Debug("tail: notify create")
				err := t.newReadCloser(e.Name)
				if err != nil {
					return err
				}
				t.logger.Info("tail: watching new file", "file", e.Name)
			}
		}
	}
	return nil
}

// note: write package for easy error handling: op, msg fields needed
func (t *Tail) readLine() ([]byte, error) {
	line, err := t.r.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return line, nil
}

func (t *Tail) trySetFile(dirPath string) error {
	var err error
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	//Look for log files, if found one set it to tail.
	// If no file found return nil, if there is file and error occurred return the error
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if re := regexp.MustCompile(t.config.FileRegex); re.MatchString(path.Base(e.Name())) != true {
			continue
		}

		err = t.newReadCloser(path.Join(dirPath, e.Name()))
		if err == nil {
			t.logger.Info("tail: file set to be watch", "file", e.Name())
			return nil
		}
	}
	// if there is no error and no file found
	if err == nil {
		t.logger.Warn("tail: no log file found", "watchPath", dirPath)
		return errors.New("no log file found") // if no file found return nil, and wait watch for new file creation.
	}
	return err
}

func (t *Tail) newReadCloser(name string) error {
	if t.rc != nil {
		t.rc.Close()
	}

	rc, err := os.Open(name)
	if err != nil {
		return err
	}

	t.rc = rc
	t.r = bufio.NewReader(rc)
	return nil
}

func (t *Tail) setWatcher(path string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.Add(path)
	if err != nil {
		return err
	}
	t.w = w
	t.logger.Info("tail: start watching changes", "path", path)
	return nil
}

func (t *Tail) WatchList() []string {
	return t.w.WatchList()
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return optionFunc(func(t *Tail) {
		t.logger = l
	})
}
