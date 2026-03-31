package reader

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

const defaultLineWidth = 80

type contentReader struct {
	content   []string
	pos       int
	rawText   string
	lineWidth int
}

func (r *contentReader) Reflow(width int) {
	if width <= 0 {
		width = defaultLineWidth
	}

	progress := r.progressRatio()
	r.lineWidth = width
	r.content = paginateContent(r.rawText, width)
	r.restoreProgress(progress)
}

func (r *contentReader) Current() string {
	if len(r.content) == 0 {
		return ""
	}

	return r.content[r.pos]
}

func (r *contentReader) CurrentView(lines int) string {
	if len(r.content) == 0 {
		return ""
	}

	if lines <= 1 {
		return r.Current()
	}

	end := r.pos + lines
	if end > len(r.content) {
		end = len(r.content)
	}

	return strings.Join(r.content[r.pos:end], "\n")
}

func (r *contentReader) Next() string {
	if len(r.content) == 0 {
		return ""
	}

	r.pos++
	if r.pos <= len(r.content)-1 {
		return r.content[r.pos]
	}

	r.pos = len(r.content) - 1
	return "END"
}

func (r *contentReader) Prev() string {
	if len(r.content) == 0 {
		return ""
	}

	r.pos--
	if r.pos < 0 {
		r.pos = 0
	}

	return r.content[r.pos]
}

func (r *contentReader) First() string {
	if len(r.content) == 0 {
		return ""
	}

	r.pos = 0
	return r.content[0]
}

func (r *contentReader) Last() string {
	if len(r.content) == 0 {
		return ""
	}

	r.pos = len(r.content) - 1
	return r.content[len(r.content)-1]
}

func (r *contentReader) CurrentPos() int {
	return r.pos
}

func (r *contentReader) Goto(pos int) string {
	if len(r.content) == 0 {
		return ""
	}

	if pos < 0 {
		pos = 0
	}

	if pos > len(r.content)-1 {
		pos = len(r.content) - 1
	}

	r.pos = pos
	return r.content[r.pos]
}

func (r *contentReader) GetProgress() string {
	if len(r.content) == 0 {
		return "(0 / 0)"
	}

	return fmt.Sprintf("(%d / %d)", r.pos+1, len(r.content))
}

func (r *contentReader) CurrentChapterTitle() string {
	return ""
}

func (r *contentReader) CurrentChapterIndex() int {
	return 0
}

func (r *contentReader) NextChapter() string {
	return r.Current()
}

func (r *contentReader) PrevChapter() string {
	return r.Current()
}

func (r *contentReader) GetTOC() string {
	return "No table of contents available."
}

func (r *contentReader) GetTOCWithSelection(selected, pageSize int) string {
	return r.GetTOC()
}

func (r *contentReader) GotoChapter(index int) string {
	return r.Current()
}

func (r *contentReader) setContent(text string) {
	r.rawText = text
	r.lineWidth = defaultLineWidth
	r.content = paginateContent(text, r.lineWidth)
	r.pos = 0
}

func (r *contentReader) progressRatio() float64 {
	if len(r.content) <= 1 {
		return 0
	}

	return float64(r.pos) / float64(len(r.content)-1)
}

func (r *contentReader) restoreProgress(progress float64) {
	if len(r.content) == 0 {
		r.pos = 0
		return
	}

	if progress <= 0 {
		r.pos = 0
		return
	}

	if progress >= 1 {
		r.pos = len(r.content) - 1
		return
	}

	r.pos = int(progress * float64(len(r.content)-1))
}

func paginateContent(text string, width int) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	paragraphs := strings.Split(text, "\n")
	lines := make([]string, 0, len(paragraphs))

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		lines = append(lines, wrapLine(paragraph, width)...)
	}

	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func wrapLine(text string, width int) []string {
	if runewidth.StringWidth(text) <= width {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) <= 1 {
		return wrapRunes(text, width)
	}

	lines := make([]string, 0, len(words))
	current := words[0]
	currentWidth := runewidth.StringWidth(current)

	for _, word := range words[1:] {
		wordWidth := runewidth.StringWidth(word)
		if currentWidth+1+wordWidth <= width {
			current += " " + word
			currentWidth += 1 + wordWidth
			continue
		}

		lines = append(lines, current)
		current = word
		currentWidth = wordWidth
	}

	lines = append(lines, current)
	return lines
}

func wrapRunes(text string, width int) []string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) == 0 {
		return nil
	}

	lines := make([]string, 0)
	start := 0
	for start < len(runes) {
		end := start
		lineWidth := 0
		for end < len(runes) {
			runeWidth := runewidth.RuneWidth(runes[end])
			if runeWidth == 0 {
				runeWidth = 1
			}

			if end > start && lineWidth+runeWidth > width {
				break
			}

			lineWidth += runeWidth
			end++
		}

		line := strings.TrimSpace(string(runes[start:end]))
		if line != "" {
			lines = append(lines, line)
		}

		start = end
	}

	return lines
}
