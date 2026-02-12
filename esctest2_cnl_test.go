package te

import "testing"

func TestEsctest2CNLBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	screen.CursorPosition(1, 1)
	if err := stream.Feed(ControlCSI + "2E"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 3)
}
