package te

import (
	"strings"
	"testing"
)

// From esctest2/esctest/tests/cub.py::test_CUB_DefaultParam
func TestEsctestCubTestCUBDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUB(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 4, Y: 3})
}

// From esctest2/esctest/tests/cub.py::test_CUB_ExplicitParam
func TestEsctestCubTestCUBExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 4})
	esctestCUB(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 3)
}

// From esctest2/esctest/tests/cub.py::test_CUB_StopsAtLeftEdge
func TestEsctestCubTestCUBStopsAtLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUB(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/cub.py::test_CUB_StopsAtLeftEdgeWhenBegunLeftOfScrollRegion
func TestEsctestCubTestCUBStopsAtLeftEdgeWhenBegunLeftOfScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 4, Y: 3})
	esctestCUB(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/cub.py::test_CUB_StopsAtLeftMarginInScrollRegion
func TestEsctestCubTestCUBStopsAtLeftMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 3})
	esctestCUB(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 5)
}

// From esctest2/esctest/tests/cub.py::test_CUB_AfterNoWrappedInlines
func TestEsctestCubTestCUBAfterNoWrappedInlines(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", size.Width-2) + "\n"
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill)
	esctestWrite(t, stream, fill)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestCUB(t, stream, size.Width*2)
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 4})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 3})
	}
}

// From esctest2/esctest/tests/cub.py::test_CUB_AfterOneWrappedInline
func TestEsctestCubTestCUBAfterOneWrappedInline(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", (size.Width+2)*2)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill+"\n"+fill)
	esctestCUB(t, stream, size.Width*5)
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 6})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 9, Y: 3})
	}
}
