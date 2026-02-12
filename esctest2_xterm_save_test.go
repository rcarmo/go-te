package te

import "testing"

func TestEsctest2XtermSaveSaveSetState(t *testing.T) {
	screen := NewScreen(80, 2)
	stream := NewStream(screen, false)
	screen.SetMode([]int{7}, true)
	if err := stream.Feed(ControlCSI + "?7s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.ResetMode([]int{7}, true)
	if err := stream.Feed(ControlCSI + "?7r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 79, 1)
	screen.Draw("xxx")
	assertCursor(t, screen, 2, 2)
}

func TestEsctest2XtermSaveSaveResetState(t *testing.T) {
	screen := NewScreen(80, 2)
	stream := NewStream(screen, false)
	screen.ResetMode([]int{7}, true)
	if err := stream.Feed(ControlCSI + "?7s"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.SetMode([]int{7}, true)
	if err := stream.Feed(ControlCSI + "?7r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 79, 1)
	screen.Draw("xxx")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at column %d, got %d", screen.Columns, screen.Cursor.Col)
	}
}
