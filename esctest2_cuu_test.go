package te

import "testing"

func TestEsctest2CUUBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	screen.CursorPosition(4, 4)
	screen.CursorUp(2)
	assertCursor(t, screen, 4, 2)
}
