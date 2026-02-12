package te

import "testing"

func TestEsctest2DECFIBasic(t *testing.T) {
	screen := NewScreen(80, 10)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlESC + "9"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 6, 6)
}

func TestEsctest2DECFINoWrapOnRightEdge(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	setCursor(screen, 10, 2)
	if err := stream.Feed(ControlESC + "9"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 10, 2)
}

func TestEsctest2DECFIScrolls(t *testing.T) {
	screen := NewScreen(5, 5)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	screen.SetMargins(2, 4)
	setCursor(screen, 4, 3)
	if err := stream.Feed(ControlESC + "9"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 1) != "fhi j" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}

func TestEsctest2DECBIBasic(t *testing.T) {
	screen := NewScreen(80, 10)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlESC + "6"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 4, 6)
}

func TestEsctest2DECBINoWrapOnLeftEdge(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	setCursor(screen, 1, 2)
	if err := stream.Feed(ControlESC + "6"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 2)
}

func TestEsctest2DECBIScrolls(t *testing.T) {
	screen := NewScreen(5, 5)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	screen.SetMargins(2, 4)
	setCursor(screen, 2, 3)
	if err := stream.Feed(ControlESC + "6"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 1) != "f ghj" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}
