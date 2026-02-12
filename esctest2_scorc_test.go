package te

import "testing"

func TestEsctest2SCORCBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed(ControlCSI + "u"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 5, 6)
}

func TestEsctest2SCORCMoveHomeWhenNotSaved(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "u"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2SCORCResetsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.SetMargins(5, 7)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 7)
	screen.SetMode([]int{6}, true)
	if err := stream.Feed(ControlCSI + "u"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.CursorPosition(1, 1)
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2SCORCWorksInLRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 2, 3)
	if err := stream.Feed(ControlCSI + "s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(1, 10)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 4, 5)
	if err := stream.Feed(ControlCSI + "u"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 2, 3)
}
