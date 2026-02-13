package te

import (
	"testing"
)

// From pyte/tests/test_screen.py::test_initialize_char
func TestPyteTestScreenInitializeChar(t *testing.T) {
	fields := map[string]func(*Cell){
		"fg": func(c *Cell) { c.Attr.Fg = Color{Name: "true", Mode: ColorANSI16} },
		"bg": func(c *Cell) { c.Attr.Bg = Color{Name: "true", Mode: ColorANSI16} },
		"bold": func(c *Cell) { c.Attr.Bold = true },
		"italics": func(c *Cell) { c.Attr.Italics = true },
		"underscore": func(c *Cell) { c.Attr.Underline = true },
		"strikethrough": func(c *Cell) { c.Attr.Strikethrough = true },
		"reverse": func(c *Cell) { c.Attr.Reverse = true },
		"blink": func(c *Cell) { c.Attr.Blink = true },
		"conceal": func(c *Cell) { c.Attr.Conceal = true },
	}
	for name, apply := range fields {
		cell := Cell{Data: string(name[0]), Attr: Attr{Fg: Color{Name: "default", Mode: ColorDefault}, Bg: Color{Name: "default", Mode: ColorDefault}}}
		apply(&cell)
		switch name {
		case "fg":
			if cell.Attr.Fg.Name != "true" {
				t.Fatalf("expected fg set")
			}
		case "bg":
			if cell.Attr.Bg.Name != "true" {
				t.Fatalf("expected bg set")
			}
		case "bold":
			if !cell.Attr.Bold {
				t.Fatalf("expected bold set")
			}
		case "italics":
			if !cell.Attr.Italics {
				t.Fatalf("expected italics set")
			}
		case "underscore":
			if !cell.Attr.Underline {
				t.Fatalf("expected underline set")
			}
		case "strikethrough":
			if !cell.Attr.Strikethrough {
				t.Fatalf("expected strikethrough set")
			}
		case "reverse":
			if !cell.Attr.Reverse {
				t.Fatalf("expected reverse set")
			}
		case "blink":
			if !cell.Attr.Blink {
				t.Fatalf("expected blink set")
			}
		case "conceal":
			if !cell.Attr.Conceal {
				t.Fatalf("expected conceal set")
			}
		}
	}
}

// From pyte/tests/test_screen.py::test_remove_non_existant_attribute
func TestPyteTestScreenRemoveNonExistantAttribute(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{24}, false)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	if screen.Cursor.Attr.Underline {
		t.Fatalf("expected underline off")
	}
}

// From pyte/tests/test_screen.py::test_attributes
func TestPyteTestScreenAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{1}, false)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	if !screen.Cursor.Attr.Bold {
		t.Fatalf("expected bold cursor")
	}
	screen.Draw("f")
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "f", withBold(true)), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
}

// From pyte/tests/test_screen.py::test_blink
func TestPyteTestScreenBlink(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{5}, false)
	screen.Draw("f")
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "f", withBlink(true)), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
}

// From pyte/tests/test_screen.py::test_colors
func TestPyteTestScreenColors(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{30}, false)
	screen.SelectGraphicRendition([]int{40}, false)
	if screen.Cursor.Attr.Fg.Name != "black" || screen.Cursor.Attr.Bg.Name != "black" {
		t.Fatalf("expected black fg/bg")
	}
	screen.SelectGraphicRendition([]int{31}, false)
	if screen.Cursor.Attr.Fg.Name != "red" || screen.Cursor.Attr.Bg.Name != "black" {
		t.Fatalf("expected red fg")
	}
}

// From pyte/tests/test_screen.py::test_colors256
func TestPyteTestScreenColors256(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{SgrFg256, 5, 0}, false)
	screen.SelectGraphicRendition([]int{SgrBg256, 5, 15}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 256 colors")
	}
	screen.SelectGraphicRendition([]int{48, 5, 100500}, false)
}

// From pyte/tests/test_screen.py::test_colors256_missing_attrs
func TestPyteTestScreenColors256MissingAttrs(t *testing.T) {
	screen := NewScreen(2, 2)
	defaultAttr := screen.Cursor.Attr
	screen.SelectGraphicRendition([]int{SgrFg256}, false)
	screen.SelectGraphicRendition([]int{SgrBg256}, false)
	if screen.Cursor.Attr != defaultAttr {
		t.Fatalf("expected default attrs")
	}
}

// From pyte/tests/test_screen.py::test_colors24bit
func TestPyteTestScreenColors24Bit(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{38, 2, 0, 0, 0}, false)
	screen.SelectGraphicRendition([]int{48, 2, 255, 255, 255}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 24-bit colors")
	}
	screen.SelectGraphicRendition([]int{48, 2, 255}, false)
}

