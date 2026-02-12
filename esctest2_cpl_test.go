package te

import "testing"

func TestEsctest2CPLBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	screen.CursorPosition(4, 4)
	if err := stream.Feed(ControlCSI + "2F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 2)
}
