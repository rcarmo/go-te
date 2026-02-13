package te

import "testing"

// From esctest2/esctest/tests/dch.py::test_DCH_DefaultParam
func TestEsctestDchTestDCHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcd")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestDCH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"acd" + esctestEmpty()})
}

// From esctest2/esctest/tests/dch.py::test_DCH_ExplicitParam
func TestEsctestDchTestDCHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcd")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestDCH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"ad" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/dch.py::test_DCH_RespectsMargins
func TestEsctestDchTestDCHRespectsMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcde")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 1})
	esctestDCH(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 1}, []string{"abd" + esctestEmpty() + "e"})
}

// From esctest2/esctest/tests/dch.py::test_DCH_DeleteAllWithMargins
func TestEsctestDchTestDCHDeleteAllWithMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcde")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 1})
	esctestDCH(t, stream, 99)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 1}, []string{"ab" + esctestEmpty() + esctestEmpty() + "e"})
}

// From esctest2/esctest/tests/dch.py::test_DCH_DoesNothingOutsideLeftRightMargin
func TestEsctestDchTestDCHDoesNothingOutsideLeftRightMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcde")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDCH(t, stream, 99)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 1}, []string{"abcde"})
}

// From esctest2/esctest/tests/dch.py::test_DCH_WorksOutsideTopBottomMargin
func TestEsctestDchTestDCHWorksOutsideTopBottomMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcde")
	esctestDECSTBM(t, stream, 2, 3)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDCH(t, stream, 99)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}
