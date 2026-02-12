package te

import "testing"

func TestEsctest2DCSUnrecognized(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlDCS + "z0" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, " ")
}

func TestEsctest2DECSTBMScrollsOnNewline(t *testing.T) {
	screen := NewScreen(10, 4)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "2;3r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 2)
	screen.Draw("1")
	screen.CarriageReturn()
	screen.LineFeed()
	screen.Draw("2")
	assertCell(t, screen, 1, 2, "1")
	assertCell(t, screen, 1, 3, "2")
	screen.CarriageReturn()
	screen.LineFeed()
	assertCell(t, screen, 1, 2, "2")
	assertCell(t, screen, 1, 3, " ")
	assertCursor(t, screen, 1, 3)
}

func TestEsctest2DECSTBMMovesCursorToOrigin(t *testing.T) {
	screen := NewScreen(10, 4)
	stream := NewStream(screen, false)
	setCursor(screen, 3, 2)
	if err := stream.Feed(ControlCSI + "2;3r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2DECSTBMDefaultRestores(t *testing.T) {
	screen := NewScreen(10, 4)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "2;3r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	setCursor(screen, 1, 2)
	screen.Draw("1")
	screen.CarriageReturn()
	screen.LineFeed()
	screen.Draw("2")
	pos := screen.Cursor
	if err := stream.Feed(ControlCSI + "r"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.Cursor = pos
	screen.CarriageReturn()
	screen.LineFeed()
	assertCell(t, screen, 1, 2, "1")
	assertCell(t, screen, 1, 3, "2")
	assertCursor(t, screen, 1, 4)
}