// From pyte/tests/test_screen.py::test_colors_aixterm
func TestPyteTestScreenColorsAixterm(t *testing.T) {
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

// From pyte/tests/test_screen.py::test_colors_ignore_invalid
func TestPyteTestScreenColorsIgnoreInvalid(t *testing.T) {
	screen := NewScreen(2, 2)
	defaultAttr := screen.Cursor.Attr
	screen.SelectGraphicRendition([]int{100500}, false)
	if screen.Cursor.Attr != defaultAttr {
		t.Fatalf("expected default attrs")
	}
	screen.SelectGraphicRendition([]int{38, 100500}, false)
	if screen.Cursor.Attr != defaultAttr {
		t.Fatalf("expected default attrs")
	}
	screen.SelectGraphicRendition([]int{48, 100500}, false)
	if screen.Cursor.Attr != defaultAttr {
		t.Fatalf("expected default attrs")
	}
}

// From pyte/tests/test_screen.py::test_reset_resets_colors
func TestPyteTestScreenResetResetsColors(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{30}, false)
	screen.SelectGraphicRendition([]int{40}, false)
	if screen.Cursor.Attr.Fg.Name != "black" || screen.Cursor.Attr.Bg.Name != "black" {
		t.Fatalf("expected black fg/bg")
	}
	screen.SelectGraphicRendition([]int{0}, false)
	if screen.Cursor.Attr != screen.defaultAttr() {
		t.Fatalf("expected default attrs")
	}
}

// From pyte/tests/test_screen.py::test_reset_works_between_attributes
func TestPyteTestScreenResetWorksBetweenAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{31, 0, 41}, false)
	if screen.Cursor.Attr.Fg.Name != "default" || screen.Cursor.Attr.Bg.Name != "red" {
		t.Fatalf("expected reset fg and red bg")
	}
}

// From pyte/tests/test_screen.py::test_multi_attribs
func TestPyteTestScreenMultiAttribs(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{1}, false)
	screen.SelectGraphicRendition([]int{3}, false)
	if !screen.Cursor.Attr.Bold || !screen.Cursor.Attr.Italics {
		t.Fatalf("expected bold/italics")
	}
}

// From pyte/tests/test_screen.py::test_attributes_reset
func TestPyteTestScreenAttributesReset(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SetMode([]int{ModeLNM}, false)
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{1}, false)
	screen.Draw("f")
	screen.Draw("o")
	screen.Draw("o")
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "f", withBold(true)), cellWith(screen, "o", withBold(true))},
		{cellWith(screen, "o", withBold(true)), screen.defaultCell()},
	})
	screen.CursorPosition(0, 0)
	screen.SelectGraphicRendition([]int{0}, false)
	screen.Draw("f")
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "f"), cellWith(screen, "o", withBold(true))},
		{cellWith(screen, "o", withBold(true)), screen.defaultCell()},
	})
}

// From pyte/tests/test_screen.py::test_resize
func TestPyteTestScreenResize(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.SetMargins(0, 1)
	if screen.Columns != 2 || screen.Lines != 2 {
		t.Fatalf("expected 2x2")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.Resize(3, 3)
	if screen.Columns != 3 || screen.Lines != 3 {
		t.Fatalf("expected 3x3")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
	})
	if _, ok := screen.Mode[ModeDECOM]; !ok {
		t.Fatalf("expected DECOM set")
	}
	if screen.Margins != nil {
		t.Fatalf("expected margins reset")
	}
	screen.Resize(2, 2)
	if screen.Columns != 2 || screen.Lines != 2 {
		t.Fatalf("expected 2x2")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, []int{})
	screen.Resize(2, 3)
	assertDisplay(t, screen, []string{"bo ", "sh "})
	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, []int{})
	screen.Resize(2, 1)
	assertDisplay(t, screen, []string{"b", "s"})
	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, []int{})
	screen.Resize(3, 2)
	assertDisplay(t, screen, []string{"bo", "sh", "  "})
	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, []int{})
	screen.Resize(1, 2)
	assertDisplay(t, screen, []string{"sh"})
}

// From pyte/tests/test_screen.py::test_resize_same
func TestPyteTestScreenResizeSame(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.Dirty = map[int]struct{}{}
	screen.Resize(2, 2)
	if len(screen.Dirty) != 0 {
		t.Fatalf("expected no dirty rows")
	}
}

// From pyte/tests/test_screen.py::test_set_mode
func TestPyteTestScreenSetMode(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{})
	screen.CursorPosition(1, 1)
	screen.SetMode([]int{ModeDECCOLM >> 5}, true)
	for line := range tolist(screen) {
		for _, char := range tolist(screen)[line] {
			if char.Data != screen.defaultCell().Data || char.Attr != screen.defaultCell().Attr {
				t.Fatalf("expected default cell")
			}
		}
	}
	if screen.Columns != 132 {
		t.Fatalf("expected 132 columns")
	}
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor origin")
	}
	screen.ResetMode([]int{ModeDECCOLM >> 5}, true)
	if screen.Columns != 3 {
		t.Fatalf("expected columns reset")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{})
	screen.CursorPosition(1, 1)
	screen.SetMode([]int{ModeDECOM}, false)
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor origin")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{})
	screen.SetMode([]int{ModeDECSCNM >> 5}, true)
	for line := range tolist(screen) {
		for _, char := range tolist(screen)[line] {
			if !char.Attr.Reverse {
				t.Fatalf("expected reverse")
			}
		}
	}
	if !screen.defaultAttr().Reverse {
		t.Fatalf("expected reverse default")
	}
	screen.ResetMode([]int{ModeDECSCNM >> 5}, true)
	for line := range tolist(screen) {
		for _, char := range tolist(screen)[line] {
			if char.Attr.Reverse {
				t.Fatalf("expected no reverse")
			}
		}
	}
	if screen.defaultAttr().Reverse {
		t.Fatalf("expected default reverse off")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{})
	screen.Cursor.Hidden = true
	screen.SetMode([]int{ModeDECTCEM}, false)
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
	screen.ResetMode([]int{ModeDECTCEM}, false)
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
}

