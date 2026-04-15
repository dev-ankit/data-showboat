package markdown

import (
	"strings"
	"testing"
)

func TestParseTitle(t *testing.T) {
	input := "# My Demo\n\n*2026-02-06T15:30:00Z by Showboat v0.3.0*\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	tb, ok := blocks[0].(TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.Title != "My Demo" {
		t.Errorf("expected title 'My Demo', got %q", tb.Title)
	}
	if tb.Timestamp != "2026-02-06T15:30:00Z" {
		t.Errorf("expected timestamp '2026-02-06T15:30:00Z', got %q", tb.Timestamp)
	}
	if tb.Version != "v0.3.0" {
		t.Errorf("expected version 'v0.3.0', got %q", tb.Version)
	}
}

func TestParseTitleNoVersion(t *testing.T) {
	input := "# My Demo\n\n*2026-02-06T15:30:00Z*\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	tb, ok := blocks[0].(TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.Timestamp != "2026-02-06T15:30:00Z" {
		t.Errorf("expected timestamp '2026-02-06T15:30:00Z', got %q", tb.Timestamp)
	}
	if tb.Version != "" {
		t.Errorf("expected empty version, got %q", tb.Version)
	}
}

func TestParseCommentary(t *testing.T) {
	input := "# Demo\n\n*2026-02-06T00:00:00Z by Showboat v0.3.0*\n\nHello world.\n\nMore text here.\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d: %+v", len(blocks), blocks)
	}
	cb, ok := blocks[1].(CommentaryBlock)
	if !ok {
		t.Fatalf("expected CommentaryBlock, got %T", blocks[1])
	}
	if cb.Text != "Hello world.\n\nMore text here." {
		t.Errorf("unexpected text: %q", cb.Text)
	}
}

func TestParseCodeAndOutput(t *testing.T) {
	input := "```bash\necho hello\n```\n\n```output\nhello\n```\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	code, ok := blocks[0].(CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if code.Lang != "bash" || code.Code != "echo hello" || code.IsImage {
		t.Errorf("unexpected code block: %+v", code)
	}
	out, ok := blocks[1].(OutputBlock)
	if !ok {
		t.Fatalf("expected OutputBlock, got %T", blocks[1])
	}
	if out.Content != "hello\n" {
		t.Errorf("unexpected output: %q", out.Content)
	}
}

func TestParseImageCodeAndOutput(t *testing.T) {
	input := "```bash {image}\npython screenshot.py\n```\n\n![Screenshot](abc-2026-02-06.png)\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	code, ok := blocks[0].(CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if !code.IsImage {
		t.Error("expected IsImage=true")
	}
	img, ok := blocks[1].(ImageOutputBlock)
	if !ok {
		t.Fatalf("expected ImageOutputBlock, got %T", blocks[1])
	}
	if img.AltText != "Screenshot" || img.Filename != "abc-2026-02-06.png" {
		t.Errorf("unexpected image output: %+v", img)
	}
}

func TestParseOutputWithLongerFence(t *testing.T) {
	input := "````output\n```bash\necho hello\n```\n````\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d: %+v", len(blocks), blocks)
	}
	out, ok := blocks[0].(OutputBlock)
	if !ok {
		t.Fatalf("expected OutputBlock, got %T", blocks[0])
	}
	expected := "```bash\necho hello\n```\n"
	if out.Content != expected {
		t.Errorf("expected content:\n%s\ngot:\n%s", expected, out.Content)
	}
}

func TestParseCodeBlockWithLongerFence(t *testing.T) {
	input := "````bash\necho ```hello```\n````\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d: %+v", len(blocks), blocks)
	}
	code, ok := blocks[0].(CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if code.Code != "echo ```hello```" {
		t.Errorf("unexpected code: %q", code.Code)
	}
}

func TestRoundTripWithBackticksInOutput(t *testing.T) {
	input := "```bash\ncat inner.md\n```\n\n````output\n# My Demo\n\n```bash\necho hello\n```\n\n```output\nhello\n```\n````\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d: %+v", len(blocks), blocks)
	}
	var buf strings.Builder
	err = Write(&buf, blocks)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("round trip mismatch.\nexpected:\n%s\ngot:\n%s", input, buf.String())
	}
}

func TestParseTitleWithDocumentID(t *testing.T) {
	input := "# My Demo\n\n*2026-02-06T15:30:00Z by Showboat v0.3.0*\n<!-- showboat-id: abc-123 -->\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	tb, ok := blocks[0].(TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.DocumentID != "abc-123" {
		t.Errorf("expected DocumentID 'abc-123', got %q", tb.DocumentID)
	}
	if tb.Title != "My Demo" {
		t.Errorf("expected title 'My Demo', got %q", tb.Title)
	}
	if tb.Timestamp != "2026-02-06T15:30:00Z" {
		t.Errorf("expected timestamp '2026-02-06T15:30:00Z', got %q", tb.Timestamp)
	}
}

func TestParseTitleWithoutDocumentID(t *testing.T) {
	input := "# My Demo\n\n*2026-02-06T15:30:00Z*\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	tb, ok := blocks[0].(TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.DocumentID != "" {
		t.Errorf("expected empty DocumentID, got %q", tb.DocumentID)
	}
}

