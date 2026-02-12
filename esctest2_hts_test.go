package te

import "testing"

func TestEsctest2HTSBasic(t *testing.T) {
	screen := NewScreen(20, 1)
	setCursor(screen, 5, 1)
	screen.SetTabStop()
	setCursor(screen, 1, 1)
	screen.Tab()
	assertCursor(t, screen, 5, 1)
}
