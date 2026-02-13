package te

import (
	"fmt"
	"testing"
	"time"
)

type esctestXtermWinopsFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestXtermWinopsFixture() esctestXtermWinopsFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestXtermWinopsFixture{screen: screen, stream: stream}
}

func (f esctestXtermWinopsFixture) getIconified(t *testing.T) bool {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 11)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) == 0 {
		t.Fatalf("expected window state response")
	}
	return params[0] == 2
}

func (f esctestXtermWinopsFixture) getWindowPosition(t *testing.T) esctestPoint {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 13)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected window position response")
	}
	return esctestPoint{X: params[1], Y: params[2]}
}

func (f esctestXtermWinopsFixture) getWindowSizePixels(t *testing.T) esctestSize {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 14)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected window size response")
	}
	return esctestSize{Width: params[2], Height: params[1]}
}

func (f esctestXtermWinopsFixture) getScreenSizePixels(t *testing.T) esctestSize {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 15)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected screen size response")
	}
	return esctestSize{Width: params[2], Height: params[1]}
}

func (f esctestXtermWinopsFixture) getCharSizePixels(t *testing.T) esctestSize {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 16)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected char size response")
	}
	return esctestSize{Width: params[2], Height: params[1]}
}

func (f esctestXtermWinopsFixture) getScreenSizeChars(t *testing.T) esctestSize {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 18)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected text area response")
	}
	return esctestSize{Width: params[2], Height: params[1]}
}

func (f esctestXtermWinopsFixture) getDisplaySizeChars(t *testing.T) esctestSize {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 19)
	})
	params := esctestReadCSI(t, response, 't', 0)
	if len(params) < 3 {
		t.Fatalf("expected screen size chars response")
	}
	return esctestSize{Width: params[2], Height: params[1]}
}

func (f esctestXtermWinopsFixture) getIconTitle(t *testing.T) string {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 20)
	})
	payload := esctestReadOSC(t, response, "L")
	return payload
}

func (f esctestXtermWinopsFixture) getWindowTitle(t *testing.T) string {
	response := esctestCaptureResponse(f.screen, func() {
		esctestXtermWinops(t, f.stream, 21)
	})
	payload := esctestReadOSC(t, response, "l")
	return payload
}

func (f esctestXtermWinopsFixture) checkAnySize(t *testing.T, desired, actual, limit esctestSize) {
	error := esctestSize{Width: absInt(actual.Width - desired.Width), Height: absInt(actual.Height - desired.Height)}
	if error.Width > limit.Width {
		t.Fatalf("expected width error <= %d, got %d", limit.Width, error.Width)
	}
	if error.Height > limit.Height {
		t.Fatalf("expected height error <= %d, got %d", limit.Height, error.Height)
	}
}

func (f esctestXtermWinopsFixture) checkActualSizePixels(t *testing.T, desired esctestSize) {
	f.checkAnySize(t, desired, f.getWindowSizePixels(t), f.getPixelErrorLimit(t))
}

func (f esctestXtermWinopsFixture) checkActualSizeChars(t *testing.T, desired, limit esctestSize) {
	f.checkAnySize(t, desired, f.getScreenSizeChars(t), limit)
}

func (f esctestXtermWinopsFixture) checkForShrinkage(t *testing.T, original, actual esctestSize) {
	if actual.Width < original.Width || actual.Height < original.Height {
		t.Fatalf("expected size to not shrink: original=%v actual=%v", original, actual)
	}
}

func (f esctestXtermWinopsFixture) averageWidth(a, b esctestSize) int {
	return (a.Width + b.Width) / 2
}

func (f esctestXtermWinopsFixture) averageHeight(a, b esctestSize) int {
	return (a.Height + b.Height) / 2
}

func (f esctestXtermWinopsFixture) getCharErrorLimit() esctestSize {
	return esctestSize{Width: 0, Height: 0}
}

