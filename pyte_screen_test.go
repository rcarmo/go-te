package te

import "testing"

func TestPyteRemoveNonExistentAttribute(t *testing.T) {
	screen := NewScreen(2, 2)
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
	screen.SelectGraphicRendition([]int{24}, false)
	if screen.Cursor.Attr.Underline {
		t.Fatalf("expected underline off")
	}
}

func TestPyteAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{1}, false)
	if !screen.Cursor.Attr.Bold {
		t.Fatalf("expected bold cursor")
	}
	screen.Draw("f")
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		{cellWith(screen, "f", withBold(true)), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
}

func TestPyteBlink(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{5}, false)
	screen.Draw("f")
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		{cellWith(screen, "f", withBlink(true)), screen.defaultCell()},
		{screen.defaultCell(), screen.defaultCell()},
	})
}

func TestPyteColors(t *testing.T) {
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

func TestPyteColors256(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{SgrFg256, 5, 0}, false)
	screen.SelectGraphicRendition([]int{SgrBg256, 5, 15}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 256 colors")
	}
	screen.SelectGraphicRendition([]int{48, 5, 100500}, false)
}

func TestPyteColors256MissingAttrs(t *testing.T) {
	screen := NewScreen(2, 2)
	defaultAttr := screen.Cursor.Attr
	screen.SelectGraphicRendition([]int{SgrFg256}, false)
	screen.SelectGraphicRendition([]int{SgrBg256}, false)
	if screen.Cursor.Attr != defaultAttr {
		t.Fatalf("expected default attrs")
	}
}

func TestPyteColors24Bit(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{38, 2, 0, 0, 0}, false)
	screen.SelectGraphicRendition([]int{48, 2, 255, 255, 255}, false)
	if screen.Cursor.Attr.Fg.Name != "000000" || screen.Cursor.Attr.Bg.Name != "ffffff" {
		t.Fatalf("expected 24-bit colors")
	}
	screen.SelectGraphicRendition([]int{48, 2, 255}, false)
}

func TestPyteColorsAixterm(t *testing.T) {
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

func TestPyteColorsIgnoreInvalid(t *testing.T) {
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

func TestPyteResetResetsColors(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{30}, false)
	screen.SelectGraphicRendition([]int{40}, false)
	screen.SelectGraphicRendition([]int{0}, false)
	if screen.Cursor.Attr.Fg.Name != "default" || screen.Cursor.Attr.Bg.Name != "default" {
		t.Fatalf("expected default colors")
	}
}

func TestPyteResetBetweenAttributes(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{31, 0, 41}, false)
	if screen.Cursor.Attr.Fg.Name != "default" || screen.Cursor.Attr.Bg.Name != "red" {
		t.Fatalf("expected reset fg and red bg")
	}
}

func TestPyteMultiAttribs(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SelectGraphicRendition([]int{1, 3}, false)
	if !screen.Cursor.Attr.Bold || !screen.Cursor.Attr.Italics {
		t.Fatalf("expected bold/italics")
	}
}

func TestPyteAttributesReset(t *testing.T) {
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

func TestPyteResize(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.SetMargins(0, 1)
	screen.Resize(3, 3)
	if screen.Columns != 3 || screen.Lines != 3 {
		t.Fatalf("expected resize to 3x3")
	}
	if screen.Margins != nil {
		t.Fatalf("expected margins reset")
	}
	screen.Resize(2, 2)
	if screen.Columns != 2 || screen.Lines != 2 {
		t.Fatalf("expected resize back")
	}

	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, nil)
	screen.Resize(2, 3)
	if screen.Display()[0] != "bo " || screen.Display()[1] != "sh " {
		t.Fatalf("expected padded columns")
	}

	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, nil)
	screen.Resize(2, 1)
	if screen.Display()[0] != "b" || screen.Display()[1] != "s" {
		t.Fatalf("expected truncated columns")
	}

	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, nil)
	screen.Resize(3, 2)
	if screen.Display()[2] != "  " {
		t.Fatalf("expected new row")
	}

	screen = updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, nil)
	screen.Resize(1, 2)
	if screen.Display()[0] != "sh" {
		t.Fatalf("expected drop top rows")
	}
}

