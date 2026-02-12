package te

import "testing"

func TestEsctest2ECHBasic(t *testing.T) {
	screen := NewScreen(6, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("abcdef"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 2, 1)
	if err := stream.Feed(ControlCSI + "2X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "a  def" {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}
