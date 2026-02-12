package te

import "testing"

func TestEsctest2VPRBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "2e"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 3)
}
