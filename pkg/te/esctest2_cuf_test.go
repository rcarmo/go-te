package te

import "testing"

// From esctest2/esctest/tests/cuf.py::test_CUF_DefaultParam
func TestEsctestCufTestCUFDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUF(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_ExplicitParam
func TestEsctestCufTestCUFExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestCUF(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 3)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_StopsAtRightSide
func TestEsctestCufTestCUFStopsAtRightSide(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	width := esctestGetScreenSize(screen).Width
	esctestCUF(t, stream, width)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, width)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_StopsAtRightEdgeWhenBegunRightOfScrollRegion
func TestEsctestCufTestCUFStopsAtRightEdgeWhenBegunRightOfScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 12, Y: 3})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 12)
	width := esctestGetScreenSize(screen).Width
	esctestCUF(t, stream, width)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, width)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_StopsAtRightMarginInScrollRegion
func TestEsctestCufTestCUFStopsAtRightMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 3})
	width := esctestGetScreenSize(screen).Width
	esctestCUF(t, stream, width)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 10)
}
