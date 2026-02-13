package te

import "testing"

// From esctest2/esctest/tests/vpa.py::test_VPA_DefaultParams
func TestEsctestVpaTestVPADefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestVPA(t, stream, 6)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 6)
	esctestVPA(t, stream)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_StopsAtBottomEdge
func TestEsctestVpaTestVPAStopsAtBottomEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	size := esctestGetScreenSize(screen)
	esctestVPA(t, stream, size.Height+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_DoesNotChangeColumn
func TestEsctestVpaTestVPADoesNotChangeColumn(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestVPA(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_IgnoresOriginMode
func TestEsctestVpaTestVPAIgnoresOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 9})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 9)
	esctestAssertEQ(t, pos.X, 7)

	esctestDECSET(t, stream, esctestModeDECOM)
	esctestVPA(t, stream, 2)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 2)
}
