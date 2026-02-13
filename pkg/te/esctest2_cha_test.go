package te

import "testing"

// From esctest2/esctest/tests/cha.py::test_CHA_DefaultParam
func TestEsctestChaTestCHADefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_ExplicitParam
func TestEsctestChaTestCHAExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 10)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_OutOfBoundsLarge
func TestEsctestChaTestCHAOutOfBoundsLarge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 9999)
	pos := esctestGetCursorPosition(screen)
	width := esctestGetScreenSize(screen).Width
	esctestAssertEQ(t, pos.X, width)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_ZeroParam
func TestEsctestChaTestCHAZeroParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 0)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_IgnoresScrollRegion
func TestEsctestChaTestCHAIgnoresScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 1)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_RespectsOriginMode
func TestEsctestChaTestCHARespectsOriginMode(t *testing.T) {
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
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestCHA(t, stream, 1)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestWrite(t, stream, "X")

	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 6, Right: 5, Bottom: 6}, []string{"X"})
}