// From pyte/tests/test_screen.py::test_draw
func TestPyteTestScreenDraw(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.SetMode([]int{ModeLNM}, false)
	if _, ok := screen.Mode[ModeDECAWM]; !ok {
		t.Fatalf("expected DECAWM")
	}
	for _, ch := range "abc" {
		screen.Draw(string(ch))
	}
	assertDisplay(t, screen, []string{"abc", "   ", "   "})
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 3 {
		t.Fatalf("unexpected cursor")
	}
	screen.Draw("a")
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("unexpected cursor after wrap")
	}

	screen = NewScreen(3, 3)
	screen.ResetMode([]int{ModeDECAWM}, false)
	for _, ch := range "abc" {
		screen.Draw(string(ch))
	}
	assertDisplay(t, screen, []string{"abc", "   ", "   "})
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 3 {
		t.Fatalf("unexpected cursor")
	}
	screen.Draw("a")
	assertDisplay(t, screen, []string{"aba", "   ", "   "})
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 3 {
		t.Fatalf("unexpected cursor without wrap")
	}

	screen.SetMode([]int{ModeIRM}, false)
	screen.CursorPosition(0, 0)
	screen.Draw("x")
	assertDisplay(t, screen, []string{"xab", "   ", "   "})
	screen.CursorPosition(0, 0)
	screen.Draw("y")
	assertDisplay(t, screen, []string{"yxa", "   ", "   "})
}

// From pyte/tests/test_screen.py::test_draw_russian
func TestPyteTestScreenDrawRussian(t *testing.T) {
	screen := NewScreen(20, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("Нерусский текст"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"Нерусский текст     "})
}

// From pyte/tests/test_screen.py::test_draw_multiple_chars
func TestPyteTestScreenDrawMultipleChars(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("foobar")
	if screen.Cursor.Col != 6 {
		t.Fatalf("expected cursor at 6")
	}
	assertDisplay(t, screen, []string{"foobar    "})
}

// From pyte/tests/test_screen.py::test_draw_utf8
func TestPyteTestScreenDrawUTF8(t *testing.T) {
	screen := NewScreen(1, 1)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte{0xE2, 0x80, 0x9D}); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"”"})
}

// From pyte/tests/test_screen.py::test_draw_width2
func TestPyteTestScreenDrawWidth2(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("コンニチハ")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
}

// From pyte/tests/test_screen.py::test_draw_width2_line_end
func TestPyteTestScreenDrawWidth2LineEnd(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw(" コンニチハ")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
}

// From pyte/tests/test_screen.py::test_draw_width2_irm
func TestPyteTestScreenDrawWidth2IRM(t *testing.T) {
	t.Skip("pyte marks this test xfail")
}

// From pyte/tests/test_screen.py::test_draw_width0_combining
func TestPyteTestScreenDrawWidth0Combining(t *testing.T) {
	screen := NewScreen(4, 2)
	screen.Draw("\u0308")
	assertDisplay(t, screen, []string{"    ", "    "})
	screen.Draw("bad")
	screen.Draw("\u0308")
	assertDisplay(t, screen, []string{"bad̈ ", "    "})
	screen.Draw("!")
	screen.Draw("\u0308")
	assertDisplay(t, screen, []string{"bad̈!̈", "    "})
}

// From pyte/tests/test_screen.py::test_draw_width0_irm
func TestPyteTestScreenDrawWidth0IRM(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.SetMode([]int{ModeIRM}, false)
	screen.Draw("\u200b")
	screen.Draw("\u0007")
	assertDisplay(t, screen, []string{"          "})
}

// From pyte/tests/test_screen.py::test_draw_width0_decawm_off
func TestPyteTestScreenDrawWidth0DECAWMOff(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.ResetMode([]int{ModeDECAWM}, false)
	screen.Draw(" コンニチハ")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
	screen.Draw("\u200b")
	screen.Draw("\u0007")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor unchanged")
	}
}

// From pyte/tests/test_screen.py::test_draw_cp437
func TestPyteTestScreenDrawCP437(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewByteStream(screen, false)
	if screen.Charset != 0 {
		t.Fatalf("expected charset 0")
	}
	screen.DefineCharset("U", "(")
	stream.SelectOtherCharset("@")
	data := []byte{0xE0, 0x20, 0xF1, 0x20, 0xEE}
	if err := stream.Feed(data); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"α ± ε"})
}

// From pyte/tests/test_screen.py::test_draw_with_carriage_return
func TestPyteTestScreenDrawWithCarriageReturn(t *testing.T) {
	line := "ipcs -s | grep nobody |awk '{print$2}'|xargs -n1 i" +
		"pcrm sem ;ps aux|grep -P 'httpd|fcgi'|grep -v grep" +
		"|awk '{print$2 \r}'|xargs kill -9;/etc/init.d/ht" +
		"tpd startssl"
	screen := NewScreen(50, 3)
	stream := NewStream(screen, false)
	if err := stream.Feed(line); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{
		"ipcs -s | grep nobody |awk '{print$2}'|xargs -n1 i",
		"pcrm sem ;ps aux|grep -P 'httpd|fcgi'|grep -v grep",
		"}'|xargs kill -9;/etc/init.d/httpd startssl       ",
	})
}

// From pyte/tests/test_screen.py::test_display_wcwidth
func TestPyteTestScreenDisplayWcwidth(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("コンニチハ")
	assertDisplay(t, screen, []string{"コンニチハ"})
}

