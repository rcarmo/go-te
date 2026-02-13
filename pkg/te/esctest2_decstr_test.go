package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECSC
func TestEsctestDecstrTestDECSTRDECSC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestDECSTR(t, stream)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_IRM
func TestEsctestDecstrTestDECSTRIRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"b"})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECOM
func TestEsctestDecstrTestDECSTRDECOM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 3, 4)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestDECSTR(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 4)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 4}, []string{"X" + strings.Repeat(esctestEmpty(), 2), strings.Repeat(esctestEmpty(), 3), strings.Repeat(esctestEmpty(), 3), strings.Repeat(esctestEmpty(), 3)})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECAWM
func TestEsctestDecstrTestDECSTRDECAWM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSTR(t, stream)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "xxx")
	position := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, position.X, 2)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_ReverseWraparound
func TestEsctestDecstrTestDECSTRReverseWraparound(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_STBM
func TestEsctestDecstrTestDECSTRSTBM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 3, 4)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 4})
	esctestWrite(t, stream, ControlCR+ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 5)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECSCA
func TestEsctestDecstrTestDECSTRDECSCA(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECSED(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECSASD
func TestEsctestDecstrTestDECSTRDECSASD(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSASD(t, stream, 1)
	esctestDECSTR(t, stream)
	position := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, position.X, 1)
	esctestAssertEQ(t, position.Y, 1)
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECRLM
func TestEsctestDecstrTestDECSTRDECRLM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECRLM)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 1}, []string{"a"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 1, Right: 3, Bottom: 1}, []string{"b"})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECLRMM
func TestEsctestDecstrTestDECSTRDECLRMM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 6)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestWrite(t, stream, "ab")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 7)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_CursorStaysPut
func TestEsctestDecstrTestDECSTRCursorStaysPut(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSTR(t, stream)
	position := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, position.X, 5)
	esctestAssertEQ(t, position.Y, 6)
}
