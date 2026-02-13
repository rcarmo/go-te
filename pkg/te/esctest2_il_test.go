package te

import (
	"fmt"
	"testing"
)

type esctestILFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestILFixture() esctestILFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestILFixture{screen: screen, stream: stream}
}

func (f esctestILFixture) prepareWide(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcde")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "fghij")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "klmno")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

func (f esctestILFixture) prepareRegion(t *testing.T) {
	lines := []string{"abcde", "fGHIj", "kLMNo", "pQRSt", "uvwxy"}
	for i, line := range lines {
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: i + 1})
		esctestWrite(t, f.stream, line)
	}
	esctestDECSET(t, f.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, f.stream, 2, 4)
	esctestDECSTBM(t, f.stream, 2, 4)
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/il.py::test_IL_DefaultParam
func TestEsctestIlTestILDefaultParam(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareWide(t)
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 4}, []string{"abcde", "fghij", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "klmno"})
}

// From esctest2/esctest/tests/il.py::test_IL_ExplicitParam
func TestEsctestIlTestILExplicitParam(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareWide(t)
	esctestWrite(t, fixture.stream, ControlCSI+"2"+EscIL)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, []string{"abcde", "fghij", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "klmno"})
}

// From esctest2/esctest/tests/il.py::test_IL_ScrollsOffBottom
func TestEsctestIlTestILScrollsOffBottom(t *testing.T) {
	fixture := newEsctestILFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	for i := 0; i < height; i++ {
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: i + 1})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", i+1))
	}
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)

	expected := 1
	for i := 0; i < height; i++ {
		y := i + 1
		if y == 2 {
			esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
		} else {
			esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", expected)})
			expected++
		}
	}
}

// From esctest2/esctest/tests/il.py::test_IL_RespectsScrollRegion
func TestEsctestIlTestILRespectsScrollRegion(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareRegion(t)
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "kGHIo", "pLMNt", "uvwxy"})
}

// From esctest2/esctest/tests/il.py::test_IL_RespectsScrollRegion_Over
func TestEsctestIlTestILRespectsScrollRegionOver(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareRegion(t)
	esctestWrite(t, fixture.stream, ControlCSI+"99"+EscIL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"})
}

// From esctest2/esctest/tests/il.py::test_IL_AboveScrollRegion
func TestEsctestIlTestILAboveScrollRegion(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareRegion(t)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, []string{"abcde", "fGHIj", "kLMNo", "pQRSt", "uvwxy"})
}
