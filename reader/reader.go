package reader

type Reader interface {
	Load(path string) error
	Reflow(width int)
	BookTitle() string
	Current() string
	CurrentView(lines int) string
	Total() int
	Search(query string, start int, forward bool) (int, bool)
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
