package exec

import (
	"testing"
)

func TestParseTSV(t *testing.T) {
	input := "name\tage\nAlice\t30\nBob\t25\n"
	headers, rows, err := ParseTSV(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 2 || headers[0] != "name" || headers[1] != "age" {
		t.Errorf("unexpected headers: %v", headers)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0][0] != "Alice" || rows[0][1] != "30" {
		t.Errorf("unexpected row 0: %v", rows[0])
	}
	if rows[1][0] != "Bob" || rows[1][1] != "25" {
		t.Errorf("unexpected row 1: %v", rows[1])
	}
}

func TestParseTSVSingleColumn(t *testing.T) {
	input := "count\n42\n"
	headers, rows, err := ParseTSV(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 1 || headers[0] != "count" {
		t.Errorf("unexpected headers: %v", headers)
	}
	if len(rows) != 1 || rows[0][0] != "42" {
		t.Errorf("unexpected rows: %v", rows)
	}
}

func TestParseTSVHeaderOnly(t *testing.T) {
	input := "name\tage\n"
	headers, rows, err := ParseTSV(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(headers))
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParseTSVEmpty(t *testing.T) {
	_, _, err := ParseTSV("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseTSVTrailingNewlines(t *testing.T) {
	input := "a\tb\n1\t2\n\n"
	headers, rows, err := ParseTSV(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(headers))
	}
	if len(rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(rows))
	}
}
