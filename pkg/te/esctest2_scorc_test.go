package te

import "testing"

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Basic
func TestEsctestScorcTestSaveRestoreCursorBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestSCOSC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestSCORC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_MoveToHomeWhenNotSaved
func TestEsctestScorcTestSaveRestoreCursorMoveToHomeWhenNotSaved(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestSCORC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Reset
func TestEsctestScorcTestSaveRestoreCursorReset(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestSCOSC(t, stream)
	esctestDECSTR(t, stream)
	esctestWrite(t, stream, "b")
	esctestSCORC(t, stream)
	esctestWrite(t, stream, "c")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"cb"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_ResetsOriginMode
func TestEsctestScorcTestSaveRestoreCursorResetsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestSCOSC(t, stream)
	esctestDECSTBM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestSCORC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_WorksInLRM
func TestEsctestScorcTestSaveRestoreCursorWorksInLRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestSCOSC(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 1, 10)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestSCOSC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 5})
	esctestSCORC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_AltVsMain
func TestEsctestScorcTestSaveRestoreCursorAltVsMain(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestSCOSC(t, stream)
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 7})
	esctestSCOSC(t, stream)
	esctestDECRESET(t, stream, esctestModeAltBuf)
	esctestSCORC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestSCORC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 7})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Protection
func TestEsctestScorcTestSaveRestoreCursorProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	esctestSCOSC(t, stream)
	esctestDECSCA(t, stream, 0)
	esctestSCORC(t, stream)
	esctestWrite(t, stream, "a")
	esctestDECSERA(t, stream, 1, 1, 1, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"a"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Wrap
func TestEsctestScorcTestSaveRestoreCursorWrap(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestSCOSC(t, stream)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestSCORC(t, stream)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "abcd")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_ReverseWrapNotAffected
func TestEsctestScorcTestSaveRestoreCursorReverseWrapNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestSCOSC(t, stream)
	esctestDECRESET(t, stream, esctestReverseWraparoundMode())
	esctestSCORC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_InsertNotAffected
func TestEsctestScorcTestSaveRestoreCursorInsertNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestSCOSC(t, stream)
	esctestRM(t, stream, esctestModeIRM)
	esctestSCORC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"b" + esctestEmpty()})
}
