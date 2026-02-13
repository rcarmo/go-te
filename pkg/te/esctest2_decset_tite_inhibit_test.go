package te

import "testing"

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_Basic
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_MoveToHomeWhenNotSaved
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorMoveToHomeWhenNotSaved(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_ResetsOriginMode
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorResetsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECSTBM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_WorksInLRM
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorWorksInLRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 1, 10)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 5})
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_AltVsMain
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorAltVsMain(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 7})
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECRESET(t, stream, esctestModeAltBuf)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 7})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_Protection
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECSCA(t, stream, 0)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestWrite(t, stream, "a")
	esctestDECSERA(t, stream, 1, 1, 1, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"a"})
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_Wrap
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorWrap(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "abcd")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_ReverseWrapNotAffected
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorReverseWrapNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestDECRESET(t, stream, esctestReverseWraparoundMode())
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/decset_tite_inhibit.py::test_SaveRestoreCursor_InsertNotAffected
func TestEsctestDecsetTiteInhibitTestSaveRestoreCursorInsertNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSET(t, stream, esctestModeSaveRestoreCursor)
	esctestRM(t, stream, esctestModeIRM)
	esctestDECRESET(t, stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"b" + esctestEmpty()})
}
