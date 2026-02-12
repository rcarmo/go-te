package te

import "testing"

func TestEsctest2REPBasic(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("a"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "3b"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0][:4] != "aaaa" {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}
