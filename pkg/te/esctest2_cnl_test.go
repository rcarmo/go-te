package te

import "testing"

// From esctest2/esctest/tests/cnl.py::test_CNL_DefaultParam
func TestEsctestCnlTestCNLDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCNL(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_ExplicitParam
func TestEsctestCnlTestCNLExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	esctestCNL(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 5)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_StopsAtBottomLine
func TestEsctestCnlTestCNLStopsAtBottomLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCNL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, height)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestCnlTestCNLStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 6})
	height := esctestGetScreenSize(screen).Height
	esctestCNL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, height)
	esctestAssertEQ(t, pos.X, 5)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_StopsAtBottomMarginInScrollRegion
func TestEsctestCnlTestCNLStopsAtBottomMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 3})
	esctestCNL(t, stream, 99)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 4)
	esctestAssertEQ(t, pos.X, 5)
}
