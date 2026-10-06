package config

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(r io.Reader) ([]string, error) {
	var out []string
	var prefix []string
	temp := ""
	lineCount := 0
	sc := bufio.NewScanner(r)

	flush := func() {
		if temp != "" {
			out = append(out, strings.Join(append(prefix, temp), " "))
			temp = ""
		}
	}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		lineCount++

		if strings.HasPrefix(line, "#") {
			continue
		}

		for strings.HasSuffix(line, "\\") {
			if !sc.Scan() {
				if err := sc.Err(); err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("file ends with a line continuation")
			}
			lineCount++
			line = strings.TrimSpace(strings.TrimSpace(strings.TrimSuffix(line, "\\")) + " " + strings.TrimSpace(sc.Text()))
		}

		switch {
		case line == "":
			continue

		case line == "}":
			if len(prefix) == 0 {
				return nil, fmt.Errorf("invalid closing at: %d", lineCount)
			}
			flush()
			prefix = prefix[:len(prefix)-1]

		case line == "{":
			if temp == "" {
				return nil, fmt.Errorf("empty config block start at: %d", lineCount)
			}
			prefix = append(prefix, temp)
			temp = ""

		case strings.HasSuffix(line, "}"):
			if len(prefix) == 0 {
				return nil, fmt.Errorf("invalid closing at: %d", lineCount)
			}
			flush()
			out = append(out, strings.Join(append(prefix, strings.TrimSpace(strings.TrimSuffix(line, "}"))), " "))
			prefix = prefix[:len(prefix)-1]

		case strings.HasSuffix(line, "{"):
			flush()
			prefix = append(prefix, strings.TrimSpace(strings.TrimSuffix(line, "{")))

		default:
			flush()
			temp = line
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(prefix) != 0 {
		return nil, fmt.Errorf("unclosed block")
	}
	flush()

	return out, nil
}
