package te

import "testing"

func TestEsctest2SaveRestoreCursorBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlESC + "7"); err != nil {
		t.Fatalf("save: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed(ControlESC + "8"); err != nil {
		t.Fatalf("restore: %v", err)
	}
	assertCursor(t, screen, 5, 6)
}
