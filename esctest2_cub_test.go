package te

import "testing"

func TestEsctest2CUBBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	screen.CursorPosition(1, 5)
	screen.CursorBack(2)
	assertCursor(t, screen, 3, 1)
}
