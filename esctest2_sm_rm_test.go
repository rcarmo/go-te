package te

import "testing"

func TestEsctest2SMIRM(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("abc"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed(ControlCSI + "4h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0][:4] != "Xabc" {
		t.Fatalf("expected Xabc, got %q", screen.Display()[0][:4])
	}
}

func TestEsctest2SMIRMMargins(t *testing.T) {
	screen := NewScreen(20, 2)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 1)
	if err := stream.Feed("abcdef"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 7, 1)
	if err := stream.Feed(ControlCSI + "4h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.ResetMode([]int{69}, true)
	if got := screen.Display()[0][4:11]; got != "abXcde " {
		t.Fatalf("expected abXcde , got %q", got)
	}
}

func TestEsctest2SMIRMWrapBehavior(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("aaaaaaaaa"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("b"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed(ControlCSI + "4h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 2, " ")
	setCursor(screen, screen.Columns, 1)
	if err := stream.Feed("YZ"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 2, "Z")
}

func TestEsctest2SMLNM(t *testing.T) {
	screen := NewScreen(10, 3)
	stream := NewStream(screen, false)
	setCursor(screen, 5, 1)
	if err := stream.Feed(ControlCSI + "20l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("\n"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 5, 2)
	if err := stream.Feed(ControlCSI + "20h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 5, 1)
	if err := stream.Feed("\n"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 2)
}

func TestEsctest2RMIRM(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "4h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed("X"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed("W"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 1)
	if err := stream.Feed(ControlCSI + "4l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("YZ"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0][:2] != "YZ" {
		t.Fatalf("expected YZ, got %q", screen.Display()[0][:2])
	}
}
