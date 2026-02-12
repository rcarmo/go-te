package te

import "testing"

func TestEsctest2DECSERABasic(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	if err := stream.Feed(ControlCSI + "5;5;7;7${"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 4) != "ABCD   H" {
		t.Fatalf("unexpected row: %q", rowString(screen, 4))
	}
}
