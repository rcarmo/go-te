package te

import "testing"

// From esctest2/esctest/tests/ind.py::test_IND_Basic
func TestEsctestIndTestINDBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestIND(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/ind.py::test_IND_Scrolls
func TestEsctestIndTestINDScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_ScrollsInTopBottomRegionStartingAbove
func TestEsctestIndTestINDScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_ScrollsInTopBottomRegionStartingWithin
func TestEsctestIndTestINDScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_MovesDoesNotScrollOutsideLeftRight
func TestEsctestIndTestINDMovesDoesNotScrollOutsideLeftRight(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 5)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	height := esctestGetScreenSize(screen).Height
	esctestCUP(t, stream, esctestPoint{X: 6, Y: height})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: 5})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: height})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})
}

// From esctest2/esctest/tests/ind.py::test_IND_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestIndTestINDStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestIND(t, stream)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}

// From esctest2/esctest/tests/ind.py::test_IND_8bit
func TestEsctestIndTestIND8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