func TestPyteResizeSame(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.Dirty = map[int]struct{}{}
	screen.Resize(2, 2)
	if len(screen.Dirty) != 0 {
		t.Fatalf("expected dirty unchanged")
	}
}

func TestPyteSetMode(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, nil)
	screen.CursorPosition(1, 1)
	screen.SetMode([]int{ModeDECCOLM}, false)
	for _, line := range screen.Buffer {
		for _, cell := range line {
			if cell.Data != screen.defaultCell().Data {
				t.Fatalf("expected cleared screen")
			}
		}
	}
	if screen.Columns != 132 || screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected 132 columns and cursor reset")
	}
	screen.ResetMode([]int{ModeDECCOLM}, false)
	if screen.Columns != 3 {
		t.Fatalf("expected restore columns")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, nil)
	screen.CursorPosition(1, 1)
	screen.SetMode([]int{ModeDECOM}, false)
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor home")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, nil)
	screen.SetMode([]int{ModeDECSCNM}, false)
	for _, line := range screen.Buffer {
		for _, cell := range line {
			if !cell.Attr.Reverse {
				t.Fatalf("expected reverse")
			}
		}
	}
	screen.ResetMode([]int{ModeDECSCNM}, false)
	for _, line := range screen.Buffer {
		for _, cell := range line {
			if cell.Attr.Reverse {
				t.Fatalf("expected reverse off")
			}
		}
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, nil)
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

func TestPyteDraw(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.SetMode([]int{ModeLNM}, false)
	for _, ch := range "abc" {
		screen.Draw(string(ch))
	}
	if screen.Display()[0] != "abc" {
		t.Fatalf("expected abc")
	}
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 3 {
		t.Fatalf("expected cursor at end")
	}
	screen.Draw("a")
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 1 {
		t.Fatalf("expected linefeed wrap")
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
		t.Fatalf("expected insert")
	}
	screen.CursorPosition(0, 0)
	screen.Draw("y")
	if screen.Display()[0] != "yxa" {
		t.Fatalf("expected insert")
	}
}

