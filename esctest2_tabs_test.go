package te

import "testing"

func TestEsctest2CBTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 17, 1)
	screen.CursorBackTab(1)
	assertCursor(t, screen, 9, 1)
}

func TestEsctest2CBTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 25, 1)
	screen.CursorBackTab(2)
	assertCursor(t, screen, 9, 1)
}

func TestEsctest2CBTStopsAtLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 25, 2)
	screen.CursorBackTab(5)
	assertCursor(t, screen, 1, 2)
}

func TestEsctest2CBTIgnoresRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 30)
	setCursor(screen, 7, 9)
	screen.CursorBackTab(2)
	assertCursor(t, screen, 1, 9)
}

func TestEsctest2CHTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.CursorForwardTab(1)
	assertCursor(t, screen, 9, 1)
}

func TestEsctest2CHTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.CursorForwardTab(2)
	assertCursor(t, screen, 17, 1)
}

func TestEsctest2CHTIgnoresScrollingRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 30)
	setCursor(screen, 7, 9)
	screen.CursorForwardTab(2)
	assertCursor(t, screen, 17, 9)
	screen.CursorForwardTab(2)
	assertCursor(t, screen, 30, 9)
	setCursor(screen, 1, 9)
	screen.CursorForwardTab(9)
	assertCursor(t, screen, 30, 9)
}
