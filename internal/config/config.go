package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type SwayConfig struct {
	filename string
	Outputs  []string
}

func NewSwayConfig(filename string) *SwayConfig {
	return &SwayConfig{filename: filename}
}

func (s *SwayConfig) Read(marker []string) ([]string, error) {
	fmt.Println(strings.Join(marker, ", "))
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", s.filename, err)
	}

	_, start, end, err := findConfigMarker(marker, f)
	if err != nil {
		return nil, fmt.Errorf("error getting configs from file %s: %w", s.filename, err)
	} else if start == -1 {
		return nil, fmt.Errorf("no matching marker found")
	}

	sr := io.NewSectionReader(f, start, end-start)

	return Parse(sr)
}
