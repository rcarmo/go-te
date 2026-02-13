package te

import "testing"

// From esctest2/esctest/tests/cpl.py::test_CPL_DefaultParam
func TestEsctestCplTestCPLDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCPL(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_ExplicitParam
func TestEsctestCplTestCPLExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestCPL(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_StopsAtTopLine
func TestEsctestCplTestCPLStopsAtTopLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCPL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_StopsAtTopLineWhenBegunAboveScrollRegion
func TestEsctestCplTestCPLStopsAtTopLineWhenBegunAboveScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCPL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 1)
	esctestAssertEQ(t, pos.X, 5)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_StopsAtTopMarginInScrollRegion
func TestEsctestCplTestCPLStopsAtTopMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 3})
	esctestCPL(t, stream, 99)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 2)
	esctestAssertEQ(t, pos.X, 5)
}
