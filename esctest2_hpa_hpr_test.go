package te

import "testing"

func TestEsctest2HPA(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "6`"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 6, 1)
	if err := stream.Feed(ControlCSI + "`"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 1)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "100`"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 80, 6)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "2`"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 2, 6)
	screen.SetMargins(6, 11)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 7, 9)
	screen.SetMode([]int{6}, true)
	if err := stream.Feed(ControlCSI + "2`"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Col != 1 {
		t.Fatalf("expected column 2, got %d", screen.Cursor.Col+1)
	}
}

func TestEsctest2HPR(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 6, 1)
	if err := stream.Feed(ControlCSI + "a"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 7, 1)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "200a"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 80, 6)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "2a"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 7, 6)
	screen.SetMargins(6, 11)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	screen.SetMode([]int{6}, true)
	setCursor(screen, 2, 2)
	screen.Draw("X")
	if err := stream.Feed(ControlCSI + "2a"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.Draw("Y")
	screen.ResetMode([]int{6}, true)
	screen.SetMargins(0, 0)
	screen.ResetMode([]int{69}, true)
	assertCell(t, screen, 6, 7, "X")
	assertCell(t, screen, 9, 7, "Y")
}