// From pyte/tests/test_screen.py::test_display_multi_char_emoji
func TestPyteTestScreenDisplayMultiCharEmoji(t *testing.T) {
	screen := NewScreen(4, 1)
	screen.Draw("👨\u200d💻a")
	assertDisplay(t, screen, []string{"👨\u200d💻a "})
}

// From pyte/tests/test_screen.py::test_display_complex_emoji
func TestPyteTestScreenDisplayComplexEmoji(t *testing.T) {
	emoji := "\U0001f926\U0001f3fd\u200d\u2642\ufe0f"
	screen := NewScreen(4, 1)
	screen.Draw(emoji + "a")
	assertDisplay(t, screen, []string{emoji + "a "})
}

// From pyte/tests/test_screen.py::test_carriage_return
func TestPyteTestScreenCarriageReturn(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.Cursor.Col = 2
	screen.CarriageReturn()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected col 0")
	}
}

// From pyte/tests/test_screen.py::test_insert_lines
func TestPyteTestScreenInsertLines(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{1})
	screen.InsertLines(1)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"   ", "sam", "is "})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{cellWith(screen, "i", withFg("red")), cellWith(screen, "s", withFg("red")), cellWith(screen, " ", withFg("red"))},
	})

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{1})
	screen.InsertLines(2)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"   ", "   ", "sam"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.InsertLines(1)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "   ", "is ", "foo", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "i"), cellWith(screen, "s"), cellWith(screen, " ")},
		{cellWith(screen, "f", withFg("red")), cellWith(screen, "o", withFg("red")), cellWith(screen, "o", withFg("red"))},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(1, 3)
	screen.Cursor.Row = 1
	screen.InsertLines(1)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "   ", "is ", "bar", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "i"), cellWith(screen, "s"), cellWith(screen, " ")},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen.InsertLines(2)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "   ", "   ", "bar", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 1
	screen.InsertLines(20)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "   ", "   ", "   ", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(2, 4)
	screen.InsertLines(5)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (0,0)")
	}
	assertDisplay(t, screen, []string{"sam", "is ", "foo", "bar", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{cellWith(screen, "i"), cellWith(screen, "s"), cellWith(screen, " ")},
		{cellWith(screen, "f", withFg("red")), cellWith(screen, "o", withFg("red")), cellWith(screen, "o", withFg("red"))},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})
}

// From pyte/tests/test_screen.py::test_delete_lines
func TestPyteTestScreenDeleteLines(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{1})
	screen.DeleteLines(1)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"is ", "foo", "   "})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "i", withFg("red")), cellWith(screen, "s", withFg("red")), cellWith(screen, " ", withFg("red"))},
		{cellWith(screen, "f"), cellWith(screen, "o"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
	})

	screen.DeleteLines(0)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"foo", "   ", "   "})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "f"), cellWith(screen, "o"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(1)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "foo", "bar", "   ", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{cellWith(screen, "f", withFg("red")), cellWith(screen, "o", withFg("red")), cellWith(screen, "o", withFg("red"))},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(2)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "bar", "   ", "   ", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{})
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(5)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"sam", "   ", "   ", "   ", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, []int{2, 3})
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 0
	screen.DeleteLines(5)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (0,0)")
	}
	assertDisplay(t, screen, []string{"sam", "is ", "foo", "bar", "baz"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "s"), cellWith(screen, "a"), cellWith(screen, "m")},
		{cellWith(screen, "i"), cellWith(screen, "s"), cellWith(screen, " ")},
		{cellWith(screen, "f", withFg("red")), cellWith(screen, "o", withFg("red")), cellWith(screen, "o", withFg("red"))},
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "a", withFg("red")), cellWith(screen, "r", withFg("red"))},
		{cellWith(screen, "b"), cellWith(screen, "a"), cellWith(screen, "z")},
	})
}

// From pyte/tests/test_screen.py::test_insert_characters
func TestPyteTestScreenInsertCharacters(t *testing.T) {
	screen := updateScreen(NewScreen(3, 4), []string{"sam", "is ", "foo", "bar"}, []int{0})
	cursor := screen.Cursor
	screen.InsertCharacters(2)
	if screen.Cursor != cursor {
		t.Fatalf("expected cursor unchanged")
	}
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, "s", withFg("red")),
	}})

	screen.Cursor.Row = 2
	screen.Cursor.Col = 1
	screen.InsertCharacters(1)
	assertCellsEqual(t, [][]Cell{screen.Buffer[2]}, [][]Cell{{
		cellWith(screen, "f"),
		screen.defaultCell(),
		cellWith(screen, "o"),
	}})

	screen.Cursor.Row = 3
	screen.Cursor.Col = 1
	screen.InsertCharacters(10)
	assertCellsEqual(t, [][]Cell{screen.Buffer[3]}, [][]Cell{{
		cellWith(screen, "b"),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{0})
	screen.CursorPosition(0, 0)
	screen.InsertCharacters(1)
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		cellWith(screen, "s", withFg("red")),
		cellWith(screen, "a", withFg("red")),
	}})

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{0})
	screen.CursorPosition(0, 0)
	screen.InsertCharacters(1)
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		cellWith(screen, "s", withFg("red")),
		cellWith(screen, "a", withFg("red")),
	}})
}

