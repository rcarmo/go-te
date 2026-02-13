package te

import "testing"

// From esctest2/esctest/tests/cbt.py::test_CBT_OneTabStopByDefault
func TestEsctestCbtTestCBTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 17, Y: 1})
	esctestCBT(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cbt.py::test_CBT_ExplicitParameter
func TestEsctestCbtTestCBTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 25, Y: 1})
	esctestCBT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cbt.py::test_CBT_StopsAtLeftEdge
func TestEsctestCbtTestCBTStopsAtLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 25, Y: 2})
	esctestCBT(t, stream, 5)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 2)
}

// From esctest2/esctest/tests/cbt.py::test_CBT_IgnoresRegion
func TestEsctestCbtTestCBTIgnoresRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 30)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 9})
	esctestCBT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}
