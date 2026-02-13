package te

import "testing"

type esctestDecsetMoreFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsetMoreFixture() esctestDecsetMoreFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsetMoreFixture{screen: screen, stream: stream}
}

func (f esctestDecsetMoreFixture) fillLineAndWriteTab(t *testing.T) {
	esctestWrite(t, f.stream, ControlCR+ControlLF)
	size := esctestGetScreenSize(f.screen)
	for i := 0; i < size.Width; i++ {
		esctestWrite(t, f.stream, "x")
	}
	esctestWrite(t, f.stream, ControlHT)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_OnRespectsLeftRightMargin
func TestEsctestDecsetMoreTestDECSETDECAWMOnRespectsLeftRightMargin(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 5, 9)
	esctestDECSTBM(t, fixture.stream, 5, 9)
	esctestCUP(t, fixture.stream, esctestPoint{X: 8, Y: 9})
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "abcdef")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 5, Top: 8, Right: 9, Bottom: 9}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + "ab", "cdef" + esctestEmpty()})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_OffRespectsLeftRightMargin
func TestEsctestDecsetMoreTestDECSETDECAWMOffRespectsLeftRightMargin(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 5, 9)
	esctestDECSTBM(t, fixture.stream, 5, 9)
	esctestCUP(t, fixture.stream, esctestPoint{X: 8, Y: 9})
	esctestDECRESET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "abcdef")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 9)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 5, Top: 8, Right: 9, Bottom: 9}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + "af"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_Allow80To132
func TestEsctestDecsetMoreTestDECSETAllow80To132(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	if esctestGetScreenSize(fixture.screen).Width == 132 {
		esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
		esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	}
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_TabDoesNotWrapAround
func TestEsctestDecsetMoreTestDECSETDECAWMTabDoesNotWrapAround(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	for i := 0; i < size.Width/8+2; i++ {
		esctestWrite(t, fixture.stream, ControlHT)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 1)
	esctestWrite(t, fixture.stream, "X")
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_NoLineWrapOnTabWithLeftRightMargin
func TestEsctestDecsetMoreTestDECSETDECAWMNoLineWrapOnTabWithLeftRightMargin(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestXtermWinops(t, fixture.stream, 8, 24, 80)
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 10, 20)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 9, Y: 1})
	esctestWrite(t, fixture.stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 17, Y: 1})
	esctestWrite(t, fixture.stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 20, Y: 1})
	esctestWrite(t, fixture.stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 20, Y: 1})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_MoreFix
func TestEsctestDecsetMoreTestDECSETMoreFix(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeMoreFix)
	fixture.fillLineAndWriteTab(t)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 9)
	esctestWrite(t, fixture.stream, "1")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 9, Top: 3, Right: 9, Bottom: 3}, []string{"1"})
	esctestDECRESET(t, fixture.stream, esctestModeMoreFix)
	fixture.fillLineAndWriteTab(t)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, esctestGetScreenSize(fixture.screen).Width)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 5, Right: 1, Bottom: 5}, []string{"2"})
}
