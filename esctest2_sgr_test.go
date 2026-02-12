package te

import "testing"

func TestEsctest2SGRBold(t *testing.T) {
	screen := NewScreen(10, 2)
	screen.Draw("x")
	screen.SelectGraphicRendition([]int{1}, false)
	screen.Draw("y")
	if screen.Buffer[0][0].Attr.Bold {
		t.Fatalf("expected first cell not bold")
	}
	if !screen.Buffer[0][1].Attr.Bold {
		t.Fatalf("expected second cell bold")
	}
}
