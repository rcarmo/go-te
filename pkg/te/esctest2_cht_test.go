package te

import "testing"

// From esctest2/esctest/tests/cht.py::test_CHT_OneTabStopByDefault
func TestEsctestChtTestCHTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCHT(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cht.py::test_CHT_ExplicitParameter
func TestEsctestChtTestCHTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCHT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
}

// From esctest2/esctest/tests/cht.py::test_CHT_IgnoresScrollingRegion
func TestEsctestChtTestCHTIgnoresScrollingRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 30)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 9})
	esctestCHT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
	esctestCHT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 30)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 9})
	esctestCHT(t, stream, 9)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 30)
}
