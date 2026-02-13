package te

import (
	"fmt"
	"testing"
)

type esctestSUFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestSUFixture() esctestSUFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestSUFixture{screen: screen, stream: stream}
}

func (f esctestSUFixture) prepare(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/su.py::test_SU_DefaultParam
func TestEsctestSuTestSUDefaultParam(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestSU(t, fixture.stream)
	expectedLines := []string{"fghij", "klmno", "pqrst", "uvwxy", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_ExplicitParam
func TestEsctestSuTestSUExplicitParam(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestSU(t, fixture.stream, 2)
	expectedLines := []string{"klmno", "pqrst", "uvwxy", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_CanClearScreen
func TestEsctestSuTestSUCanClearScreen(t *testing.T) {
	fixture := newEsctestSUFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", y))
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestSU(t, fixture.stream, height)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_RespectsTopBottomScrollRegion
func TestEsctestSuTestSURespectsTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "pqrst", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_OutsideTopBottomScrollRegion
func TestEsctestSuTestSUOutsideTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "pqrst", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_RespectsLeftRightScrollRegion
func TestEsctestSuTestSURespectsLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSU(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"almne", "fqrsj", "kvwxo", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "u" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "y"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_OutsideLeftRightScrollRegion
func TestEsctestSuTestSUOutsideLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"almne", "fqrsj", "kvwxo", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "u" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "y"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_LeftRightAndTopBottomScrollRegion
func TestEsctestSuTestSULeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "fqrsj", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_BigScrollLeftRightAndTopBottomScrollRegion
func TestEsctestSuTestSUBigScrollLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSU(t, fixture.stream, 99)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}
