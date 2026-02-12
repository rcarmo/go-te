package te

import "testing"

func TestEsctest2INDBasic(t *testing.T) {
	screen := NewScreen(10, 3)
	stream := NewStream(screen, false)
	setCursor(screen, 2, 2)
	if err := stream.Feed(ControlESC + "D"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 2, 3)
}
