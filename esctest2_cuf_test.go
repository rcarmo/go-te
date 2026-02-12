package te

import "testing"

func TestEsctest2CUFBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	screen.CursorPosition(1, 1)
	screen.CursorForward(3)
	assertCursor(t, screen, 4, 1)
}
