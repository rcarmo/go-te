package te

import (
	"fmt"
	"testing"
)

type esctestSDFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestSDFixture() esctestSDFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestSDFixture{screen: screen, stream: stream}
}

func (f esctestSDFixture) prepare(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/sd.py::test_SD_DefaultParam
func TestEsctestSdTestSDDefaultParam(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestSD(t, fixture.stream)
	expectedLines := []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "abcde", "fghij", "klmno", "pqrst"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_ExplicitParam
func TestEsctestSdTestSDExplicitParam(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestSD(t, fixture.stream, 2)
	expectedLines := []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "abcde", "fghij", "klmno"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_CanClearScreen
func TestEsctestSdTestSDCanClearScreen(t *testing.T) {
	fixture := newEsctestSDFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", y))
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestSD(t, fixture.stream, height)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_RespectsTopBottomScrollRegion
func TestEsctestSdTestSDRespectsTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "fghij", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_OutsideTopBottomScrollRegion
func TestEsctestSdTestSDOutsideTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "fghij", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_RespectsLeftRightScrollRegion
func TestEsctestSdTestSDRespectsLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSD(t, fixture.stream, 2)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"a" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "e", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "kbcdo", "pghit", "ulmny"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_OutsideLeftRightScrollRegion
func TestEsctestSdTestSDOutsideLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"a" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "e", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "kbcdo", "pghit", "ulmny", esctestEmpty() + "qrs" + esctestEmpty(), esctestEmpty() + "vwx" + esctestEmpty()}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 7}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_LeftRightAndTopBottomScrollRegion
func TestEsctestSdTestSDLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "pghit", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_BigScrollLeftRightAndTopBottomScrollRegion
func TestEsctestSdTestSDBigScrollLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestSD(t, fixture.stream, 99)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}
