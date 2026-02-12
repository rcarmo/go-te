package te

import "testing"

func TestEsctest2NELBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 5, 3)
	screen.NextLine()
	assertCursor(t, screen, 1, 4)
}

func TestEsctest2NELScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	height := screen.Lines
	setCursor(screen, 2, height-1)
	screen.Draw("a")
	setCursor(screen, 2, height)
	screen.Draw("b")
	setCursor(screen, 2, height-1)
	screen.NextLine()
	assertCursor(t, screen, 1, height)
	assertCell(t, screen, 2, height-1, "a")
	assertCell(t, screen, 2, height, "b")
	screen.NextLine()
	assertCursor(t, screen, 1, height)
	assertCell(t, screen, 2, height-1, "b")
	assertCell(t, screen, 2, height, " ")
}

func TestEsctest2NELScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(4, 5)
	setCursor(screen, 2, 5)
	screen.Draw("x")
	setCursor(screen, 2, 3)
	screen.NextLine()
	screen.NextLine()
	screen.NextLine()
	assertCursor(t, screen, 1, 5)
	assertCell(t, screen, 2, 4, "x")
	assertCell(t, screen, 2, 5, " ")
}

func TestEsctest2NELScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(4, 5)
	setCursor(screen, 2, 5)
	screen.Draw("x")
	setCursor(screen, 2, 4)
	screen.NextLine()
	screen.NextLine()
	assertCursor(t, screen, 1, 5)
	assertCell(t, screen, 2, 4, "x")
	assertCell(t, screen, 2, 5, " ")
}

func TestEsctest2NELMovesDoesNotScrollOutsideLeftRight(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(2, 5)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 5)
	setCursor(screen, 3, 5)
	screen.Draw("x")
	setCursor(screen, 6, 5)
	screen.NextLine()
	assertCursor(t, screen, 2, 5)
	assertCell(t, screen, 3, 5, "x")
	setCursor(screen, 6, 4)
	screen.NextLine()
	assertCursor(t, screen, 2, 5)
	setCursor(screen, 6, screen.Lines)
	screen.NextLine()
	assertCursor(t, screen, 2, screen.Lines)
	assertCell(t, screen, 3, 5, "x")
	setCursor(screen, 1, 5)
	screen.NextLine()
	assertCursor(t, screen, 1, 5)
	assertCell(t, screen, 3, 5, "x")
	setCursor(screen, 1, screen.Lines)
	screen.NextLine()
	assertCursor(t, screen, 1, screen.Lines)
	assertCell(t, screen, 3, 5, "x")
}

func TestEsctest2NELStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(4, 5)
	setCursor(screen, 1, 6)
	screen.Draw("x")
	for i := 0; i < screen.Lines; i++ {
		screen.NextLine()
	}
	assertCursor(t, screen, 1, screen.Lines)
	assertCell(t, screen, 1, 6, "x")
}

func TestEsctest2NEL8Bit(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 3)
	if err := stream.Feed("\x1b G\u0085"); err != nil {
		t.Fatalf("feed error: %v", err)
	}
	assertCursor(t, screen, 1, 4)
}
