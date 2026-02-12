package te

import "testing"

func TestEsctest2DECSEDEraseDisplay(t *testing.T) {
	screen := NewScreen(5, 2)
	stream := NewStream(screen, false)
	screen.Draw("abcde")
	setCursor(screen, 1, 2)
	screen.Draw("fghij")
	if err := stream.Feed(ControlCSI + "?2J"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, " ")
	assertCell(t, screen, 1, 2, " ")
}

func TestEsctest2DECSELEraseLine(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewStream(screen, false)
	screen.Draw("abcde")
	if err := stream.Feed(ControlCSI + "?2K"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, " ")
	assertCell(t, screen, 5, 1, " ")
}
