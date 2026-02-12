package te

import "testing"

func TestEsctest2ELBasic(t *testing.T) {
	screen := NewScreen(6, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("abcdef"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 3, 1)
	if err := stream.Feed(ControlCSI + "0K"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "ab    " {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}
