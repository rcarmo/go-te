package te

import "testing"

func TestEsctest2ECHDefault(t *testing.T) {
	screen := NewScreen(6, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("abcdef"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 2, 1)
	if err := stream.Feed(ControlCSI + "X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "a cdef" {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}
