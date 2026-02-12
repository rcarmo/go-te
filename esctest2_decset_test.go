package te

import "testing"

func TestEsctest2DECSETDECCOLM(t *testing.T) {
	screen := NewScreen(80, 5)
	screen.SetMargins(1, 2)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(1, 2)
	setCursor(screen, 5, 5)
	screen.Draw("x")
	screen.SetMode([]int{3}, true)
	if screen.Columns != 132 {
		t.Fatalf("expected 132 columns")
	}
	assertCursor(t, screen, 1, 1)
	assertCell(t, screen, 5, 5, " ")
}

func TestEsctest2DECSETDECOM(t *testing.T) {
	screen := NewScreen(80, 10)
	screen.SetMargins(5, 7)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 7)
	screen.SetMode([]int{6}, true)
	setCursor(screen, 1, 1)
	screen.Draw("X")
	screen.ResetMode([]int{6}, true)
	screen.SetMargins(0, 0)
	screen.ResetMode([]int{69}, true)
	assertCell(t, screen, 5, 5, "X")
}
