package te

import "testing"

type esctestDecsedFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsedFixture() esctestDecsedFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsedFixture{screen: screen, stream: stream}
}

func (f esctestDecsedFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "a")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "bcd")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 5})
	esctestWrite(t, f.stream, "e")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

func (f esctestDecsedFixture) prepareWide(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcde")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "fghij")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "klmno")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_Default
func TestEsctestDecsedTestDecsedDefault(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"b" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0
func TestEsctestDecsedTestDecsed0(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"b" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1
func TestEsctestDecsedTestDecsed1(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestBlank() + esctestBlank() + "d",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2
func TestEsctestDecsedTestDecsed2(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_3
func TestEsctestDecsedTestDecsed3(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0_WithScrollRegion
func TestEsctestDecsedTestDecsed0WithScrollRegion(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 0)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		"abcde",
		"fg" + esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1_WithScrollRegion
func TestEsctestDecsedTestDecsed1WithScrollRegion(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 1)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestBlank() + esctestBlank() + esctestBlank() + "ij",
		"klmno",
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2_WithScrollRegion
func TestEsctestDecsedTestDecsed2WithScrollRegion(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepareWide(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_Default_Protection
func TestEsctestDecsedTestDecsedDefaultProtection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_DECSCA_2
func TestEsctestDecsedTestDecsedDecsca2(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 2)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0_Protection
func TestEsctestDecsedTestDecsed0Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1_Protection
func TestEsctestDecsedTestDecsed1Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2_Protection
func TestEsctestDecsedTestDecsed2Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_3_Protection
func TestEsctestDecsedTestDecsed3Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"aX" + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0_WithScrollRegion_Protection
func TestEsctestDecsedTestDecsed0WithScrollRegionProtection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepareWide(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 0)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		"abcde",
		"fghij",
		esctestBlank() + "lmno",
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1_WithScrollRegion_Protection
func TestEsctestDecsedTestDecsed1WithScrollRegionProtection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepareWide(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 1)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		esctestBlank() + "bcde",
		"fghij",
		"klmno",
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2_WithScrollRegion_Protection
func TestEsctestDecsedTestDecsed2WithScrollRegionProtection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepareWide(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSED(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 3}, []string{
		esctestBlank() + "bcde",
		"fghij",
		"klmno",
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_doesNotRespectISOProtect
func TestEsctestDecsedTestDecsedDoesNotRespectISOProtect(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"W")
	esctestDECSED(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{esctestBlank() + esctestBlank()})
}
