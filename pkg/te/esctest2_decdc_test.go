package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/decdc.py::test_DECDC_DefaultParam
func TestEsctestDecdcTestDecdcDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestDECDC(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 2}, []string{
		"acdefg" + esctestEmpty(),
		"ACDEFG" + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_ExplicitParam
func TestEsctestDecdcTestDecdcExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG"+ControlCR+ControlLF+"zyxwvut")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestDECDC(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 3}, []string{
		"adefg" + strings.Repeat(esctestEmpty(), 2),
		"ADEFG" + strings.Repeat(esctestEmpty(), 2),
		"zwvut" + strings.Repeat(esctestEmpty(), 2),
	})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_CursorWithinTopBottom
func TestEsctestDecdcTestDecdcCursorWithinTopBottom(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 1, 20)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG"+ControlCR+ControlLF+"zyxwvut"+ControlCR+ControlLF+"ZYXWVUT")
	esctestDECSTBM(t, stream, 2, 3)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestDECDC(t, stream, 2)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 4}, []string{
		"abcdefg",
		"ADEFG" + strings.Repeat(esctestEmpty(), 2),
		"zwvut" + strings.Repeat(esctestEmpty(), 2),
		"ZYXWVUT",
	})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_IsNoOpWhenCursorBeginsOutsideScrollRegion
func TestEsctestDecdcTestDecdcIsNoOpWhenCursorBeginsOutsideScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefg"+ControlCR+ControlLF+"ABCDEFG")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECDC(t, stream, 10)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 2}, []string{
		"abcdefg",
		"ABCDEFG",
	})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_DeleteAll
func TestEsctestDecdcTestDecdcDeleteAll(t *testing.T) {
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
	esctestDECDC(t, stream, width+10)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: startX, Top: 1, Right: width, Bottom: 2}, []string{
		"a" + strings.Repeat(esctestEmpty(), 6),
		"A" + strings.Repeat(esctestEmpty(), 6),
	})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{esctestEmpty(), esctestEmpty()})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_DeleteWithLeftRightMargins
func TestEsctestDecdcTestDecdcDeleteWithLeftRightMargins(t *testing.T) {
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
	esctestDECDC(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 2}, []string{
		"abde" + esctestEmpty() + "fg",
		"ABDE" + esctestEmpty() + "FG",
	})
}

// From esctest2/esctest/tests/decdc.py::test_DECDC_DeleteAllWithLeftRightMargins
func TestEsctestDecdcTestDecdcDeleteAllWithLeftRightMargins(t *testing.T) {
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
	esctestDECDC(t, stream, 99)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 2}, []string{
		"ab" + strings.Repeat(esctestEmpty(), 3) + "fg",
		"AB" + strings.Repeat(esctestEmpty(), 3) + "FG",
	})
}
