package te

import "testing"

// From esctest2/esctest/tests/ff.py::test_FF_Basic
func TestEsctestFfTestFFBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestWrite(t, stream, ControlFF)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/ff.py::test_FF_Scrolls
func TestEsctestFfTestFFScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_ScrollsInTopBottomRegionStartingAbove
func TestEsctestFfTestFFScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_ScrollsInTopBottomRegionStartingWithin
func TestEsctestFfTestFFScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_MovesDoesNotScrollOutsideLeftRight
func TestEsctestFfTestFFMovesDoesNotScrollOutsideLeftRight(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 5)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 6, Y: 4})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 5})

	height := esctestGetScreenSize(screen).Height
	esctestCUP(t, stream, esctestPoint{X: 6, Y: height})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: 5})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: height})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})
}

// From esctest2/esctest/tests/ff.py::test_FF_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestFfTestFFStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestWrite(t, stream, ControlFF)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}
