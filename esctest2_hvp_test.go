package te

import "testing"

func TestEsctest2HVPBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "2;3f"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 3, 2)
}
