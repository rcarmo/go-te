package te

import "testing"

func TestEsctest2BSBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 3, 3)
	screen.Backspace()
	assertCursor(t, screen, 2, 3)
}

func TestEsctest2BSNoWrapByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 1, 3)
	screen.Backspace()
	assertCursor(t, screen, 1, 3)
}

func TestEsctest2BSWrapsInWraparoundMode(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 1045}, true)
	setCursor(screen, 1, 3)
	screen.Backspace()
	assertCursor(t, screen, 80, 2)
}

func TestEsctest2BSInitialReverseWraparound(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 45}, true)
	setCursor(screen, 1, 1)
	screen.NextLine()
	screen.Backspace()
	assertCursor(t, screen, 1, 2)
}

func TestEsctest2BSReverseWrapRequiresDECAWM(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.ResetMode([]int{7}, true)
	screen.SetMode([]int{1045}, true)
	setCursor(screen, 1, 3)
	screen.Backspace()
	assertCursor(t, screen, 1, 3)
}

func TestEsctest2BSReverseWrapWithLeftRight(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 1045, 69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 5, 3)
	screen.Backspace()
	assertCursor(t, screen, 10, 2)
}

func TestEsctest2BSReversewrapFromLeftEdgeToRightMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 1045, 69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 1, 3)
	screen.Backspace()
	assertCursor(t, screen, 10, 2)
}

func TestEsctest2BSReverseWrapGoesToBottom(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 1045}, true)
	screen.SetMargins(2, 5)
	setCursor(screen, 1, 2)
	screen.Backspace()
	assertCursor(t, screen, 80, 5)
}

func TestEsctest2BSStopsAtLeftMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 5, 1)
	screen.Backspace()
	assertCursor(t, screen, 5, 1)
}

func TestEsctest2BSMovesLeftWhenLeftOfLeftMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 4, 1)
	screen.Backspace()
	assertCursor(t, screen, 3, 1)
}

func TestEsctest2BSStopsAtOrigin(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 1, 1)
	screen.Backspace()
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2BSCursorStartsInDoWrapPosition(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 79, 1)
	screen.Draw("ab")
	screen.Backspace()
	screen.Draw("X")
	if screen.Display()[0][78:80] != "Xb" {
		t.Fatalf("expected Xb at end, got %q", screen.Display()[0][78:80])
	}
}

func TestEsctest2BSReverseWrapStartingInDoWrapPosition(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{7, 1045}, true)
	setCursor(screen, 79, 1)
	screen.Draw("ab")
	screen.Backspace()
	screen.Draw("X")
	if screen.Display()[0][78:80] != "aX" {
		t.Fatalf("expected aX at end, got %q", screen.Display()[0][78:80])
	}
}

func TestEsctest2CRBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	setCursor(screen, 3, 3)
	screen.CarriageReturn()
	assertCursor(t, screen, 1, 3)
}

func TestEsctest2CRWithLeftRightMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 6, 1)
	screen.CarriageReturn()
	assertCursor(t, screen, 5, 1)
	setCursor(screen, 4, 1)
	screen.CarriageReturn()
	assertCursor(t, screen, 1, 1)
}

func TestEsctest2CROriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMode([]int{69, 6}, true)
	screen.SetLeftRightMargins(5, 10)
	setCursor(screen, 4, 1)
	screen.CarriageReturn()
	assertCursor(t, screen, 5, 1)
}
