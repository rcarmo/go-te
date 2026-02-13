package te

import "testing"

// From esctest2/esctest/tests/hpr.py::test_HPR_DefaultParams
func TestEsctestHprTestHPRDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 1})
	esctestHPR(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_StopsAtRightEdge
func TestEsctestHprTestHPRStopsAtRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestHPR(t, stream, size.Width+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_DoesNotChangeRow
func TestEsctestHprTestHPRDoesNotChangeRow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestHPR(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_IgnoresOriginMode
func TestEsctestHprTestHPRIgnoresOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestWrite(t, stream, "X")
	esctestHPR(t, stream, 2)
	esctestWrite(t, stream, "Y")
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 7, Right: 9, Bottom: 7}, []string{esctestEmpty() + "X" + esctestEmpty() + esctestEmpty() + "Y"})
}
