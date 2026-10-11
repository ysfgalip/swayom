package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"slices"
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

func findConfigMarker(marker []string, r io.Reader) (line int, start int64, length int64, err error) {
	var blockMarker []string
	lineNr := 0
	var offset int64 = 0
	var block int64 = -1
	blockLine := -1
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, 0, 0, fmt.Errorf("error while reading file at byte %d: %w", offset, err)
		}
		lineNr++
		lineStart := offset
		offset += int64(len(line))
		line = strings.TrimSpace(line)

		switch {
		case line == "@@":
			if block == -1 {
				return 0, 0, 0, fmt.Errorf("invalid block ending at byte %d", offset)
			}
			if sameElements(blockMarker, marker) {
				return blockLine, block, lineStart, nil
			}
			block = -1
			blockLine = -1
			blockMarker = nil

		case strings.HasPrefix(line, "@@"):
			if block > -1 {
				return 0, 0, 0, fmt.Errorf("block start before closing")
			}
			line = strings.TrimSpace(strings.TrimPrefix(line, "@@"))
			if errJSON := json.Unmarshal([]byte(line), &blockMarker); errJSON != nil {
				return 0, 0, 0, fmt.Errorf("invalid json array found at byte %d: %w", offset, errJSON)
			}
			block = offset
			blockLine = lineNr
		}

		if err == io.EOF {
			if block > -1 {
				return 0, 0, 0, fmt.Errorf("block unclosed at EOF")
			}
			return 0, block, 0, nil
		}
	}
}

func sameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	a2, b2 := slices.Clone(a), slices.Clone(b)
	slices.Sort(a2)
	slices.Sort(b2)
	return slices.Equal(a2, b2)
}
