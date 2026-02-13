package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/sm.py::test_SM_IRM
func TestEsctestSmTestSMIRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestSM(t, stream, esctestModeIRM)
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"Xabc"})
}

// From esctest2/esctest/tests/sm.py::test_SM_IRM_DoesNotWrapUnlessCursorAtMargin
func TestEsctestSmTestSMIRmDoesNotWrapUnlessCursorAtMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestWrite(t, stream, strings.Repeat("a", size.Width-1))
	esctestWrite(t, stream, "b")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestSM(t, stream, esctestModeIRM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
	esctestCUP(t, stream, esctestPoint{X: size.Width, Y: 1})
	esctestWrite(t, stream, "YZ")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"Z"})
}

// From esctest2/esctest/tests/sm.py::test_SM_IRM_TruncatesAtRightMargin
func TestEsctestSmTestSMIRmTruncatesAtRightMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, "abcdef")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 1})
	esctestSM(t, stream, esctestModeIRM)
	esctestWrite(t, stream, "X")
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 1, Right: 11, Bottom: 1}, []string{"abXcde" + esctestEmpty()})
}

func esctestDoLinefeedModeTest(t *testing.T, screen *Screen, stream *Stream, code string) {
	esctestRM(t, stream, esctestModeLNM)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, code)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 2})
	esctestSM(t, stream, esctestModeLNM)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, code)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 2})
}

// From esctest2/esctest/tests/sm.py::test_SM_LNM
func TestEsctestSmTestSMLnm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoLinefeedModeTest(t, screen, stream, ControlLF)
	esctestDoLinefeedModeTest(t, screen, stream, ControlVT)
	esctestDoLinefeedModeTest(t, screen, stream, ControlFF)
}