// From pyte/tests/test_screen.py::test_delete_characters
func TestPyteTestScreenDeleteCharacters(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{0})
	screen.DeleteCharacters(2)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"m  ", "is ", "foo"})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "m", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen.Cursor.Row = 2
	screen.Cursor.Col = 2
	screen.DeleteCharacters(1)
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"m  ", "is ", "fo "})

	screen.Cursor.Row = 1
	screen.Cursor.Col = 1
	screen.DeleteCharacters(1)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"m  ", "i  ", "fo "})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.Cursor.Col = 1
	screen.DeleteCharacters(3)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 1 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"15   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "1", withFg("red")),
		cellWith(screen, "5", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.Cursor.Col = 2
	screen.DeleteCharacters(10)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"12   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "1", withFg("red")),
		cellWith(screen, "2", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.DeleteCharacters(4)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"5    "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "5", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
	}})
}

// From pyte/tests/test_screen.py::test_erase_character
func TestPyteTestScreenEraseCharacter(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, []int{0})
	screen.EraseCharacters(2)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	assertDisplay(t, screen, []string{"  m", "is ", "foo"})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, "m", withFg("red")),
	}})

	screen.Cursor.Row = 2
	screen.Cursor.Col = 2
	screen.EraseCharacters(1)
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"  m", "is ", "fo "})

	screen.Cursor.Row = 1
	screen.Cursor.Col = 1
	screen.EraseCharacters(1)
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"  m", "i  ", "fo "})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.Cursor.Col = 1
	screen.EraseCharacters(3)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 1 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"1   5"})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "1", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, "5", withFg("red")),
	}})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.Cursor.Col = 2
	screen.EraseCharacters(10)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"12   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "1", withFg("red")),
		cellWith(screen, "2", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, []int{0})
	screen.EraseCharacters(4)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor unchanged")
	}
	assertDisplay(t, screen, []string{"    5"})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, "5", withFg("red")),
	}})
}

// From pyte/tests/test_screen.py::test_erase_in_line
func TestPyteTestScreenEraseInLine(t *testing.T) {
	screen := updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{0})
	screen.CursorPosition(1, 3)
	screen.EraseInLine(0, false)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (0,2)")
	}
	assertDisplay(t, screen, []string{"sa   ", "s foo", "but a", "re yo", "u?   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		cellWith(screen, "s", withFg("red")),
		cellWith(screen, "a", withFg("red")),
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{0})
	screen.EraseInLine(1, false)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (0,2)")
	}
	assertDisplay(t, screen, []string{"    i", "s foo", "but a", "re yo", "u?   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(),
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, " ", withFg("red")),
		cellWith(screen, "i", withFg("red")),
	}})

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{0})
	screen.EraseInLine(2, false)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (0,2)")
	}
	assertDisplay(t, screen, []string{"     ", "s foo", "but a", "re yo", "u?   "})
	assertCellsEqual(t, [][]Cell{screen.Buffer[0]}, [][]Cell{{
		screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(),
	}})
}

// From pyte/tests/test_screen.py::test_erase_in_display
func TestPyteTestScreenEraseInDisplay(t *testing.T) {
	screen := updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{2, 3})
	screen.CursorPosition(3, 3)
	screen.EraseInDisplay(0, false)
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (2,2)")
	}
	assertDisplay(t, screen, []string{"sam i", "s foo", "bu   ", "     ", "     "})
	assertCellsEqual(t, screen.Buffer[2:], [][]Cell{
		{cellWith(screen, "b", withFg("red")), cellWith(screen, "u", withFg("red")), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{2, 3})
	screen.EraseInDisplay(1, false)
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (2,2)")
	}
	assertDisplay(t, screen, []string{"     ", "     ", "    a", "re yo", "u?   "})
	assertCellsEqual(t, screen.Buffer[:3], [][]Cell{
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), cellWith(screen, " ", withFg("red")), cellWith(screen, "a", withFg("red"))},
	})

	screen.EraseInDisplay(2, false)
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor at (2,2)")
	}
	assertDisplay(t, screen, []string{"     ", "     ", "     ", "     ", "     "})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{2, 3})
	screen.EraseInDisplay(3, true)
	assertDisplay(t, screen, []string{"     ", "     ", "     ", "     ", "     "})

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{2, 3})
	screen.EraseInDisplay(3, false, 0)
	assertDisplay(t, screen, []string{"     ", "     ", "     ", "     ", "     "})

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, []int{2, 3})
	screen.EraseInDisplay(3, true, 0)
	assertDisplay(t, screen, []string{"     ", "     ", "     ", "     ", "     "})
}

// From pyte/tests/test_screen.py::test_cursor_up
func TestPyteTestScreenCursorUp(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorUp(1)
	if screen.Cursor.Row != 0 {
		t.Fatalf("expected row 0")
	}
	screen.Cursor.Row = 1
	screen.CursorUp(10)
	if screen.Cursor.Row != 0 {
		t.Fatalf("expected row 0")
	}
	screen.Cursor.Row = 5
	screen.CursorUp(3)
	if screen.Cursor.Row != 2 {
		t.Fatalf("expected row 2")
	}
}

// From pyte/tests/test_screen.py::test_cursor_down
func TestPyteTestScreenCursorDown(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.Cursor.Row = 9
	screen.CursorDown(1)
	if screen.Cursor.Row != 9 {
		t.Fatalf("expected row 9")
	}
	screen.Cursor.Row = 8
	screen.CursorDown(10)
	if screen.Cursor.Row != 9 {
		t.Fatalf("expected row 9")
	}
	screen.Cursor.Row = 5
	screen.CursorDown(3)
	if screen.Cursor.Row != 8 {
		t.Fatalf("expected row 8")
	}
}

