package te

import "testing"

func TestEsctest2DECSTRDECSC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlESC + "7"); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if err := stream.Feed(ControlESC + "8"); err != nil {
		t.Fatalf("restore: %v", err)
	}
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2DECSTRIRM(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewStream(screen, false)
	screen.SetMode([]int{4}, false)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	setCursor(screen, 1, 1)
	screen.Draw("a")
	setCursor(screen, 1, 1)
	screen.Draw("b")
	assertCell(t, screen, 1, 1, "b")
}

func TestEsctest2DECSTRDECOM(t *testing.T) {
	screen := NewScreen(10, 6)
	stream := NewStream(screen, false)
	screen.SetMargins(3, 4)
	screen.SetMode([]int{6}, true)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(3, 4)
	screen.SetMargins(4, 5)
	setCursor(screen, 1, 1)
	screen.Draw("X")
	screen.ResetMode([]int{6}, true)
	screen.SetMargins(0, 0)
	screen.ResetMode([]int{69}, true)
	assertCell(t, screen, 1, 1, "X")
}

func TestEsctest2DECSTRReverseWraparound(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	screen.SetMode([]int{1045}, true)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	setCursor(screen, 1, 2)
	screen.Backspace()
	assertCursor(t, screen, 1, 2)
}

func TestEsctest2DECSTRSTBM(t *testing.T) {
	screen := NewScreen(10, 6)
	stream := NewStream(screen, false)
	screen.SetMargins(3, 4)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	setCursor(screen, 1, 4)
	screen.CarriageReturn()
	screen.LineFeed()
	assertCursor(t, screen, 1, 5)
}

func TestEsctest2DECSTRDECLRMM(t *testing.T) {
	screen := NewScreen(10, 6)
	stream := NewStream(screen, false)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 6)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	setCursor(screen, 5, 5)
	screen.Draw("ab")
	assertCursor(t, screen, 7, 5)
}

func TestEsctest2DECSTRCursorStaysPut(t *testing.T) {
	screen := NewScreen(10, 6)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "!p"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	assertCursor(t, screen, 5, 6)
}
