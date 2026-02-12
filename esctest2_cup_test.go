package te

import "testing"

func TestEsctest2CUPBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "3;4H"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 4, 3)
}
