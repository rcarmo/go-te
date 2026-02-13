package te

import "testing"

// From esctest2/esctest/tests/hpa.py::test_HPA_DefaultParams
func TestEsctestHpaTestHPADefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHPA(t, stream, 6)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestHPA(t, stream)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
}

// From esctest2/esctest/tests/hpa.py::test_HPA_StopsAtRightEdge
func TestEsctestHpaTestHPAStopsAtRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestHPA(t, stream, size.Width+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpa.py::test_HPA_DoesNotChangeRow
func TestEsctestHpaTestHPADoesNotChangeRow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestHPA(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpa.py::test_HPA_IgnoresOriginMode
func TestEsctestHpaTestHPAIgnoresOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 9})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
	esctestAssertEQ(t, pos.Y, 9)

	esctestDECSET(t, stream, esctestModeDECOM)
	esctestHPA(t, stream, 2)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
}
