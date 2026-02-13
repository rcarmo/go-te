package te

import (
	"fmt"
	"testing"
)

type esctestDECSTBMFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDECSTBMFixture() esctestDECSTBMFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDECSTBMFixture{screen: screen, stream: stream}
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_ScrollsOnNewline
func TestEsctestDecstbmTestDECSTBMScrollsOnNewline(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"2", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 3)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_NewlineBelowRegion
func TestEsctestDecstbmTestDECSTBMNewlineBelowRegion(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 4})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_MovsCursorToOrigin
func TestEsctestDecstbmTestDECSTBMMovsCursorToOrigin(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_TopBelowBottom
func TestEsctestDecstbmTestDECSTBMTopBelowBottom(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 3, 3)
	for i := 0; i < size.Height; i++ {
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", i))
		y := i + 1
		if y != size.Height {
			esctestWrite(t, fixture.stream, ControlCR+ControlLF)
		}
	}
	for i := 0; i < size.Height; i++ {
		y := i + 1
		esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", i)})
	}
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlLF)
	for i := 0; i < size.Height-1; i++ {
		y := i + 1
		esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", i+1)})
	}
	y := size.Height
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_DefaultRestores
func TestEsctestDecstbmTestDECSTBMDefaultRestores(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	position := esctestGetCursorPosition(fixture.screen)
	esctestDECSTBM(t, fixture.stream)
	esctestCUP(t, fixture.stream, position)
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 4)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_CursorBelowRegionAtBottomTriesToScroll
func TestEsctestDecstbmTestDECSTBMCursorBelowRegionAtBottomTriesToScroll(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, "3"+ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: size.Height, Right: 1, Bottom: size.Height}, []string{"3"})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_MaxSizeOfRegionIsPageSize
func TestEsctestDecstbmTestDECSTBMMaxSizeOfRegionIsPageSize(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "x")
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 1, size.Height+10)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 2}, []string{"x", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_TopOfZeroIsTopOfScreen
func TestEsctestDecstbmTestDECSTBMTopOfZeroIsTopOfScreen(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 0, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "3"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "4")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 3}, []string{"2", "3", "4"})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_BottomOfZeroIsBottomOfScreen
func TestEsctestDecstbmTestDECSTBMBottomOfZeroIsBottomOfScreen(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, fixture.stream, "x")
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 2, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"x", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}
