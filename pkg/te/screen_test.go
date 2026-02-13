package te

import (
	"testing"
)

// From pyte/tests/test_screen.py::test_attributes
func TestAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{1}, false)
	if !screen.Cursor.Attr.Bold {
		t.Fatalf("expected bold cursor attr")
	}
	screen.Draw("f")
	cell := screen.Buffer[0][0]
	if cell.Data != "f" || !cell.Attr.Bold {
		t.Fatalf("expected bold cell")
	}
}

// From pyte/tests/test_screen.py::test_blink
func TestBlink(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{5}, false)
	screen.Draw("f")
	cell := screen.Buffer[0][0]
	if !cell.Attr.Blink {
		t.Fatalf("expected blink cell")
	}
}

// From pyte/tests/test_screen.py::test_colors
func TestColors(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{30}, false)
	screen.SelectGraphicRendition([]int{40}, false)
	if screen.Cursor.Attr.Fg.Name != "black" || screen.Cursor.Attr.Bg.Name != "black" {
		t.Fatalf("expected black colors")
	}
	screen.SelectGraphicRendition([]int{31}, false)
	if screen.Cursor.Attr.Fg.Name != "red" {
		t.Fatalf("expected red foreground")
	}
}

// From pyte/tests/test_screen.py::test_colors256
func TestColors256(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{SgrFg256, 5, 0}, false)
	screen.SelectGraphicRendition([]int{SgrBg256, 5, 15}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 256 colors")
	}
}

// From pyte/tests/test_screen.py::test_colors24bit
func TestColors24Bit(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{38, 2, 0, 0, 0}, false)
	screen.SelectGraphicRendition([]int{48, 2, 255, 255, 255}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 24-bit colors")
	}
}

// From pyte/tests/test_screen.py::test_colors_aixterm
func TestColorsAixterm(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{94}, false)
	if screen.Cursor.Attr.Fg.Name != "brightblue" {
		t.Fatalf("expected brightblue fg")
	}
	screen.SelectGraphicRendition([]int{104}, false)
	if screen.Cursor.Attr.Bg.Name != "brightblue" {
		t.Fatalf("expected brightblue bg")
	}
}

// From pyte/tests/test_screen.py::test_reset_resets_colors
func TestResetResetsColors(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{30}, false)
	screen.SelectGraphicRendition([]int{40}, false)
	screen.SelectGraphicRendition([]int{0}, false)
	if screen.Cursor.Attr.Fg.Name != "default" || screen.Cursor.Attr.Bg.Name != "default" {
		t.Fatalf("expected reset colors")
	}
}

func TestResetBetweenAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{31, 0, 41}, false)
	if screen.Cursor.Attr.Fg.Name != "default" || screen.Cursor.Attr.Bg.Name != "red" {
		t.Fatalf("expected reset fg and red bg")
	}
}

// From pyte/tests/test_screen.py::test_multi_attribs
func TestMultiAttribs(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{1, 3}, false)
	if !screen.Cursor.Attr.Bold || !screen.Cursor.Attr.Italics {
		t.Fatalf("expected bold and italics")
	}
}

// From pyte/tests/test_screen.py::test_attributes_reset
func TestAttributesReset(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SetMode([]int{ModeLNM}, false)
	screen.SelectGraphicRendition([]int{1}, false)
	screen.Draw("foo")
	screen.CursorPosition(0, 0)
	screen.SelectGraphicRendition([]int{0}, false)
	screen.Draw("f")
	if screen.Buffer[0][0].Attr.Bold {
		t.Fatalf("expected reset attr")
	}
}

// From pyte/tests/test_screen.py::test_resize
func TestResize(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.SetMargins(0, 1)
	screen.Resize(3, 3)
	if screen.Columns != 3 || screen.Lines != 3 {
		t.Fatalf("expected resized 3x3")
	}
	screen.Resize(2, 2)
	if screen.Columns != 2 || screen.Lines != 2 {
		t.Fatalf("expected resized 2x2")
	}
}

func TestDrawAutowrapAndIRM(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.SetMode([]int{ModeLNM}, false)
	for _, ch := range "abc" {
		screen.Draw(string(ch))
	}
	if screen.Display()[0] != "abc" {
		t.Fatalf("expected abc")
	}
	screen.Draw("a")
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("expected wrapped cursor")
	}

	screen = NewScreen(3, 3)
	screen.ResetMode([]int{ModeDECAWM}, false)
	for _, ch := range "abc" {
		screen.Draw(string(ch))
	}
	screen.Draw("a")
	if screen.Display()[0] != "aba" {
		t.Fatalf("expected overwrite")
	}

	screen.SetMode([]int{ModeIRM}, false)
	screen.CursorPosition(0, 0)
	screen.Draw("x")
	if screen.Display()[0] != "xab" {
		t.Fatalf("expected insert mode")
	}
}

func TestDrawUTF8(t *testing.T) {
	screen := NewScreen(1, 1)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("\xE2\x80\x9D")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "”" {
		t.Fatalf("expected utf8 char")
	}
}

func TestDrawWidthCombining(t *testing.T) {
	screen := NewScreen(4, 2)
	screen.Draw("\u0308")
	screen.Draw("bad")
	screen.Draw("\u0308")
	if screen.Display()[0] != "bad̈ " {
		t.Fatalf("expected combining")
	}
	screen.Draw("!")
	screen.Draw("\u0308")
	if screen.Display()[0] != "bad̈!̈" {
		t.Fatalf("expected combining previous line")
	}
}

func TestCarriageReturnAndIndex(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.Cursor.Col = 2
	screen.CarriageReturn()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected carriage return")
	}

	screen = NewScreen(2, 2)
	screen.Draw("wo")
	screen.Cursor.Row = 1
	screen.Index()
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected cursor row")
	}
}

func TestTabStops(t *testing.T) {
	screen := NewScreen(10, 10)
	if _, ok := screen.TabStops[8]; !ok {
		t.Fatalf("expected default tabstop")
	}
	screen.ClearTabStop(3)
	screen.Cursor.Col = 1
	screen.SetTabStop()
	screen.Cursor.Col = 8
	screen.SetTabStop()
	screen.Cursor.Col = 0
	screen.Tab()
	if screen.Cursor.Col != 1 {
		t.Fatalf("expected tab to col1")
	}
}

func TestHideCursorMode(t *testing.T) {
	screen := NewScreen(10, 10)
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
	screen.ResetMode([]int{ModeDECTCEM}, false)
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
	screen.SetMode([]int{ModeDECTCEM}, false)
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
}
