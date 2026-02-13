package te

import "testing"

// From esctest2/esctest/tests/ri.py::test_RI_Basic
func TestEsctestRiTestRIBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestRI(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/ri.py::test_RI_Scrolls
func TestEsctestRiTestRIScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestWrite(t, stream, "b")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 3}, []string{"a", "b", esctestEmpty()})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 3}, []string{esctestEmpty(), "a", "b"})
}

// From esctest2/esctest/tests/ri.py::test_RI_ScrollsInTopBottomRegionStartingBelow
func TestEsctestRiTestRIScrollsInTopBottomRegionStartingBelow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 6})
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 4})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{esctestEmpty(), "x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_ScrollsInTopBottomRegionStartingWithin
func TestEsctestRiTestRIScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 4})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{esctestEmpty(), "x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_MovesDoesNotScrollOutsideLeftRight
func TestEsctestRiTestRIMovesDoesNotScrollOutsideLeftRight(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 5)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 6, Y: 2})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 2})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 6, Y: 1})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 1})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 2})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})

	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 5, Right: 3, Bottom: 5}, []string{"x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_StopsAtTopLineWhenBegunAboveScrollRegion
func TestEsctestRiTestRIStopsAtTopLineWhenBegunAboveScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, "x")
	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestRI(t, stream)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 3, Right: 1, Bottom: 3}, []string{"x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_8bit
func TestEsctestRiTestRI8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
