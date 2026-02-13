package te

import "testing"

// From esctest2/esctest/tests/cud.py::test_CUD_DefaultParam
func TestEsctestCudTestCUDDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUD(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/cud.py::test_CUD_ExplicitParam
func TestEsctestCudTestCUDExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUD(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 5)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomLine
func TestEsctestCudTestCUDStopsAtBottomLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCUD(t, stream, height)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestCudTestCUDStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	height := esctestGetScreenSize(screen).Height
	esctestCUD(t, stream, height)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomMarginInScrollRegion
func TestEsctestCudTestCUDStopsAtBottomMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUD(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 4)
}
