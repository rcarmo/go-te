package te

import (
	"fmt"
	"testing"
)

type esctestDLFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDLFixture() esctestDLFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDLFixture{screen: screen, stream: stream}
}

func (f esctestDLFixture) prepare(t *testing.T) {
	height := esctestGetScreenSize(f.screen).Height
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, fmt.Sprintf("%04d", y))
	}
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
}

func (f esctestDLFixture) prepareForRegion(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/dl.py::test_DL_DefaultParam
func TestEsctestDlTestDLDefaultParam(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	height := esctestGetScreenSize(fixture.screen).Height
	y := 1
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		if y != 2 {
			expectedLines = append(expectedLines, fmt.Sprintf("%04d", y))
		}
		y++
	}
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_ExplicitParam
func TestEsctestDlTestDLExplicitParam(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	esctestWrite(t, fixture.stream, ControlCSI+"2"+EscDL)
	height := esctestGetScreenSize(fixture.screen).Height
	y := 1
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		if y < 2 || y > 3 {
			expectedLines = append(expectedLines, fmt.Sprintf("%04d", y))
		}
		y++
	}
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_DeleteMoreThanVisible
func TestEsctestDlTestDLDeleteMoreThanVisible(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	height := esctestGetScreenSize(fixture.screen).Height
	esctestWrite(t, fixture.stream, ControlCSI+fmt.Sprintf("%d", height*2)+EscDL)
	expectedLines := []string{"0001"}
	for i := 0; i < height-1; i++ {
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_InScrollRegion
func TestEsctestDlTestDLInScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "klmno", "pqrst", esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_OutsideScrollRegion
func TestEsctestDlTestDLOutsideScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 1})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_InLeftRightScrollRegion
func TestEsctestDlTestDLInLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "flmnj", "kqrso", "pvwxt", "u" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "y"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_OutsideLeftRightScrollRegion
func TestEsctestDlTestDLOutsideLeftRightScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	expectedLines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_InLeftRightAndTopBottomScrollRegion
func TestEsctestDlTestDLInLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "flmnj", "kqrso", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_ClearOutLeftRightAndTopBottomScrollRegion
func TestEsctestDlTestDLClearOutLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+"99"+EscDL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "f" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "j", "k" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "o", "p" + esctestEmpty() + esctestEmpty() + esctestEmpty() + "t", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_OutsideLeftRightAndTopBottomScrollRegion
func TestEsctestDlTestDLOutsideLeftRightAndTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 4)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}