func TestPyteDrawRussian(t *testing.T) {
	screen := NewScreen(20, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("Нерусский текст"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "Нерусский текст     " {
		t.Fatalf("unexpected display")
	}
}

func TestPyteDrawMultipleChars(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("foobar")
	if screen.Cursor.Col != 6 {
		t.Fatalf("expected cursor at 6")
	}
	if screen.Display()[0] != "foobar    " {
		t.Fatalf("unexpected display")
	}
}

func TestPyteDrawUTF8(t *testing.T) {
	screen := NewScreen(1, 1)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("\xE2\x80\x9D")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "”" {
		t.Fatalf("expected utf8")
	}
}

func TestPyteDrawWidth2(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("コンニチハ")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
}

func TestPyteDrawWidth2LineEnd(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw(" コンニチハ")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
}

func TestPyteDrawWidth2IRM(t *testing.T) {
	t.Skip("pyte marks as xfail")
}

func TestPyteDrawWidth0Combining(t *testing.T) {
	screen := NewScreen(4, 2)
	screen.Draw("\u0308")
	if screen.Display()[0] != "    " {
		t.Fatalf("expected empty display")
	}
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

func TestPyteDrawWidth0IRM(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.SetMode([]int{ModeIRM}, false)
	screen.Draw("\u200b")
	screen.Draw("\u0007")
	if screen.Display()[0] != "          " {
		t.Fatalf("expected blanks")
	}
}

func TestPyteDrawWidth0DecawmOff(t *testing.T) {
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

func TestPyteDrawCP437(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewByteStream(screen, false)
	screen.DefineCharset("U", "(")
	stream.SelectOtherCharset("@")
	if err := stream.Feed([]byte{0xe0, 0x20, 0xf1, 0x20, 0xee}); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "α ± ε" {
		t.Fatalf("unexpected display: %q", screen.Display()[0])
	}
}

func TestPyteDrawWithCarriageReturn(t *testing.T) {
	line := "ipcs -s | grep nobody |awk '{print$2}'|xargs -n1 i" +
		"pcrm sem ;ps aux|grep -P 'httpd|fcgi'|grep -v grep" +
		"|awk '{print$2 \r}'|xargs kill -9;/etc/init.d/ht" +
		"tpd startssl"

	screen := NewScreen(50, 3)
	stream := NewStream(screen, false)
	if err := stream.Feed(line); err != nil {
		t.Fatalf("feed: %v", err)
	}
	lines := screen.Display()
	if lines[0] != "ipcs -s | grep nobody |awk '{print$2}'|xargs -n1 i" ||
		lines[1] != "pcrm sem ;ps aux|grep -P 'httpd|fcgi'|grep -v grep" ||
		lines[2] != "}'|xargs kill -9;/etc/init.d/httpd startssl       " {
		t.Fatalf("unexpected display")
	}
}

func TestPyteDisplayWidth(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.Draw("コンニチハ")
	if screen.Display()[0] != "コンニチハ" {
		t.Fatalf("unexpected display")
	}
}

func TestPyteDisplayEmoji(t *testing.T) {
	screen := NewScreen(4, 1)
	screen.Draw("👨\u200d💻a")
	if screen.Display()[0] != "👨\u200d💻a " {
		t.Fatalf("unexpected display")
	}
}

func TestPyteDisplayComplexEmoji(t *testing.T) {
	emoji := "\U0001f926\U0001f3fd\u200d\u2642\ufe0f"
	screen := NewScreen(4, 1)
	screen.Draw(emoji + "a")
	if screen.Display()[0] != emoji+"a " {
		t.Fatalf("unexpected display")
	}
}

func TestPyteCarriageReturn(t *testing.T) {
	screen := NewScreen(3, 3)
	screen.Cursor.Col = 2
	screen.CarriageReturn()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected col 0")
	}
}

func TestPyteIndexAndReverseIndex(t *testing.T) {
	screen := updateScreen(NewScreen(2, 2), []string{"wo", "ot"}, coloredLines(1))
	screen.Index()
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected row 1")
	}

	screen.Index()
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		{cellWith(screen, "o", withFg("red")), cellWith(screen, "t", withFg("red"))},
		{screen.defaultCell(), screen.defaultCell()},
	})

	screen = updateScreen(NewScreen(2, 2), []string{"wo", "ot"}, coloredLines(0))
	screen.ReverseIndex()
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		{screen.defaultCell(), screen.defaultCell()},
		{cellWith(screen, "w", withFg("red")), cellWith(screen, "o", withFg("red"))},
	})
}

func TestPyteLinefeed(t *testing.T) {
	screen := updateScreen(NewScreen(2, 2), []string{"bo", "sh"}, nil)
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

func TestPyteLinefeedMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(3, 27)
	screen.CursorPosition(0, 0)
	screen.LineFeed()
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor moved")
	}
}

func TestPyteTabStops(t *testing.T) {
	screen := NewScreen(10, 10)
	if len(screen.TabStops) != 1 {
		t.Fatalf("expected default tabstop")
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
}

func TestPyteClearTabStops(t *testing.T) {
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
	screen.Cursor.Col = 9
	screen.SetTabStop()
	screen.ClearTabStop(3)
	if len(screen.TabStops) != 0 {
		t.Fatalf("expected cleared tabstops")
	}
}

func TestPyteBackspace(t *testing.T) {
	screen := NewScreen(2, 2)
	screen.Backspace()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected col 0")
	}
	screen.Cursor.Col = 1
	screen.Backspace()
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected backspace")
	}
}

func TestPyteSaveRestoreCursor(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.SaveCursor()
	screen.Cursor.Col = 3
	screen.Cursor.Row = 5
	screen.SaveCursor()
	screen.Cursor.Col = 4
	screen.Cursor.Row = 4
	screen.RestoreCursor()
	if screen.Cursor.Col != 3 || screen.Cursor.Row != 5 {
		t.Fatalf("expected restore cursor")
	}
	screen.RestoreCursor()
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected restore cursor")
	}

	screen = NewScreen(10, 10)
	screen.SetMode([]int{ModeDECAWM, ModeDECOM}, false)
	screen.SaveCursor()
	screen.ResetMode([]int{ModeDECAWM}, false)
	screen.RestoreCursor()
	if !screen.isModeSet(ModeDECAWM) || !screen.isModeSet(ModeDECOM) {
		t.Fatalf("expected modes restored")
	}

	screen = NewScreen(10, 10)
	screen.SelectGraphicRendition([]int{4}, false)
	screen.SaveCursor()
	screen.SelectGraphicRendition([]int{24}, false)
	if screen.Cursor.Attr != screen.defaultAttr() {
		t.Fatalf("expected reset attrs")
	}
	screen.RestoreCursor()
	if !screen.Cursor.Attr.Underline {
		t.Fatalf("expected underline restored")
	}
}

