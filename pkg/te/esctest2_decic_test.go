package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/decic.py::test_DECIC_DefaultParam
func TestEsctestDecicTestDecicDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestDECIC(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 2}, []string{
		"a" + esctestBlank() + "bcdefg",
		"A" + esctestBlank() + "BCDEFG",
	})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_ExplicitParam
func TestEsctestDecicTestDecicExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG"+ControlCR+ControlLF+"zyxwvut")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestDECIC(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 9, Bottom: 3}, []string{
		"a" + strings.Repeat(esctestBlank(), 2) + "bcdefg",
		"A" + strings.Repeat(esctestBlank(), 2) + "BCDEFG",
		"z" + strings.Repeat(esctestBlank(), 2) + "yxwvut",
	})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_CursorWithinTopBottom
func TestEsctestDecicTestDecicCursorWithinTopBottom(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 1, 20)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG"+ControlCR+ControlLF+"zyxwvut"+ControlCR+ControlLF+"ZYXWVUT")
	esctestDECSTBM(t, stream, 2, 3)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestDECIC(t, stream, 2)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 9, Bottom: 4}, []string{
		"abcdefg" + strings.Repeat(esctestEmpty(), 2),
		"A" + strings.Repeat(esctestBlank(), 2) + "BCDEFG",
		"z" + strings.Repeat(esctestBlank(), 2) + "yxwvut",
		"ZYXWVUT" + strings.Repeat(esctestEmpty(), 2),
	})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_IsNoOpWhenCursorBeginsOutsideScrollRegion
func TestEsctestDecicTestDecicIsNoOpWhenCursorBeginsOutsideScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECIC(t, stream, 10)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 2}, []string{
		"abcdefg",
		"ABCDEFG",
	})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_ScrollOffRightEdge
func TestEsctestDecicTestDecicScrollOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	s := "abcdefg"
	startX := width - len(s) + 1
	esctestCUP(t, stream, esctestPoint{X: startX, Y: 1})
	esctestWrite(t, stream, s)
	esctestCUP(t, stream, esctestPoint{X: startX, Y: 2})
	esctestWrite(t, stream, strings.ToUpper(s))
	esctestCUP(t, stream, esctestPoint{X: startX + 1, Y: 1})
	esctestDECIC(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: startX, Top: 1, Right: width, Bottom: 2}, []string{
		"a" + esctestBlank() + "bcdef",
		"A" + esctestBlank() + "BCDEF",
	})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{esctestEmpty(), esctestEmpty()})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_ScrollEntirelyOffRightEdge
func TestEsctestDecicTestDecicScrollEntirelyOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, strings.Repeat("x", width))
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, strings.Repeat("x", width))
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECIC(t, stream, width)
	expectedLine := strings.Repeat(esctestBlank(), width)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: width, Bottom: 2}, []string{expectedLine, expectedLine})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{esctestBlank(), esctestBlank()})
}

// From esctest2/esctest/tests/decic.py::test_DECIC_ScrollOffRightMarginInScrollRegion
func TestEsctestDecicTestDecicScrollOffRightMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	s := "abcdefg"
	esctestWrite(t, stream, s)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, strings.ToUpper(s))
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 1})
	esctestDECIC(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: len(s), Bottom: 2}, []string{
		"ab" + esctestBlank() + "cdfg",
		"AB" + esctestBlank() + "CDFG",
	})
}