func (f esctestXtermWinopsFixture) getPixelErrorLimit(t *testing.T) esctestSize {
	frame := f.getFrameSizePixels(t)
	chars := f.getCharSizePixels(t)
	cells := 3
	return esctestSize{Width: frame.Width + cells*chars.Width, Height: frame.Height + cells*chars.Height}
}

func (f esctestXtermWinopsFixture) getFrameSizePixels(t *testing.T) esctestSize {
	inner := f.getScreenSizeChars(t)
	chars := f.getCharSizePixels(t)
	outer := f.getWindowSizePixels(t)
	return esctestSize{Width: outer.Width - inner.Width*chars.Width, Height: outer.Height - inner.Height*chars.Height}
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_IconifyDeiconfiy
func TestEsctestXtermWinopsTestXtermWinopsIconifyDeiconfiy(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	esctestXtermWinops(t, fixture.stream, 2)
	esctestAssertEQ(t, fixture.getIconified(t), true)
	esctestXtermWinops(t, fixture.stream, 1)
	esctestAssertEQ(t, fixture.getIconified(t), false)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_MoveToXY
func TestEsctestXtermWinopsTestXtermWinopsMoveToXY(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	esctestXtermWinops(t, fixture.stream, 3, 0, 0)
	origin := fixture.getWindowPosition(t)
	limit := 10
	for n := 1; n < limit; n++ {
		esctestXtermWinops(t, fixture.stream, 3, origin.X+n, origin.Y+n)
		esctestAssertEQ(t, fixture.getWindowPosition(t), esctestPoint{X: origin.X + n, Y: origin.Y + n})
	}
	for n := limit; n > 1; n-- {
		esctestXtermWinops(t, fixture.stream, 3, origin.X+n, origin.Y+limit)
		esctestAssertEQ(t, fixture.getWindowPosition(t), esctestPoint{X: origin.X + n, Y: origin.Y + limit})
	}
	for n := limit; n > 1; n-- {
		esctestXtermWinops(t, fixture.stream, 3, origin.X+limit, origin.Y+n)
		esctestAssertEQ(t, fixture.getWindowPosition(t), esctestPoint{X: origin.X + limit, Y: origin.Y + n})
	}
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_MoveToXY_Defaults
func TestEsctestXtermWinopsTestXtermWinopsMoveToXYDefaults(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	esctestXtermWinops(t, fixture.stream, 3, 0, 0)
	origin := fixture.getWindowPosition(t)
	limit := 10
	wanted := esctestPoint{X: origin.X + limit, Y: origin.Y + limit}
	esctestXtermWinops(t, fixture.stream, 3, wanted.X, wanted.Y)
	esctestAssertEQ(t, fixture.getWindowPosition(t), wanted)
	esctestXtermWinops(t, fixture.stream, 3, origin.X+limit)
	esctestAssertEQ(t, fixture.getWindowPosition(t), esctestPoint{X: origin.X + limit, Y: origin.Y})
	esctestXtermWinops(t, fixture.stream, 3, origin.X, origin.Y+limit)
	esctestAssertEQ(t, fixture.getWindowPosition(t), esctestPoint{X: origin.X, Y: origin.Y + limit})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizePixels_BothParameters
func TestEsctestXtermWinopsTestXtermWinopsResizePixelsBothParameters(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getScreenSizePixels(t)
	originalSize := fixture.getWindowSizePixels(t)
	desiredSize := esctestSize{Width: 400, Height: 200}
	if maximumSize.Width > 0 && maximumSize.Height > 0 {
		desiredSize = esctestSize{Width: fixture.averageWidth(maximumSize, originalSize), Height: fixture.averageHeight(maximumSize, originalSize)}
	}
	esctestXtermWinops(t, fixture.stream, 4, desiredSize.Height, desiredSize.Width)
	fixture.checkActualSizePixels(t, desiredSize)
	esctestXtermWinops(t, fixture.stream, 4, originalSize.Height, originalSize.Width)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizePixels_OmittedHeight
func TestEsctestXtermWinopsTestXtermWinopsResizePixelsOmittedHeight(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getScreenSizePixels(t)
	originalSize := fixture.getWindowSizePixels(t)
	desiredSize := esctestSize{Width: 400, Height: originalSize.Height}
	if maximumSize.Width > 0 {
		desiredSize = esctestSize{Width: maximumSize.Width, Height: originalSize.Height}
	}
	esctestXtermWinops(t, fixture.stream, 4, desiredSize.Height, desiredSize.Width)
	fixture.checkActualSizePixels(t, desiredSize)
	esctestXtermWinops(t, fixture.stream, 4, originalSize.Height, originalSize.Width)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizePixels_OmittedWidth
func TestEsctestXtermWinopsTestXtermWinopsResizePixelsOmittedWidth(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getScreenSizePixels(t)
	originalSize := fixture.getWindowSizePixels(t)
	desiredSize := esctestSize{Width: originalSize.Width, Height: 200}
	if maximumSize.Height > 0 {
		desiredSize = esctestSize{Width: originalSize.Width, Height: fixture.averageHeight(maximumSize, originalSize)}
	}
	esctestXtermWinops(t, fixture.stream, 4, desiredSize.Height)
	fixture.checkActualSizePixels(t, desiredSize)
	esctestXtermWinops(t, fixture.stream, 4, originalSize.Height, originalSize.Width)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizePixels_ZeroWidth
func TestEsctestXtermWinopsTestXtermWinopsResizePixelsZeroWidth(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getScreenSizePixels(t)
	originalSize := fixture.getWindowSizePixels(t)
	desiredHeight := 200
	if maximumSize.Height > 0 {
		desiredHeight = fixture.averageHeight(maximumSize, originalSize)
	}
	esctestXtermWinops(t, fixture.stream, 4, desiredHeight, 0)
	actualSize := fixture.getWindowSizePixels(t)
	maxError := 20
	if absInt(actualSize.Height-desiredHeight) > maxError {
		t.Fatalf("expected height error <= %d, got %d", maxError, absInt(actualSize.Height-desiredHeight))
	}
	displaySize := fixture.getDisplaySizeChars(t)
	screenSize := fixture.getScreenSizeChars(t)
	maxError = 5
	if absInt(displaySize.Width-screenSize.Width) >= maxError {
		t.Fatalf("expected display width close to screen width")
	}
	esctestXtermWinops(t, fixture.stream, 4, originalSize.Height, originalSize.Width)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizePixels_ZeroHeight
func TestEsctestXtermWinopsTestXtermWinopsResizePixelsZeroHeight(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getScreenSizePixels(t)
	originalSize := fixture.getWindowSizePixels(t)
	desiredWidth := 400
	if maximumSize.Width > 0 {
		desiredWidth = fixture.averageWidth(maximumSize, originalSize)
	}
	esctestXtermWinops(t, fixture.stream, 4, 0, desiredWidth)
	actualSize := fixture.getWindowSizePixels(t)
	maxError := 20
	if absInt(actualSize.Width-desiredWidth) > maxError {
		t.Fatalf("expected width error <= %d, got %d", maxError, absInt(actualSize.Width-desiredWidth))
	}
	displaySize := fixture.getDisplaySizeChars(t)
	screenSize := fixture.getScreenSizeChars(t)
	maxError = 5
	if absInt(displaySize.Height-screenSize.Height) >= maxError {
		t.Fatalf("expected display height close to screen height")
	}
	esctestXtermWinops(t, fixture.stream, 4, originalSize.Height, originalSize.Width)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizeChars_BothParameters
func TestEsctestXtermWinopsTestXtermWinopsResizeCharsBothParameters(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getDisplaySizeChars(t)
	originalSize := fixture.getScreenSizeChars(t)
	desiredSize := esctestSize{Width: 20, Height: 21}
	if maximumSize.Width > 0 && maximumSize.Height > 0 {
		desiredSize = esctestSize{Width: fixture.averageWidth(maximumSize, originalSize), Height: fixture.averageHeight(maximumSize, originalSize)}
	}
	esctestXtermWinops(t, fixture.stream, 8, desiredSize.Height, desiredSize.Width)
	fixture.checkActualSizeChars(t, desiredSize, fixture.getCharErrorLimit())
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizeChars_ZeroWidth
func TestEsctestXtermWinopsTestXtermWinopsResizeCharsZeroWidth(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getDisplaySizeChars(t)
	originalSize := fixture.getScreenSizeChars(t)
	desiredSize := esctestSize{Width: maximumSize.Width, Height: originalSize.Height}
	esctestXtermWinops(t, fixture.stream, 8, desiredSize.Height, 0)
	limit := fixture.getCharErrorLimit()
	limit = esctestSize{Width: limit.Width, Height: 0}
	fixture.checkActualSizeChars(t, desiredSize, limit)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizeChars_ZeroHeight
func TestEsctestXtermWinopsTestXtermWinopsResizeCharsZeroHeight(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	maximumSize := fixture.getDisplaySizeChars(t)
	originalSize := fixture.getScreenSizeChars(t)
	desiredSize := esctestSize{Width: originalSize.Width, Height: maximumSize.Height}
	esctestXtermWinops(t, fixture.stream, 8, 0, desiredSize.Width)
	limit := fixture.getCharErrorLimit()
	limit = esctestSize{Width: 0, Height: limit.Height}
	fixture.checkActualSizeChars(t, desiredSize, limit)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizeChars_DefaultWidth
func TestEsctestXtermWinopsTestXtermWinopsResizeCharsDefaultWidth(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	originalSize := fixture.getScreenSizeChars(t)
	displaySize := fixture.getDisplaySizeChars(t)
	desiredSize := esctestSize{Width: originalSize.Width, Height: 21}
	if displaySize.Height > 0 {
		desiredSize = esctestSize{Width: originalSize.Width, Height: fixture.averageHeight(originalSize, displaySize)}
	}
	esctestXtermWinops(t, fixture.stream, 8, desiredSize.Height)
	fixture.checkActualSizeChars(t, desiredSize, esctestSize{Width: 0, Height: 20})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ResizeChars_DefaultHeight
func TestEsctestXtermWinopsTestXtermWinopsResizeCharsDefaultHeight(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	originalSize := fixture.getScreenSizeChars(t)
	displaySize := fixture.getDisplaySizeChars(t)
	desiredSize := esctestSize{Width: 20, Height: originalSize.Height}
	if displaySize.Width > 0 {
		desiredSize = esctestSize{Width: fixture.averageWidth(originalSize, displaySize), Height: originalSize.Height}
	}
	esctestXtermWinops(t, fixture.stream, 8, desiredSize.Height, desiredSize.Width)
	fixture.checkActualSizeChars(t, desiredSize, esctestSize{Width: 0, Height: 0})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_MaximizeWindow_HorizontallyAndVertically
func TestEsctestXtermWinopsTestXtermWinopsMaximizeWindowHorizontallyAndVertically(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	esctestXtermWinops(t, fixture.stream, 9, 1)
	actualSize := fixture.getScreenSizeChars(t)
	desiredSize := fixture.getDisplaySizeChars(t)
	esctestXtermWinops(t, fixture.stream, 9, 0)
	fixture.checkAnySize(t, desiredSize, actualSize, esctestSize{Width: 3, Height: 3})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_MaximizeWindow_Horizontally
func TestEsctestXtermWinopsTestXtermWinopsMaximizeWindowHorizontally(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	desiredSize := esctestSize{Width: fixture.getDisplaySizeChars(t).Width, Height: fixture.getScreenSizeChars(t).Height}
	esctestXtermWinops(t, fixture.stream, 9, 3)
	actualSize := fixture.getScreenSizeChars(t)
	esctestXtermWinops(t, fixture.stream, 9, 0)
	fixture.checkAnySize(t, desiredSize, actualSize, esctestSize{Width: 3, Height: 0})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_MaximizeWindow_Vertically
func TestEsctestXtermWinopsTestXtermWinopsMaximizeWindowVertically(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	desiredSize := esctestSize{Width: fixture.getScreenSizeChars(t).Width, Height: fixture.getDisplaySizeChars(t).Height}
	esctestXtermWinops(t, fixture.stream, 9, 2)
	actualSize := fixture.getScreenSizeChars(t)
	esctestXtermWinops(t, fixture.stream, 9, 0)
	fixture.checkAnySize(t, desiredSize, actualSize, esctestSize{Width: 0, Height: 5})
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_Fullscreen
func TestEsctestXtermWinopsTestXtermWinopsFullscreen(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	originalSize := fixture.getScreenSizeChars(t)
	displaySize := fixture.getDisplaySizeChars(t)
	esctestXtermWinops(t, fixture.stream, 10, 1)
	actualSize := fixture.getScreenSizeChars(t)
	fixture.checkAnySize(t, displaySize, actualSize, esctestSize{Width: 3, Height: 3})
	fixture.checkForShrinkage(t, originalSize, actualSize)
	esctestXtermWinops(t, fixture.stream, 10, 0)
	fixture.checkForShrinkage(t, originalSize, fixture.getScreenSizeChars(t))
	esctestXtermWinops(t, fixture.stream, 10, 2)
	fixture.checkForShrinkage(t, originalSize, fixture.getScreenSizeChars(t))
	esctestXtermWinops(t, fixture.stream, 10, 2)
	fixture.checkForShrinkage(t, originalSize, fixture.getScreenSizeChars(t))
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ReportIconLabel
func TestEsctestXtermWinopsTestXtermWinopsReportIconLabel(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("test %d", time.Now().UnixNano())
	esctestChangeIconTitle(t, fixture.stream, value)
	esctestAssertEQ(t, fixture.getIconTitle(t), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_ReportWindowLabel
func TestEsctestXtermWinopsTestXtermWinopsReportWindowLabel(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("test %d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, value)
	esctestAssertEQ(t, fixture.getWindowTitle(t), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushIconAndWindow_PopIconAndWindow
func TestEsctestXtermWinopsTestXtermWinopsPushIconAndWindowPopIconAndWindow(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, value)
	esctestChangeIconTitle(t, fixture.stream, value)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
	esctestXtermWinops(t, fixture.stream, 22, 0)
	esctestChangeWindowTitle(t, fixture.stream, "x")
	esctestChangeIconTitle(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
	esctestXtermWinops(t, fixture.stream, 23, 0)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushIconAndWindow_PopIcon
func TestEsctestXtermWinopsTestXtermWinopsPushIconAndWindowPopIcon(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, value)
	esctestChangeIconTitle(t, fixture.stream, value)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
	esctestXtermWinops(t, fixture.stream, 22, 0)
	esctestChangeWindowTitle(t, fixture.stream, "x")
	esctestChangeIconTitle(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
	esctestXtermWinops(t, fixture.stream, 23, 1)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
	esctestXtermWinops(t, fixture.stream, 23, 2)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushIconAndWindow_PopWindow
func TestEsctestXtermWinopsTestXtermWinopsPushIconAndWindowPopWindow(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, value)
	esctestChangeIconTitle(t, fixture.stream, value)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
	esctestXtermWinops(t, fixture.stream, 22, 0)
	esctestChangeWindowTitle(t, fixture.stream, "x")
	esctestChangeIconTitle(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
	esctestXtermWinops(t, fixture.stream, 23, 2)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
	esctestXtermWinops(t, fixture.stream, 23, 1)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushIcon_PopIcon
func TestEsctestXtermWinopsTestXtermWinopsPushIconPopIcon(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, "x")
	esctestChangeIconTitle(t, fixture.stream, value)
	esctestXtermWinops(t, fixture.stream, 22, 1)
	esctestChangeIconTitle(t, fixture.stream, "y")
	esctestXtermWinops(t, fixture.stream, 23, 1)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushWindow_PopWindow
func TestEsctestXtermWinopsTestXtermWinopsPushWindowPopWindow(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	esctestChangeIconTitle(t, fixture.stream, "x")
	esctestChangeWindowTitle(t, fixture.stream, value)
	esctestXtermWinops(t, fixture.stream, 22, 2)
	esctestChangeWindowTitle(t, fixture.stream, "y")
	esctestXtermWinops(t, fixture.stream, 23, 2)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), "x")
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushIconThenWindowThenPopBoth
func TestEsctestXtermWinopsTestXtermWinopsPushIconThenWindowThenPopBoth(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value1 := fmt.Sprintf("a%d", time.Now().UnixNano())
	value2 := fmt.Sprintf("b%d", time.Now().UnixNano())
	esctestChangeWindowTitle(t, fixture.stream, value1)
	esctestChangeIconTitle(t, fixture.stream, value2)
	esctestXtermWinops(t, fixture.stream, 22, 1)
	esctestXtermWinops(t, fixture.stream, 22, 2)
	esctestChangeWindowTitle(t, fixture.stream, "y")
	esctestChangeIconTitle(t, fixture.stream, "z")
	esctestXtermWinops(t, fixture.stream, 23, 0)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value1)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value2)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushMultiplePopMultiple_Icon
func TestEsctestXtermWinopsTestXtermWinopsPushMultiplePopMultipleIcon(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value1 := fmt.Sprintf("a%d", time.Now().UnixNano())
	value2 := fmt.Sprintf("b%d", time.Now().UnixNano())
	for _, title := range []string{value1, value2} {
		esctestChangeIconTitle(t, fixture.stream, title)
		esctestXtermWinops(t, fixture.stream, 22, 1)
	}
	esctestChangeIconTitle(t, fixture.stream, "z")
	esctestXtermWinops(t, fixture.stream, 23, 1)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value2)
	esctestXtermWinops(t, fixture.stream, 23, 1)
	esctestAssertEQ(t, esctestGetIconTitle(fixture.screen), value1)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_PushMultiplePopMultiple_Window
func TestEsctestXtermWinopsTestXtermWinopsPushMultiplePopMultipleWindow(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	value1 := fmt.Sprintf("a%d", time.Now().UnixNano())
	value2 := fmt.Sprintf("b%d", time.Now().UnixNano())
	for _, title := range []string{value1, value2} {
		esctestChangeWindowTitle(t, fixture.stream, title)
		esctestXtermWinops(t, fixture.stream, 22, 2)
	}
	esctestChangeWindowTitle(t, fixture.stream, "z")
	esctestXtermWinops(t, fixture.stream, 23, 2)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value2)
	esctestXtermWinops(t, fixture.stream, 23, 2)
	esctestAssertEQ(t, esctestGetWindowTitle(fixture.screen), value1)
}

// From esctest2/esctest/tests/xterm_winops.py::test_XtermWinops_DECSLPP
func TestEsctestXtermWinopsTestXtermWinopsDecslpp(t *testing.T) {
	fixture := newEsctestXtermWinopsFixture()
	esctestXtermWinops(t, fixture.stream, 8, 10, 90)
	esctestAssertEQ(t, fixture.getScreenSizeChars(t), esctestSize{Width: 90, Height: 10})
	esctestXtermWinops(t, fixture.stream, 8, 24)
	esctestAssertEQ(t, fixture.getScreenSizeChars(t), esctestSize{Width: 90, Height: 24})
	esctestXtermWinops(t, fixture.stream, 8, 30)
	esctestAssertEQ(t, fixture.getScreenSizeChars(t), esctestSize{Width: 90, Height: 30})
}