// From pyte/tests/test_screen.py::test_cursor_back
func TestPyteTestScreenCursorBack(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.Cursor.Col = 0
	screen.CursorBack(1)
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected col 0")
	}
	screen.Cursor.Col = 3
	screen.CursorBack(10)
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected col 0")
	}
	screen.Cursor.Col = 5
	screen.CursorBack(3)
	if screen.Cursor.Col != 2 {
		t.Fatalf("expected col 2")
	}
}

// From pyte/tests/test_screen.py::test_cursor_back_last_column
func TestPyteTestScreenCursorBackLastColumn(t *testing.T) {
	screen := NewScreen(13, 1)
	screen.Draw("Hello, world!")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
	screen.CursorBack(5)
	if screen.Cursor.Col != (screen.Columns-1)-5 {
		t.Fatalf("expected cursor position")
	}
}

// From pyte/tests/test_screen.py::test_cursor_forward
func TestPyteTestScreenCursorForward(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.Cursor.Col = 9
	screen.CursorForward(1)
	if screen.Cursor.Col != 9 {
		t.Fatalf("expected col 9")
	}
	screen.Cursor.Col = 8
	screen.CursorForward(10)
	if screen.Cursor.Col != 9 {
		t.Fatalf("expected col 9")
	}
	screen.Cursor.Col = 5
	screen.CursorForward(3)
	if screen.Cursor.Col != 8 {
		t.Fatalf("expected col 8")
	}
}

// From pyte/tests/test_screen.py::test_cursor_position
func TestPyteTestScreenCursorPosition(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorPosition(5, 10)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 9 {
		t.Fatalf("expected cursor (4,9)")
	}
	screen.CursorPosition(0, 10)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 9 {
		t.Fatalf("expected cursor (0,9)")
	}
	screen.CursorPosition(100, 5)
	if screen.Cursor.Row != 9 || screen.Cursor.Col != 4 {
		t.Fatalf("expected cursor (9,4)")
	}
	screen.CursorPosition(5, 100)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 9 {
		t.Fatalf("expected cursor (4,9)")
	}
	screen.SetMargins(5, 9)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.CursorPosition(0, 0)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor (4,0)")
	}
	screen.CursorPosition(2, 0)
	if screen.Cursor.Row != 5 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor (5,0)")
	}
	screen.CursorPosition(10, 0)
	if screen.Cursor.Row != 5 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor unchanged")
	}
}

// From pyte/tests/test_screen.py::test_unicode
func TestPyteTestScreenUnicode(t *testing.T) {
	screen := NewScreen(4, 2)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("тест")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"тест", "    "})
}

// From pyte/tests/test_screen.py::test_alignment_display
func TestPyteTestScreenAlignmentDisplay(t *testing.T) {
	screen := NewScreen(5, 5)
	screen.SetMode([]int{ModeLNM}, false)
	screen.Draw("a")
	screen.LineFeed()
	screen.LineFeed()
	screen.Draw("b")
	assertDisplay(t, screen, []string{"a    ", "     ", "b    ", "     ", "     "})
	screen.AlignmentDisplay()
	assertDisplay(t, screen, []string{"EEEEE", "EEEEE", "EEEEE", "EEEEE", "EEEEE"})
}

// From pyte/tests/test_screen.py::test_set_margins
func TestPyteTestScreenSetMargins(t *testing.T) {
	screen := NewScreen(10, 10)
	if screen.Margins != nil {
		t.Fatalf("expected nil margins")
	}
	screen.SetMargins(1, 5)
	if screen.Margins == nil || screen.Margins.Top != 0 || screen.Margins.Bottom != 4 {
		t.Fatalf("expected margins (0,4)")
	}
	screen.SetMargins(100, 10)
	if screen.Margins == nil || screen.Margins.Top != 0 || screen.Margins.Bottom != 4 {
		t.Fatalf("expected margins unchanged")
	}
	screen.SetMargins(0, 0)
	if screen.Margins != nil {
		t.Fatalf("expected margins reset")
	}
}

// From pyte/tests/test_screen.py::test_set_margins_zero
func TestPyteTestScreenSetMarginsZero(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(1, 5)
	if screen.Margins == nil || screen.Margins.Top != 0 || screen.Margins.Bottom != 4 {
		t.Fatalf("expected margins (0,4)")
	}
	screen.SetMargins(0, 0)
	if screen.Margins != nil {
		t.Fatalf("expected margins reset")
	}
}

// From pyte/tests/test_screen.py::test_hide_cursor
func TestPyteTestScreenHideCursor(t *testing.T) {
	screen := NewScreen(10, 10)
	if _, ok := screen.Mode[ModeDECTCEM]; !ok {
		t.Fatalf("expected DECTCEM")
	}
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

// From pyte/tests/test_screen.py::test_report_device_attributes
func TestPyteTestScreenReportDeviceAttributes(t *testing.T) {
	screen := NewScreen(10, 10)
	var acc []string
	screen.WriteProcessInput = func(data string) { acc = append(acc, data) }
	screen.ReportDeviceAttributes(42, false, 0)
	if len(acc) != 0 {
		t.Fatalf("expected no response")
	}
	screen.ReportDeviceAttributes(0, false, 0)
	if len(acc) != 1 || acc[0] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response")
	}
}

