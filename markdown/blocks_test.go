package markdown

import "testing"

func TestCommentaryBlock(t *testing.T) {
	b := CommentaryBlock{Text: "Hello world\n\nThis is a test."}
	if b.Type() != "commentary" {
		t.Errorf("expected type commentary, got %s", b.Type())
	}
}

func TestCodeBlock(t *testing.T) {
	b := CodeBlock{Lang: "bash", Code: "echo hello", IsImage: false}
	if b.Type() != "code" {
		t.Errorf("expected type code, got %s", b.Type())
	}
}

func TestOutputBlock(t *testing.T) {
	b := OutputBlock{Content: "hello\n"}
	if b.Type() != "output" {
		t.Errorf("expected type output, got %s", b.Type())
	}
}

func TestImageOutputBlock(t *testing.T) {
	b := ImageOutputBlock{AltText: "Screenshot", Filename: "abc-2026-02-06.png"}
	if b.Type() != "output-image" {
		t.Errorf("expected type output-image, got %s", b.Type())
	}
}

func TestTitleBlock(t *testing.T) {
	b := TitleBlock{Title: "My Demo", Timestamp: "2026-02-06T15:30:00Z"}
	if b.Type() != "title" {
		t.Errorf("expected type title, got %s", b.Type())
	}
}

func TestTableOutputBlock(t *testing.T) {
	b := TableOutputBlock{
		Headers: []string{"name", "age"},
		Rows:    [][]string{{"Alice", "30"}, {"Bob", "25"}},
	}
	if b.Type() != "output-table" {
		t.Errorf("expected type output-table, got %s", b.Type())
	}
	if len(b.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(b.Headers))
	}
	if len(b.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(b.Rows))
	}
}
