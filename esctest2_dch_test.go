package te

import "testing"

func TestEsctest2DCHBasic(t *testing.T) {
	screen := NewScreen(6, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("abcdef"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 2, 1)
	if err := stream.Feed(ControlCSI + "2P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "adef  " {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}
