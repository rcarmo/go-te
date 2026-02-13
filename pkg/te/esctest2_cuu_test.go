package te

import "testing"

// From esctest2/esctest/tests/cuu.py::test_CUU_DefaultParam
func TestEsctestCuuTestCUUDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUU(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_ExplicitParam
func TestEsctestCuuTestCUUExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopLine
func TestEsctestCuuTestCUUStopsAtTopLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopLineWhenBegunAboveScrollRegion
func TestEsctestCuuTestCUUStopsAtTopLineWhenBegunAboveScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopMarginInScrollRegion
func TestEsctestCuuTestCUUStopsAtTopMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 2)
}
