package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/ich.py::test_ICH_DefaultParam
func TestEsctestIchTestICHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestICH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 1}, []string{"a" + esctestBlank() + "bcdefg"})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ExplicitParam
func TestEsctestIchTestICHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestICH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 9, Bottom: 1}, []string{"a" + esctestBlank() + esctestBlank() + "bcdefg"})
}

// From esctest2/esctest/tests/ich.py::test_ICH_IsNoOpWhenCursorBeginsOutsideScrollRegion
func TestEsctestIchTestICHIsNoOpWhenCursorBeginsOutsideScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	s := "abcdefg"
	esctestWrite(t, stream, s)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestICH(t, stream, 10)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: len(s), Bottom: 1}, []string{s})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ScrollOffRightEdge
func TestEsctestIchTestICHScrollOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	s := "abcdefg"
	startX := width - len(s) + 1
	esctestCUP(t, stream, esctestPoint{X: startX, Y: 1})
	esctestWrite(t, stream, s)
	esctestCUP(t, stream, esctestPoint{X: startX + 1, Y: 1})
	esctestICH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: startX, Top: 1, Right: width, Bottom: 1}, []string{"a" + esctestBlank() + "bcdef"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ScrollEntirelyOffRightEdge
func TestEsctestIchTestICHScrollEntirelyOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, strings.Repeat("x", width))
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestICH(t, stream, width)
	expectedLine := strings.Repeat(esctestBlank(), width)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: width, Bottom: 1}, []string{expectedLine})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ScrollOffRightMarginInScrollRegion
func TestEsctestIchTestICHScrollOffRightMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	s := "abcdefg"
	esctestWrite(t, stream, s)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 1})
	esctestICH(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: len(s), Bottom: 1}, []string{"ab" + esctestBlank() + "cdfg"})
}