func TestPyteRestoreCursorNoneSaved(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.Cursor.Col = 5
	screen.Cursor.Row = 5
	screen.RestoreCursor()
	if screen.Cursor.Col != 0 || screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor home")
	}
	if screen.isModeSet(ModeDECOM) {
		t.Fatalf("expected DECOM reset")
	}
}

func TestPyteRestoreCursorOutOfBounds(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorPosition(5, 5)
	screen.SaveCursor()
	screen.Resize(3, 3)
	screen.Reset()
	screen.RestoreCursor()
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 2 {
		t.Fatalf("expected clamped cursor")
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
		t.Fatalf("expected origin clamped cursor")
	}
}

func TestPyteInsertDeleteLinesAndChars(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(1))
	screen.InsertLines(1)
	if screen.Display()[0] != "   " || screen.Display()[1] != "sam" {
		t.Fatalf("expected insert lines")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(1))
	screen.DeleteLines(1)
	if screen.Display()[0] != "is " || screen.Display()[1] != "foo" {
		t.Fatalf("expected delete lines")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.InsertCharacters(1)
	if screen.Buffer[0][0].Data != " " || screen.Buffer[0][1].Data != "s" {
		t.Fatalf("expected insert characters")
	}

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.DeleteCharacters(2)
	if screen.Display()[0] != "m  " {
		t.Fatalf("expected delete characters")
	}
}

func TestPyteErase(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.EraseCharacters(2)
	if screen.Display()[0] != "  m" {
		t.Fatalf("expected erase characters")
	}

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(0))
	screen.CursorPosition(1, 3)
	screen.EraseInLine(0, false)
	if screen.Display()[0] != "sa   " {
		t.Fatalf("expected erase in line")
	}

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.CursorPosition(3, 3)
	screen.EraseInDisplay(0, false)
	if screen.Display()[2] != "bu   " {
		t.Fatalf("expected erase in display")
	}
	screen.EraseInDisplay(2, false)
	if screen.Display()[0] != "     " {
		t.Fatalf("expected erase full display")
	}
}

func TestPyteCursorMovement(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorUp(1)
	if screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor top")
	}
	screen.Cursor.Row = 1
	screen.CursorUp(10)
	if screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor top")
	}
	screen.Cursor.Row = 5
	screen.CursorUp(3)
	if screen.Cursor.Row != 2 {
		t.Fatalf("expected cursor row 2")
	}

	screen.Cursor.Row = 9
	screen.CursorDown(1)
	if screen.Cursor.Row != 9 {
		t.Fatalf("expected cursor bottom")
	}
	screen.Cursor.Row = 8
	screen.CursorDown(10)
	if screen.Cursor.Row != 9 {
		t.Fatalf("expected cursor bottom")
	}
	screen.Cursor.Row = 5
	screen.CursorDown(3)
	if screen.Cursor.Row != 8 {
		t.Fatalf("expected cursor row 8")
	}

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

func TestPyteCursorPosition(t *testing.T) {
	screen := NewScreen(10, 10)
	screen.CursorPosition(5, 10)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 9 {
		t.Fatalf("expected row 4 col 9")
	}
	screen.CursorPosition(0, 10)
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 9 {
		t.Fatalf("expected row 0 col 9")
	}
	screen.CursorPosition(100, 5)
	if screen.Cursor.Row != 9 || screen.Cursor.Col != 4 {
		t.Fatalf("expected bounds")
	}
	screen.CursorPosition(5, 100)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 9 {
		t.Fatalf("expected bounds")
	}

	screen.SetMargins(5, 9)
	screen.SetMode([]int{ModeDECOM}, false)
	screen.CursorPosition(0, 0)
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 0 {
		t.Fatalf("expected origin mode position")
	}
	screen.CursorPosition(2, 0)
	if screen.Cursor.Row != 5 {
		t.Fatalf("expected origin offset")
	}
	screen.CursorPosition(10, 0)
	if screen.Cursor.Row != 5 {
		t.Fatalf("expected unchanged")
	}
}

