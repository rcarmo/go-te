package te

import "testing"

func TestEsctest2FFBasic(t *testing.T) {
	screen := NewScreen(10, 3)
	setCursor(screen, 2, 2)
	screen.LineFeed()
	assertCursor(t, screen, 2, 3)
}
