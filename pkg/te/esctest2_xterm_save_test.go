package te

import "testing"

// From esctest2/esctest/tests/xterm_save.py::test_XtermSave_SaveSetState
func TestEsctestXtermSaveTestXtermSaveSaveSetState(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestXtermSave(t, stream, esctestModeDECAWM)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestXtermRestore(t, stream, esctestModeDECAWM)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "xxx")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
}

// From esctest2/esctest/tests/xterm_save.py::test_XtermSave_SaveResetState
func TestEsctestXtermSaveTestXtermSaveSaveResetState(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestXtermSave(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestXtermRestore(t, stream, esctestModeDECAWM)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "xxx")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, size.Width)
}
