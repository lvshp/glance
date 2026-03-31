package reader

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestPaginateContentWrapsLongParagraphs(t *testing.T) {
	lines := paginateContent(strings.Repeat("word ", 25), 20)
	if len(lines) < 2 {
		t.Fatalf("expected wrapped lines, got %d", len(lines))
	}

	for _, line := range lines {
		if runewidth.StringWidth(line) > 20 {
			t.Fatalf("line too wide: %q", line)
		}
	}
}

func TestPaginateContentSkipsBlankLines(t *testing.T) {
	lines := paginateContent("first\n\n second \n", 80)
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(lines))
	}

	if lines[0] != "first" || lines[1] != "second" {
		t.Fatalf("unexpected lines: %#v", lines)
	}
}

func TestCurrentViewReturnsMultipleLines(t *testing.T) {
	r := &contentReader{
		content: []string{"one", "two", "three"},
		pos:     0,
	}

	if got := r.CurrentView(2); got != "one\ntwo" {
		t.Fatalf("current view = %q, want %q", got, "one\ntwo")
	}
}

func TestWrapLineWrapsChineseWithoutSpaces(t *testing.T) {
	lines := wrapLine("这是一个没有空格但是需要自动换行的中文段落", 8)
	if len(lines) < 2 {
		t.Fatalf("expected wrapped lines, got %d", len(lines))
	}

	for _, line := range lines {
		if runewidth.StringWidth(line) > 8 {
			t.Fatalf("line too wide: %q", line)
		}
	}
}
