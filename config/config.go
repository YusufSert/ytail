package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"ytail/client"
	"ytail/tailer"
)

// Config is a global config
type Config struct {
	// Tailer
	TailerConfig tailer.Config `yaml:"tailer"`

	// Client
	ClientConfig client.Config `yaml:"client"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(file)
	c := &Config{}
	if err != nil {
		return nil, err
	}
	err = dec.Decode(c)
	if err != nil {
		return nil, err
	}
	err = c.Validate()
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *Config) Validate() error {
	if c.TailerConfig.ScrapePath == "" {
		return errors.New("config: empty watch_path field")
	}
	if c.ClientConfig.Retry < 0 {
		return errors.New("config: empty retry field")
	}
	if c.ClientConfig.Backoff < 1 {
		return errors.New("config: empty back_off field")
	}
	if c.ClientConfig.MaxBackoff < 1 {
		return errors.New("config: empty max_back_off field")
	}
	_, err := url.Parse(c.ClientConfig.PushURL)
	if err != nil {
		return err
	}
	if c.ClientConfig.BatchMaxSize < 0 {
		return errors.New("config: empty batch_max_size field")
	}
	if c.ClientConfig.BatchMaxWait < 1 {
		return errors.New("config: empty batch_max_wait field")
	}
	if len(c.ClientConfig.Labels) < 1 {
		return errors.New("config: empty labels field")
	}
	return nil
}
