package te

import "testing"

func TestEsctest2RISClearsScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	screen.Draw("x")
	if err := stream.Feed(ControlESC + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, " ")
}

func TestEsctest2RISCursorToOrigin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlESC + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2RISResetTabs(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 1)
	screen.SetTabStop()
	setCursor(screen, 10, 1)
	screen.SetTabStop()
	if err := stream.Feed(ControlESC + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.Tab()
	assertCursor(t, screen, 9, 1)
}

func TestEsctest2RISResetDECCOLM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	screen.SetMode([]int{3}, true)
	if screen.Columns != 132 {
		t.Fatalf("expected 132 columns")
	}
	if err := stream.Feed(ControlESC + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Columns != 80 {
		t.Fatalf("expected 80 columns, got %d", screen.Columns)
	}
}

func TestEsctest2RISResetDECOMAndMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	screen.SetMargins(5, 7)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 7)
	screen.SetMode([]int{6}, true)
	if err := stream.Feed(ControlESC + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	screen.Draw("X")
	assertCell(t, screen, 1, 1, "X")
	setCursor(screen, 3, 4)
	screen.CursorBack(1)
	assertCursor(t, screen, 2, 4)
	screen.CursorUp(1)
	assertCursor(t, screen, 2, 3)
}
