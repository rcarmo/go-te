package te

import "testing"

// From esctest2/esctest/tests/cr.py::test_CR_Basic
func TestEsctestCrTestCRBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 3})
	esctestWrite(t, stream, ControlCR)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 3})
}

// From esctest2/esctest/tests/cr.py::test_CR_MovesToLeftMarginWhenRightOfLeftMargin
func TestEsctestCrTestCRMovesToLeftMarginWhenRightOfLeftMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 1})
	esctestWrite(t, stream, ControlCR)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 1})
}

// From esctest2/esctest/tests/cr.py::test_CR_MovesToLeftOfScreenWhenLeftOfLeftMargin
func TestEsctestCrTestCRMovesToLeftOfScreenWhenLeftOfLeftMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 1})
	esctestWrite(t, stream, ControlCR)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/cr.py::test_CR_StaysPutWhenAtLeftMargin
func TestEsctestCrTestCRStaysPutWhenAtLeftMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, ControlCR)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 1})
}

// From esctest2/esctest/tests/cr.py::test_CR_MovesToLeftMarginWhenLeftOfLeftMarginInOriginMode
func TestEsctestCrTestCRMovesToLeftMarginWhenLeftOfLeftMarginInOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 1})
	esctestWrite(t, stream, ControlCR)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestWrite(t, stream, "x")
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 1, Right: 5, Bottom: 1}, []string{"x"})
}
