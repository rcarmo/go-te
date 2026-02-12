package te

import "testing"

func TestEsctest2S8C1TCSI(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed("\x1b G\u009b1;2H\x1b F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 2, 1)
}