func TestParseTitleWithDocumentIDFollowedByContent(t *testing.T) {
	input := "# My Demo\n\n*2026-02-06T15:30:00Z*\n<!-- showboat-id: abc-123 -->\n\nHello world.\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d: %+v", len(blocks), blocks)
	}
	tb, ok := blocks[0].(TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.DocumentID != "abc-123" {
		t.Errorf("expected DocumentID 'abc-123', got %q", tb.DocumentID)
	}
	cb, ok := blocks[1].(CommentaryBlock)
	if !ok {
		t.Fatalf("expected CommentaryBlock, got %T", blocks[1])
	}
	if cb.Text != "Hello world." {
		t.Errorf("unexpected text: %q", cb.Text)
	}
}

func TestRoundTripWithDocumentID(t *testing.T) {
	input := "# Demo\n\n*2026-02-06T00:00:00Z by Showboat v0.3.0*\n<!-- showboat-id: test-uuid-456 -->\n\nLet's begin.\n\n```bash\necho hi\n```\n\n```output\nhi\n```\n\nDone.\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	err = Write(&buf, blocks)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("round trip mismatch.\nexpected:\n%s\ngot:\n%s", input, buf.String())
	}
}

func TestParseTableCodeBlock(t *testing.T) {
	input := "```python3 {table}\nprint('hello')\n```\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	code, ok := blocks[0].(CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if code.Lang != "python3" {
		t.Errorf("expected lang 'python3', got %q", code.Lang)
	}
	if !code.IsTable {
		t.Error("expected IsTable=true")
	}
	if code.IsImage {
		t.Error("expected IsImage=false")
	}
}

func TestParseTableOutput(t *testing.T) {
	input := "```python3 {table}\nprint('hello')\n```\n\n| name | age |\n| --- | --- |\n| Alice | 30 |\n| Bob | 25 |\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d: %+v", len(blocks), blocks)
	}
	code, ok := blocks[0].(CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if !code.IsTable {
		t.Error("expected IsTable=true")
	}

	tbl, ok := blocks[1].(TableOutputBlock)
	if !ok {
		t.Fatalf("expected TableOutputBlock, got %T", blocks[1])
	}
	if len(tbl.Headers) != 2 || tbl.Headers[0] != "name" || tbl.Headers[1] != "age" {
		t.Errorf("unexpected headers: %v", tbl.Headers)
	}
	if len(tbl.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows))
	}
	if tbl.Rows[0][0] != "Alice" || tbl.Rows[0][1] != "30" {
		t.Errorf("unexpected row 0: %v", tbl.Rows[0])
	}
	if tbl.Rows[1][0] != "Bob" || tbl.Rows[1][1] != "25" {
		t.Errorf("unexpected row 1: %v", tbl.Rows[1])
	}
}

func TestParseTableOutputHeadersOnly(t *testing.T) {
	input := "| name | age |\n| --- | --- |\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d: %+v", len(blocks), blocks)
	}
	tbl, ok := blocks[0].(TableOutputBlock)
	if !ok {
		t.Fatalf("expected TableOutputBlock, got %T", blocks[0])
	}
	if len(tbl.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(tbl.Headers))
	}
	if len(tbl.Rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(tbl.Rows))
	}
}

func TestRoundTripWithTable(t *testing.T) {
	input := "# Demo\n\n*2026-02-06T00:00:00Z by Showboat v0.3.0*\n\nQuerying the database.\n\n```python3 {table}\nimport sqlite3\n```\n\n| name | age |\n| --- | --- |\n| Alice | 30 |\n| Bob | 25 |\n\nDone.\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 5 {
		t.Fatalf("expected 5 blocks, got %d: %+v", len(blocks), blocks)
	}
	// Verify block types
	if _, ok := blocks[0].(TitleBlock); !ok {
		t.Errorf("block 0: expected TitleBlock, got %T", blocks[0])
	}
	if _, ok := blocks[1].(CommentaryBlock); !ok {
		t.Errorf("block 1: expected CommentaryBlock, got %T", blocks[1])
	}
	if cb, ok := blocks[2].(CodeBlock); !ok || !cb.IsTable {
		t.Errorf("block 2: expected CodeBlock with IsTable=true, got %T %+v", blocks[2], blocks[2])
	}
	if _, ok := blocks[3].(TableOutputBlock); !ok {
		t.Errorf("block 3: expected TableOutputBlock, got %T", blocks[3])
	}
	if _, ok := blocks[4].(CommentaryBlock); !ok {
		t.Errorf("block 4: expected CommentaryBlock, got %T", blocks[4])
	}

	var buf strings.Builder
	err = Write(&buf, blocks)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("round trip mismatch.\nexpected:\n%s\ngot:\n%s", input, buf.String())
	}
}

func TestRoundTrip(t *testing.T) {
	input := "# Demo\n\n*2026-02-06T00:00:00Z by Showboat v0.3.0*\n\nLet's begin.\n\n```bash\necho hi\n```\n\n```output\nhi\n```\n\nDone.\n"
	blocks, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	err = Write(&buf, blocks)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("round trip mismatch.\nexpected:\n%s\ngot:\n%s", input, buf.String())
	}
}