func TestPyteUnicode(t *testing.T) {
	screen := NewScreen(4, 2)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("тест")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "тест" {
		t.Fatalf("unexpected display")
	}
}

func TestPyteAlignmentDisplay(t *testing.T) {
	screen := NewScreen(5, 5)
	screen.SetMode([]int{ModeLNM}, false)
	screen.Draw("a")
	screen.LineFeed()
	screen.LineFeed()
	screen.Draw("b")
	screen.AlignmentDisplay()
	for _, line := range screen.Display() {
		if line != "EEEEE" {
			t.Fatalf("expected alignment display")
		}
	}
}

func TestPyteSetMargins(t *testing.T) {
	screen := NewScreen(10, 10)
	if screen.Margins != nil {
		t.Fatalf("expected nil margins")
	}
	screen.SetMargins(1, 5)
	if screen.Margins == nil || screen.Margins.Top != 0 || screen.Margins.Bottom != 4 {
		t.Fatalf("expected margins set")
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

func TestPyteSetMarginsZero(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetMargins(1, 5)
	if screen.Margins == nil || screen.Margins.Top != 0 || screen.Margins.Bottom != 4 {
		t.Fatalf("expected margins set")
	}
	screen.SetMargins(0, 0)
	if screen.Margins != nil {
		t.Fatalf("expected margins reset")
	}
}

func TestPyteHideCursor(t *testing.T) {
	screen := NewScreen(10, 10)
	if screen.Cursor.Hidden {
		t.Fatalf("expected visible cursor")
	}
	screen.ResetMode([]int{ModeDECTCEM}, false)
	if !screen.Cursor.Hidden {
		t.Fatalf("expected hidden cursor")
	}
	screen.SetMode([]int{ModeDECTCEM}, false)
	if screen.Cursor.Hidden {
		t.Fatalf("expected visible cursor")
	}
}

func TestPyteReportDeviceAttributes(t *testing.T) {
	screen := NewScreen(10, 10)
	acc := []string{}
	screen.WriteProcessInput = func(data string) { acc = append(acc, data) }
	screen.ReportDeviceAttributes(42, false, 0)
	if len(acc) != 0 {
		t.Fatalf("expected no output")
	}
	screen.ReportDeviceAttributes(0, false, 0)
	if acc[len(acc)-1] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response")
	}
}

func TestPyteReportDeviceStatus(t *testing.T) {
	screen := NewScreen(10, 10)
	acc := []string{}
	screen.WriteProcessInput = func(data string) { acc = append(acc, data) }
	screen.ReportDeviceStatus(42, false, 0)
	if len(acc) != 0 {
		t.Fatalf("expected no output")
	}
	screen.ReportDeviceStatus(5, false, 0)
	if acc[len(acc)-1] != ControlCSI+"0n" {
		t.Fatalf("unexpected status")
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

func TestPyteScreenSetIconNameTitle(t *testing.T) {
	screen := NewScreen(10, 1)
	screen.SetIconName("±")
	if screen.IconName != "±" {
		t.Fatalf("expected icon name")
	}
	screen.SetTitle("±")
	if screen.Title != "±" {
		t.Fatalf("expected title")
	}
}

func TestPytePrivateSGRIgnored(t *testing.T) {
	screen := NewHistoryScreen(2, 2, 10)
	stream := NewByteStream(screen, false)
	display := screen.Display()
	mode := screen.Mode
	cursor := screen.Cursor
	if err := stream.Feed([]byte("\x1b[?4m")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if got := screen.Display(); got[0] != display[0] || got[1] != display[1] {
		t.Fatalf("display changed")
	}
	if len(mode) != len(screen.Mode) {
		t.Fatalf("mode changed")
	}
	if cursor != screen.Cursor {
		t.Fatalf("cursor changed")
	}
}
