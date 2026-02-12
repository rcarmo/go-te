package te

import "testing"

func fillRectScreen(screen *Screen) {
	lines := []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	}
	for i, line := range lines {
		setCursor(screen, 1, i+1)
		screen.Draw(line)
	}
}

func TestEsctest2DECERABasic(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	if err := stream.Feed(ControlCSI + "5;5;7;7$z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if rowString(screen, 4) != "ABCD   H" {
		t.Fatalf("unexpected row: %q", rowString(screen, 4))
	}
}

func TestEsctest2DECFRAOriginMode(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 8)
	screen.SetMargins(2, 8)
	screen.SetMode([]int{6}, true)
	if err := stream.Feed(ControlCSI + "37;1;1;3;3$x"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.ResetMode([]int{69}, true)
	screen.SetMargins(0, 0)
	screen.ResetMode([]int{6}, true)
	if rowString(screen, 1) != "i%%%mnop" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}

func TestEsctest2DECFRAIgnoresMargins(t *testing.T) {
	screen := NewScreen(8, 8)
	stream := NewStream(screen, false)
	fillRectScreen(screen)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(3, 6)
	screen.SetMargins(3, 6)
	if err := stream.Feed(ControlCSI + "37;5;5;7;7$x"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.ResetMode([]int{69}, true)
	screen.SetMargins(0, 0)
	if rowString(screen, 4) != "ABCD%%%H" {
		t.Fatalf("unexpected row: %q", rowString(screen, 4))
	}
}
