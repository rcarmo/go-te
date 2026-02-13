package te

import "testing"

type esctestEDFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestEDFixture() esctestEDFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestEDFixture{screen: screen, stream: stream}
}

func (f esctestEDFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "a")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "bcd")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 5})
	esctestWrite(t, f.stream, "e")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

func (f esctestEDFixture) prepareWide(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcde")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "fghij")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "klmno")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/ed.py::test_ED_Default
func TestEsctestEdTestEDDefault(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "b" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_0
func TestEsctestEdTestED0(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "b" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_1
func TestEsctestEdTestED1(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestBlank() + esctestBlank() + "d", esctestEmpty() + esctestEmpty() + esctestEmpty(), "e" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_2
func TestEsctestEdTestED2(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_3
func TestEsctestEdTestED3(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "bcd", esctestEmpty() + esctestEmpty() + esctestEmpty(), "e" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_0_WithScrollRegion
func TestEsctestEdTestED0WithScrollRegion(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestED(t, fixture.stream, 0)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{"abcde", "fg" + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_1_WithScrollRegion
func TestEsctestEdTestED1WithScrollRegion(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestED(t, fixture.stream, 1)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestBlank() + esctestBlank() + esctestBlank() + "ij", "klmno"})
}

// From esctest2/esctest/tests/ed.py::test_ED_2_WithScrollRegion
func TestEsctestEdTestED2WithScrollRegion(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestED(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_doesNotRespectDECProtection
func TestEsctestEdTestEDDoesNotRespectDECProtection(t *testing.T) {
	fixture := newEsctestEDFixture()
	esctestWrite(t, fixture.stream, "a")
	esctestWrite(t, fixture.stream, "b")
	esctestWrite(t, fixture.stream, ControlCSI+"1\"q")
	esctestWrite(t, fixture.stream, "c")
	esctestWrite(t, fixture.stream, ControlCSI+"0\"q")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_respectsISOProtection
func TestEsctestEdTestEDRespectsISOProtection(t *testing.T) {
	fixture := newEsctestEDFixture()
	esctestWrite(t, fixture.stream, "a")
	esctestWrite(t, fixture.stream, "b")
	esctestWrite(t, fixture.stream, ControlESC+"V")
	esctestWrite(t, fixture.stream, "c")
	esctestWrite(t, fixture.stream, ControlESC+"W")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}
