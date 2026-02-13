package te

import "testing"

// From esctest2/esctest/tests/decset.py::test_DECSET_DECCOLM
func TestEsctestDecsetTestDecsetDeccolm(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECRESET(t, fixture.stream, esctestModeDECNCSM)
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 5})
	esctestWrite(t, fixture.stream, "xyz")
	esctestDECSTBM(t, fixture.stream, 1, 2)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 1, 2)
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
	position := esctestGetCursorPosition(fixture.screen)
	esctestAssertEQ(t, position.X, 1)
	esctestAssertEQ(t, position.Y, 1)
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "Hello")
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "World")
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 5, Bottom: 3}, []string{"Hello", "World"})
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 5, Top: 5, Right: 5, Bottom: 5}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECOM
func TestEsctestDecsetTestDecsetDecom(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSTBM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECOM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 5, Top: 5, Right: 5, Bottom: 5}, []string{"X"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECOM_SoftReset
func TestEsctestDecsetTestDecsetDecomSoftReset(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSTBM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECOM)
	esctestDECSTR(t, fixture.stream)
	esctestCHA(t, fixture.stream, 1)
	esctestVPA(t, fixture.stream, 1)
	esctestWrite(t, fixture.stream, "X")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECOM_DECRQCRA
func TestEsctestDecsetTestDecsetDecomDecrqcra(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSTBM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 5, 7)
	esctestDECSET(t, fixture.stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"X"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECLRMM
func TestEsctestDecsetTestDecsetDeclrmm(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestWrite(t, fixture.stream, "abcdefgh")
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 3}, []string{"abcd", esctestEmpty() + "efg", esctestEmpty() + "h" + esctestEmpty() + esctestEmpty()})
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "ABCDEFGH")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 1}, []string{"ABCDEFGH"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECLRMM_MarginsResetByDECSTR
func TestEsctestDecsetTestDecsetDeclrmmMarginsResetByDecstr(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTR(t, fixture.stream)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 3})
	esctestWrite(t, fixture.stream, "abc")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 6)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECLRMM_ModeResetByDECSTR
func TestEsctestDecsetTestDecsetDeclrmmModeResetByDecstr(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTR(t, fixture.stream)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 3})
	esctestWrite(t, fixture.stream, "abc")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 6)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECNCSM
func TestEsctestDecsetTestDecsetDecncsm(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestDECSET(t, fixture.stream, esctestModeDECNCSM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "1")
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"1"})
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestDECSET(t, fixture.stream, esctestModeDECNCSM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "2")
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"2"})
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestDECRESET(t, fixture.stream, esctestModeDECNCSM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "3")
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestDECRESET(t, fixture.stream, esctestModeDECNCSM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "4")
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_SaveRestoreCursor
func TestEsctestDecsetTestDecsetSaveRestoreCursor(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, fixture.stream, esctestModeSaveRestoreCursor)
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 5})
	esctestDECRESET(t, fixture.stream, esctestModeSaveRestoreCursor)
	cursor := esctestGetCursorPosition(fixture.screen)
	esctestAssertEQ(t, cursor.X, 2)
	esctestAssertEQ(t, cursor.Y, 3)
}
