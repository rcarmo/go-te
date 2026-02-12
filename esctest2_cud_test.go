package te

import "testing"

func TestEsctest2CUDBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	screen.CursorPosition(1, 1)
	screen.CursorDown(2)
	assertCursor(t, screen, 1, 3)
}
