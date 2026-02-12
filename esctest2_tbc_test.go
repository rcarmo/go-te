package te

import "testing"

func TestEsctest2TBCDefault(t *testing.T) {
	screen := NewScreen(80, 1)
	screen.Tab()
	assertCursor(t, screen, 9, 1)
	screen.ClearTabStop(0)
	setCursor(screen, 1, 1)
	screen.Tab()
	assertCursor(t, screen, 17, 1)
}

func TestEsctest2TBC0(t *testing.T) {
	screen := NewScreen(80, 1)
	screen.Tab()
	assertCursor(t, screen, 9, 1)
	screen.ClearTabStop(0)
	setCursor(screen, 1, 1)
	screen.Tab()
	assertCursor(t, screen, 17, 1)
}

func TestEsctest2TBC3(t *testing.T) {
	screen := NewScreen(80, 1)
	screen.ClearTabStop(3)
	setCursor(screen, 30, 1)
	screen.SetTabStop()
	setCursor(screen, 1, 1)
	screen.Tab()
	assertCursor(t, screen, 30, 1)
}

func TestEsctest2TBCNoOp(t *testing.T) {
	screen := NewScreen(80, 1)
	setCursor(screen, 10, 1)
	screen.ClearTabStop(0)
	setCursor(screen, 1, 1)
	screen.Tab()
	assertCursor(t, screen, 9, 1)
	screen.Tab()
	assertCursor(t, screen, 17, 1)
}
