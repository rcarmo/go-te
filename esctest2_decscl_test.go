package te

import "testing"

func TestEsctest2DECSCLReset(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	screen.Draw("x")
	if err := stream.Feed(ControlCSI + "62;1\"p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "          " {
		t.Fatalf("expected clear screen")
	}
	if screen.conformanceLevel != 2 {
		t.Fatalf("expected conformance level 2, got %d", screen.conformanceLevel)
	}
}
