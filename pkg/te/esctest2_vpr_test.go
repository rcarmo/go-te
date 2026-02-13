package te

import "testing"

// From esctest2/esctest/tests/vpr.py::test_VPR_DefaultParams
func TestEsctestVprTestVPRDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestVPR(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 7)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_StopsAtBottomEdge
func TestEsctestVprTestVPRStopsAtBottomEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestVPR(t, stream, size.Height+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_DoesNotChangeColumn
func TestEsctestVprTestVPRDoesNotChangeColumn(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestVPR(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 8)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_IgnoresOriginMode
func TestEsctestVprTestVPRIgnoresOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestWrite(t, stream, "X")
	esctestVPR(t, stream, 2)
	esctestWrite(t, stream, "Y")
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 6, Top: 7, Right: 7, Bottom: 9}, []string{"X" + esctestEmpty(), esctestEmpty() + esctestEmpty(), esctestEmpty() + "Y"})
}
