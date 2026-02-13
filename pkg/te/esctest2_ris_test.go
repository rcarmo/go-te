package te

import "testing"

// From esctest2/esctest/tests/ris.py::test_RIS_ClearsScreen
func TestEsctestRisTestRISClearsScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "x")
	esctestRIS(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ris.py::test_RIS_CursorToOrigin
func TestEsctestRisTestRISCursorToOrigin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestRIS(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetTabs
func TestEsctestRisTestRISResetTabs(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHTS(t, stream)
	esctestCUF(t, stream)
	esctestHTS(t, stream)
	esctestCUF(t, stream)
	esctestHTS(t, stream)
	esctestRIS(t, stream)
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 9, Y: 1})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetTitleMode
func TestEsctestRisTestRISResetTitleMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryUTF8)
	esctestSMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryHex)
	esctestRIS(t, stream)
	esctestChangeWindowTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "ab")
	esctestChangeWindowTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "a")
	esctestChangeIconTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "ab")
	esctestChangeIconTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "a")
}

// From esctest2/esctest/tests/ris.py::test_RIS_ExitAltScreen
func TestEsctestRisTestRISExitAltScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "m")
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestRIS(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetDECCOLM
func TestEsctestRisTestRISResetDECCOLM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeAllow80To132)
	esctestDECSET(t, stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(screen).Width, 132)
	esctestRIS(t, stream)
	esctestAssertEQ(t, esctestGetScreenSize(screen).Width, 80)
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetDECOM
func TestEsctestRisTestRISResetDECOM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestRIS(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/ris.py::test_RIS_RemoveMargins
func TestEsctestRisTestRISRemoveMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 5)
	esctestDECSTBM(t, stream, 4, 6)
	esctestRIS(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 4})
	esctestCUB(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 4})
	esctestCUU(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestCUF(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 6})
	esctestCUD(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 7})
}
