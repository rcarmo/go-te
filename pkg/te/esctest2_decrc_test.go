package te

import "testing"

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Basic
func TestEsctestDecrcTestSaveRestoreCursorBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_MoveToHomeWhenNotSaved
func TestEsctestDecrcTestSaveRestoreCursorMoveToHomeWhenNotSaved(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Reset
func TestEsctestDecrcTestSaveRestoreCursorReset(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestDECSC(t, stream)
	esctestDECSTR(t, stream)
	esctestWrite(t, stream, "b")
	esctestDECRC(t, stream)
	esctestWrite(t, stream, "c")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"cb"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_ResetsOriginMode
func TestEsctestDecrcTestSaveRestoreCursorResetsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestDECSTBM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 7)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestDECRC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_WorksInLRM
func TestEsctestDecrcTestSaveRestoreCursorWorksInLRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestDECSC(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 1, 10)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 5})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_AltVsMain
func TestEsctestDecrcTestSaveRestoreCursorAltVsMain(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestDECSC(t, stream)
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 7})
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestModeAltBuf)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 7})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Protection
func TestEsctestDecrcTestSaveRestoreCursorProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	esctestDECSC(t, stream)
	esctestDECSCA(t, stream, 0)
	esctestDECRC(t, stream)
	esctestWrite(t, stream, "a")
	esctestDECSERA(t, stream, 1, 1, 1, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"a"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Wrap
func TestEsctestDecrcTestSaveRestoreCursorWrap(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestDECRC(t, stream)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "abcd")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_ReverseWrapNotAffected
func TestEsctestDecrcTestSaveRestoreCursorReverseWrapNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestReverseWraparoundMode())
	esctestDECRC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_InsertNotAffected
func TestEsctestDecrcTestSaveRestoreCursorInsertNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSC(t, stream)
	esctestRM(t, stream, esctestModeIRM)
	esctestDECRC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"b" + esctestEmpty()})
}
