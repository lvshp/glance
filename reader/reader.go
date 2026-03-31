package reader

type Reader interface {
	Load(path string) error
	Reflow(width int)
	Current() string
	CurrentView(lines int) string
	Next() string
	Prev() string
	First() string
	Last() string
	CurrentPos() int
	Goto(pos int) string
	GetProgress() string
	CurrentChapterTitle() string
	CurrentChapterIndex() int
	NextChapter() string
	PrevChapter() string
	GetTOC() string
	GetTOCWithSelection(selected, pageSize int) string
	GotoChapter(index int) string
}
