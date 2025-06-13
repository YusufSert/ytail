package client

import "time"

type Config struct {
	Labels       map[string]string `yaml:"labels"`
	Retry        int               `yaml:"retry"`
	Backoff      time.Duration     `yaml:"backoff"`
	MaxBackoff   time.Duration     `yaml:"max_backoff"`
	PushURL      string            `yaml:"push_url"`
	BatchMaxSize int               `yaml:"batch_max_size"`
	BatchMaxWait time.Duration     `yaml:"batch_max_wait"`
}