// From pyte/tests/test_screen.py::test_report_device_status
func TestPyteTestScreenReportDeviceStatus(t *testing.T) {
	screen := NewScreen(10, 10)
	var acc []string
	screen.WriteProcessInput = func(data string) { acc = append(acc, data) }
	screen.ReportDeviceStatus(42, false, 0)
	if len(acc) != 0 {
		t.Fatalf("expected no response")
	}
	screen.ReportDeviceStatus(5, false, 0)
	if acc[len(acc)-1] != ControlCSI+"0n" {
		t.Fatalf("unexpected terminal status")
	}
	screen.CursorToColumn(5)
	screen.ReportDeviceStatus(6, false, 0)
	if acc[len(acc)-1] != ControlCSI+"1;5R" {
		t.Fatalf("unexpected cursor report")
	}
	screen.CursorPosition(0, 0)
	screen.SetMargins(5, 9)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.CursorToLine(5)
	screen.ReportDeviceStatus(6, false, 0)
	if acc[len(acc)-1] != ControlCSI+"5;1R" {
		t.Fatalf("unexpected origin cursor report")
	}
}

// From pyte/tests/test_screen.py::test_screen_set_icon_name_title
func TestPyteTestScreenSetIconNameTitle(t *testing.T) {
	screen := NewScreen(10, 1)
	text := "±"
	screen.SetIconName(text)
	if screen.IconName != text {
		t.Fatalf("expected icon name")
	}
	screen.SetTitle(text)
	if screen.Title != text {
		t.Fatalf("expected title")
	}
}

// From pyte/tests/test_screen.py::test_private_sgr_sequence_ignored
func TestPyteTestScreenPrivateSGRSequenceIgnored(t *testing.T) {
	cursorState := func(cursor Cursor) Cursor { return cursor }
	screen := NewHistoryScreen(2, 2, 0)
	stream := NewByteStream(screen, false)
	display := append([]string(nil), screen.Display()...)
	mode := make(map[int]struct{})
	for k := range screen.Mode {
		mode[k] = struct{}{}
	}
	cursor := cursorState(screen.Cursor)
	if err := stream.Feed([]byte{0x1b, '[', '?', '4', 'm'}); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if got := screen.Display(); len(got) != len(display) || got[0] != display[0] || got[1] != display[1] {
		t.Fatalf("display changed")
	}
	if len(screen.Mode) != len(mode) {
		t.Fatalf("mode changed")
	}
	for k := range mode {
		if _, ok := screen.Mode[k]; !ok {
			t.Fatalf("mode changed")
		}
	}
	if screen.Cursor != cursor {
		t.Fatalf("cursor changed")
	}
}

// From pyte/tests/test_screen.py::test_index
func TestPyteTestScreenIndex(t *testing.T) {
	screen := updateScreen(NewScreen(2, 2), []string{"wo", "ot"}, []int{1})
	screen.Index()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "w"), cellWith(screen, "o")},
		{cellWith(screen, "o", withFg("red")), cellWith(screen, "t", withFg("red"))},
	})

	screen.Index()
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected cursor row 1")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "o", withFg("red")), cellWith(screen, "t", withFg("red"))},
		{screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(NewScreen(2, 5), []string{"bo", "sh", "th", "er", "oh"}, []int{1, 2})
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 3
	screen.Index()
	if screen.Cursor.Row != 3 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (3,0)")
	}
	assertDisplay(t, screen, []string{"bo", "th", "er", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{cellWith(screen, "t", withFg("red")), cellWith(screen, "h", withFg("red"))},
		{cellWith(screen, "e"), cellWith(screen, "r")},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.Index()
	if screen.Cursor.Row != 3 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (3,0)")
	}
	assertDisplay(t, screen, []string{"bo", "er", "  ", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{cellWith(screen, "e"), cellWith(screen, "r")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.Index()
	if screen.Cursor.Row != 3 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (3,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "  ", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.Index()
	if screen.Cursor.Row != 3 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (3,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "  ", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})
}

// From pyte/tests/test_screen.py::test_reverse_index
func TestPyteTestScreenReverseIndex(t *testing.T) {
	screen := updateScreen(NewScreen(2, 2), []string{"wo", "ot"}, []int{0})
	screen.ReverseIndex()
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (0,0)")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "w", withFg("red")), cellWith(screen, "o", withFg("red"))},
	})

	screen.ReverseIndex()
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (0,0)")
	}
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(NewScreen(2, 5), []string{"bo", "sh", "th", "er", "oh"}, []int{2, 3})
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 1
	screen.ReverseIndex()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "sh", "th", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "s"), cellWith(screen, "h")},
		{cellWith(screen, "t", withFg("red")), cellWith(screen, "h", withFg("red"))},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.ReverseIndex()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "  ", "sh", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "s"), cellWith(screen, "h")},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.ReverseIndex()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "  ", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})

	screen.ReverseIndex()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor at (1,0)")
	}
	assertDisplay(t, screen, []string{"bo", "  ", "  ", "  ", "oh"})
	assertCellsEqual(t, tolist(screen), [][]Cell{
		{cellWith(screen, "b"), cellWith(screen, "o")},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "o"), cellWith(screen, "h")},
	})
}

// From pyte/tests/test_screen.py::test_linefeed
func TestPyteTestScreenLinefeed(t *testing.T) {
	screen := updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, []int{})
	screen.SetMode([]int{ModeLNM}, false)
	screen.Cursor.Col = 1
	screen.Cursor.Row = 0
	screen.LineFeed()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor moved with CR")
	}

	screen.ResetMode([]int{ModeLNM}, false)
	screen.Cursor.Col = 1
	screen.Cursor.Row = 0
	screen.LineFeed()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("expected cursor moved without CR")
	}
}

