package reader

import "testing"

func TestContentReaderReflowKeepsApproximateProgress(t *testing.T) {
	r := &contentReader{}
	r.setContent("第一行\n第二行\n第三行\n第四行\n第五行\n第六行")
	r.Goto(3)

	r.Reflow(4)

	if len(r.content) == 0 {
		t.Fatal("expected reflowed content")
	}

	if r.CurrentPos() < 0 || r.CurrentPos() >= len(r.content) {
		t.Fatalf("current pos out of range: %d", r.CurrentPos())
	}
}

func TestTXTAnchorRestoreKeepsChapterAfterReflow(t *testing.T) {
	r := NewTxtReader()
	r.setContent("第1章 开始\n第一段内容第一段内容第一段内容\n第2章 中段\n第二段内容第二段内容第二段内容\n第3章 结尾\n第三段内容第三段内容第三段内容")
	r.buildChapters()
	r.Reflow(10)
	r.GotoChapter(1)
	r.Next()

	anchor := AnchorFromReader(r)
	r.Reflow(18)
	RestoreFromAnchor(r, anchor)

	if got := r.CurrentChapterIndex(); got != 1 {
		t.Fatalf("CurrentChapterIndex() = %d, want 1", got)
	}
}

func TestEPUBAnchorRestoreKeepsChapterAfterReflow(t *testing.T) {
	r := NewEpubReader()
	r.title = "demo"
	r.chapterText = []string{
		"第一章 开始\n第一章内容第一章内容第一章内容第一章内容",
		"第二章 中段\n第二章内容第二章内容第二章内容第二章内容",
		"第三章 结尾\n第三章内容第三章内容第三章内容第三章内容",
	}
	r.rebuildChapters([]string{"第一章 开始", "第二章 中段", "第三章 结尾"}, 10)
	r.GotoChapter(1)
	r.Next()

	anchor := AnchorFromReader(r)
	r.Reflow(18)
	RestoreFromAnchor(r, anchor)

	if got := r.CurrentChapterIndex(); got != 1 {
		t.Fatalf("CurrentChapterIndex() = %d, want 1", got)
	}
}

func TestFindChapterIndexByTitle(t *testing.T) {
	r := NewEpubReader()
	r.chapterText = []string{"a", "b", "c"}
	r.rebuildChapters([]string{"第一章 开始", "第二十章 投资价值", "第二十六章 故居"}, 20)

	if got := FindChapterIndexByTitle(r, "第二十章  投资价值"); got != 1 {
		t.Fatalf("FindChapterIndexByTitle() = %d, want 1", got)
	}
}
