package te

import "testing"

func TestEsctest2DECICDefault(t *testing.T) {
	screen := NewScreen(8, 2)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcdefg", "ABCDEFG"})
	setCursor(screen, 2, 1)
	if err := stream.Feed(ControlCSI + "'}"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 0) != "a bcdefg" {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 1) != "A BCDEFG" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}

func TestEsctest2DECICExplicit(t *testing.T) {
	screen := NewScreen(9, 3)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcdefg", "ABCDEFG", "zyxwvut"})
	setCursor(screen, 2, 2)
	if err := stream.Feed(ControlCSI + "2'}"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 0) != "a  bcdefg" {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 2) != "z  yxwvut" {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}

func TestEsctest2DECDCDefault(t *testing.T) {
	screen := NewScreen(8, 2)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcdefg", "ABCDEFG"})
	setCursor(screen, 2, 1)
	if err := stream.Feed(ControlCSI + "'~"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 0) != "acdefg  " {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 1) != "ACDEFG  " {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}

func TestEsctest2DECDCExplicit(t *testing.T) {
	screen := NewScreen(9, 3)
	stream := NewStream(screen, false)
	fillScreenLines(screen, []string{"abcdefg", "ABCDEFG", "zyxwvut"})
	setCursor(screen, 2, 2)
	if err := stream.Feed(ControlCSI + "2'~"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 0) != "adefg    " {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 2) != "zwvut    " {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}
