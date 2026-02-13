package te

import (
	"strings"
	"testing"
)

type esctestDecselFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecselFixture() esctestDecselFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecselFixture{screen: screen, stream: stream}
}

func (f esctestDecselFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcdefghij")
	esctestCUP(t, f.stream, esctestPoint{X: 5, Y: 1})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_Default
func TestEsctestDecselTestDecselDefault(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + strings.Repeat(esctestEmpty(), 6)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_0
func TestEsctestDecselTestDecsel0(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + strings.Repeat(esctestEmpty(), 6)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_1
func TestEsctestDecselTestDecsel1(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{strings.Repeat(esctestBlank(), 5) + "fghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_2
func TestEsctestDecselTestDecsel2(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{strings.Repeat(esctestEmpty(), 10)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_IgnoresScrollRegion
func TestEsctestDecselTestDecselIgnoresScrollRegion(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{strings.Repeat(esctestEmpty(), 10)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_Default_Protection
func TestEsctestDecselTestDecselDefaultProtection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 10, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcdefghi" + esctestEmpty()})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_0_Protection
func TestEsctestDecselTestDecsel0Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 10, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcdefghi" + esctestEmpty()})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_1_Protection
func TestEsctestDecselTestDecsel1Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + "bcdefghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_2_Protection
func TestEsctestDecselTestDecsel2Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSEL(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + "bcdefghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_IgnoresScrollRegion_Protection
func TestEsctestDecselTestDecselIgnoresScrollRegionProtection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + "bcdefghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_doesNotRespectISOProtect
func TestEsctestDecselTestDecselDoesNotRespectISOProtect(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"W")
	esctestDECSEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{strings.Repeat(esctestBlank(), 2)})
}
