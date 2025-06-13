package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	_, err := ParseConfig(".././ytail-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
}