// From pyte/tests/test_screen.py::test_linefeed_margins
func TestPyteTestScreenLinefeedMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(3, 27)
	screen.CursorPosition(0, 0)
	screen.LineFeed()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor moved")
	}
}

// From pyte/tests/test_screen.py::test_tabstops
func TestPyteTestScreenTabstops(t *testing.T) {
	screen := NewScreen(10, 10)
	if len(screen.TabStops) != 1 {
		t.Fatalf("expected one tabstop")
	}
	if _, ok := screen.TabStops[8]; !ok {
		t.Fatalf("expected tabstop at 8")
	}
	screen.ClearTabStop(3)
	if len(screen.TabStops) != 0 {
		t.Fatalf("expected cleared tabstops")
	}
	screen.Cursor.Col = 1
	screen.SetTabStop()
	screen.Cursor.Col = 8
	screen.SetTabStop()
	screen.Cursor.Col = 0
	screen.Tab()
	if screen.Cursor.Col != 1 {
		t.Fatalf("expected tab to col 1")
	}
	screen.Tab()
	if screen.Cursor.Col != 8 {
		t.Fatalf("expected tab to col 8")
	}
	screen.Tab()
	if screen.Cursor.Col != 9 {
		t.Fatalf("expected tab to col 9")
	}
	screen.Tab()
	if screen.Cursor.Col != 9 {
		t.Fatalf("expected tab to stay at col 9")
	}
}

// From pyte/tests/test_screen.py::test_clear_tabstops
func TestPyteTestScreenClearTabstops(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.ClearTabStop(3)
	screen.Cursor.Col = 1
	screen.SetTabStop()
	screen.Cursor.Col = 5
	screen.SetTabStop()
	screen.ClearTabStop(0)
	if len(screen.TabStops) != 1 {
		t.Fatalf("expected one tabstop")
	}
	screen.SetTabStop()
	screen.ClearTabStop(0)
	if len(screen.TabStops) != 1 {
		t.Fatalf("expected one tabstop")
	}
	screen.SetTabStop()
	screen.Cursor.Col = 9
	screen.SetTabStop()
	screen.ClearTabStop(3)
	if len(screen.TabStops) != 0 {
		t.Fatalf("expected cleared tabstops")
	}
}

// From pyte/tests/test_screen.py::test_backspace
func TestPyteTestScreenBackspace(t *testing.T) {
	screen := NewScreen(2, 2)
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor col 0")
	}
	screen.Backspace()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor col 0")
	}
	screen.Cursor.Col = 1
	screen.Backspace()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor col 0")
	}
}

// From pyte/tests/test_screen.py::test_save_cursor
func TestPyteTestScreenSaveCursor(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.SaveCursor()
	screen.Cursor.Col = 3
	screen.Cursor.Row = 5
	screen.SaveCursor()
	screen.Cursor.Col = 4
	screen.Cursor.Row = 4
	screen.RestoreCursor()
	if screen.Cursor.Col != 3 || screen.Cursor.Row != 5 {
		t.Fatalf("expected cursor (3,5)")
	}
	screen.RestoreCursor()
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor (0,0)")
	}

	screen = NewScreen(10, 10)
	screen.SetMode([]int{ModeDECAWM, ModeDECOM}, false)
	screen.SaveCursor()
	screen.ResetMode([]int{ModeDECAWM}, false)
	screen.RestoreCursor()
	if _, ok := screen.Mode[ModeDECAWM]; !ok {
		t.Fatalf("expected DECAWM")
	}
	if _, ok := screen.Mode[ModeDECOM]; !ok {
		t.Fatalf("expected DECOM")
	}

	screen = NewScreen(10, 10)
	screen.SelectGraphicRendition([]int{4}, false)
	screen.SaveCursor()
	screen.SelectGraphicRendition([]int{24}, false)
	if screen.Cursor.Attr != screen.defaultAttr() {
		t.Fatalf("expected default attrs")
	}
	screen.RestoreCursor()
	if screen.Cursor.Attr == screen.defaultAttr() {
		t.Fatalf("expected underline attr")
	}
	if !screen.Cursor.Attr.Underline {
		t.Fatalf("expected underline")
	}
}

// From pyte/tests/test_screen.py::test_restore_cursor_with_none_saved
func TestPyteTestScreenRestoreCursorWithNoneSaved(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.Cursor.Col = 5
	screen.Cursor.Row = 5
	screen.RestoreCursor()
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor origin")
	}
	if _, ok := screen.Mode[ModeDECOM]; ok {
		t.Fatalf("expected DECOM reset")
	}
}

// From pyte/tests/test_screen.py::test_restore_cursor_out_of_bounds
func TestPyteTestScreenRestoreCursorOutOfBounds(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorPosition(5, 5)
	screen.SaveCursor()
	screen.Resize(3, 3)
	screen.Reset()
	screen.RestoreCursor()
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor (2,2)")
	}

	screen.Resize(10, 10)
	screen.CursorPosition(8, 8)
	screen.SaveCursor()
	screen.Resize(5, 5)
	screen.Reset()
	screen.SetMode([]int{ModeDECOM}, false)
	screen.SetMargins(2, 3)
	screen.RestoreCursor()
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 4 {
		t.Fatalf("expected cursor (2,4)")
	}
}
