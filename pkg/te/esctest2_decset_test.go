package te

import "testing"

type esctestDecsetFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsetFixture() esctestDecsetFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsetFixture{screen: screen, stream: stream}
}

func (f esctestDecsetFixture) doAltBufTest(t *testing.T, code int, altClearsBeforeMain bool, cursorSaved bool, movesCursorOnEnter bool) {
	esctestWrite(t, f.stream, "abc"+ControlCR+ControlLF+"abc")
	var mainCursorPosition esctestPoint
	if cursorSaved {
		mainCursorPosition = esctestGetCursorPosition(f.screen)
	}
	before := esctestGetCursorPosition(f.screen)
	esctestDECSET(t, f.stream, code)
	after := esctestGetCursorPosition(f.screen)
	if !movesCursorOnEnter {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	esctestED(t, f.stream, 2)
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "def"+ControlCR+ControlLF+"def")
	esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), "def", "def"})
	before = esctestGetCursorPosition(f.screen)
	esctestDECRESET(t, f.stream, code)
	after = esctestGetCursorPosition(f.screen)
	if cursorSaved {
		esctestAssertEQ(t, mainCursorPosition.X, after.X)
		esctestAssertEQ(t, mainCursorPosition.Y, after.Y)
	} else {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{"abc", "abc", esctestEmpty() + esctestEmpty() + esctestEmpty()})
	before = esctestGetCursorPosition(f.screen)
	esctestDECSET(t, f.stream, code)
	after = esctestGetCursorPosition(f.screen)
	if !movesCursorOnEnter {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	if altClearsBeforeMain {
		esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
	} else {
		esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), "def", "def"})
	}
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM
func TestEsctestDecsetTestDECSETDECAWM(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "abc")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"c"})
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestDECRESET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "ABC")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: size.Width - 1, Top: 1, Right: size.Width, Bottom: 1}, []string{"AC"})
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"c"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_CursorAtRightMargin
func TestEsctestDecsetTestDECSETDECAWMCursorAtRightMargin(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
	for i := 0; i < size.Width-2; i++ {
		esctestWrite(t, fixture.stream, "x")
	}
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width-1)
	esctestWrite(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_BS
func TestEsctestDecsetTestDECSETReverseWraparoundBS(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, esctestGetScreenSize(fixture.screen).Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparoundLastCol_BS
func TestEsctestDecsetTestDECSETReverseWraparoundLastColBS(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, fixture.stream, "a")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, "b")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_Multi
func TestEsctestDecsetTestDECSETReverseWraparoundMulti(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, fixture.stream, "abcd")
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUB(t, fixture.stream, 4)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width-1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ResetReverseWraparoundDisablesIt
func TestEsctestDecsetTestDECSETResetReverseWraparoundDisablesIt(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECRESET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_RequiresDECAWM
func TestEsctestDecsetTestDECSETReverseWraparoundRequiresDECAWM(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECRESET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_OPT_ALTBUF
func TestEsctestDecsetTestDECSETOptAltBuf(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeOptAltBuf, true, false, false)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ALTBUF
func TestEsctestDecsetTestDECSETAltBuf(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeAltBuf, false, false, false)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_OPT_ALTBUF_CURSOR
func TestEsctestDecsetTestDECSETOptAltBufCursor(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeOptAltBufCursor, true, true, true)
}
