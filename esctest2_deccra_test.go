package te

import "testing"

func TestEsctest2DECCRAOverlap(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	if err := stream.Feed(ControlCSI + "2;2;4;4;1;3;3;1$v"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 2) != "qrjklvwx" {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}

func TestEsctest2DECCRANonOverlap(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	if err := stream.Feed(ControlCSI + "2;2;4;4;1;5;5;1$v"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 4) != "ABCDjklH" {
		t.Fatalf("unexpected row: %q", rowString(screen, 4))
	}
}
