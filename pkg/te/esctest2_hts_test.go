package te

import "testing"

// From esctest2/esctest/tests/hts.py::test_HTS_Basic
func TestEsctestHtsTestHTSBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestTBC(t, stream, 3)
	esctestCUP(t, stream, esctestPoint{X: 20, Y: 1})
	esctestHTS(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 20)
}

// From esctest2/esctest/tests/hts.py::test_HTS_8bit
func TestEsctestHtsTestHTS8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
