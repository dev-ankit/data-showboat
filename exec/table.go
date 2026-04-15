package exec

import (
	"fmt"
	"strings"
)

// ParseTSV parses tab-separated output into headers and rows.
// The first line is treated as the header row. Empty trailing lines are ignored.
func ParseTSV(output string) ([]string, [][]string, error) {
	lines := strings.Split(output, "\n")
	// Remove empty trailing lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("empty output: no header row found")
	}

	headers := strings.Split(lines[0], "\t")
	var rows [][]string
	for _, line := range lines[1:] {
		rows = append(rows, strings.Split(line, "\t"))
	}
	return headers, rows, nil
}
