package core

import "testing"

func TestHighlightSearchMatches(t *testing.T) {
	got := highlightSearchMatches("hello world\nHello again", "hello")
	want := "[hello](fg:black,bg:yellow,mod:bold) world\n[Hello](fg:black,bg:yellow,mod:bold) again"
	if got != want {
		t.Fatalf("highlightSearchMatches() = %q, want %q", got, want)
	}
}

func TestHighlightSearchMatchesChinese(t *testing.T) {
	got := highlightSearchMatches("我在三国当混蛋\n混蛋也要讲道理", "混蛋")
	want := "我在三国当[混蛋](fg:black,bg:yellow,mod:bold)\n[混蛋](fg:black,bg:yellow,mod:bold)也要讲道理"
	if got != want {
		t.Fatalf("highlightSearchMatches() = %q, want %q", got, want)
	}
}
