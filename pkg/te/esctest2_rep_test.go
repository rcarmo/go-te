package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/rep.py::test_REP_DefaultParam
func TestEsctestRepTestREPDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestREP(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{"aa" + esctestEmpty()})
}

// From esctest2/esctest/tests/rep.py::test_REP_ExplicitParam
func TestEsctestRepTestREPExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestREP(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"aaa" + esctestEmpty()})
}

// From esctest2/esctest/tests/rep.py::test_REP_RespectsLeftRightMargins
func TestEsctestRepTestREPRespectsLeftRightMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, stream, "a")
	esctestREP(t, stream, 3)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 2}, []string{esctestEmpty() + "aaa" + esctestEmpty(), esctestEmpty() + "a" + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/rep.py::test_REP_RespectsTopBottomMargins
func TestEsctestRepTestREPRespectsTopBottomMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 2, Y: 4})
	esctestWrite(t, stream, "a")
	esctestREP(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 3, Right: size.Width, Bottom: 4}, []string{strings.Repeat(esctestEmpty(), size.Width-3) + "aaa", "a" + strings.Repeat(esctestEmpty(), size.Width-1)})
}
