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
