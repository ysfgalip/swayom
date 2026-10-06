package config

import (
	"fmt"
	"os"
)

type SwayConfig struct {
	filename string
	Outputs  []string
}

func NewSwayConfig(filename string) *SwayConfig {
	return &SwayConfig{filename: filename}
}

func (s *SwayConfig) Read() ([]string, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", s.filename, err)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting output configs from file %s: %w", s.filename, err)
	}

	return Parse(f)
}
