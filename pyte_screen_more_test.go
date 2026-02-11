package te

import "testing"

func TestPyteCursorBackLastColumn(t *testing.T) {
	screen := NewScreen(13, 1)
	screen.Draw("Hello, world!")
	if screen.Cursor.Col != screen.Columns {
		t.Fatalf("expected cursor at end")
	}
	screen.CursorBack(5)
	if screen.Cursor.Col != (screen.Columns-1)-5 {
		t.Fatalf("unexpected cursor")
	}
}

func TestPyteInsertLinesFull(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(1))
	screen.InsertLines(1)
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		lineCells(screen, "   "),
		lineCells(screen, "sam"),
		lineCells(screen, "is ", withFg("red")),
	})

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(1))
	screen.InsertLines(2)
	assertCellsEqual(t, screen.Buffer, [][]Cell{
		lineCells(screen, "   "),
		lineCells(screen, "   "),
		lineCells(screen, "sam"),
	})

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.InsertLines(1)
	if screen.Display()[0] != "sam" || screen.Display()[1] != "   " || screen.Display()[2] != "is " {
		t.Fatalf("unexpected display")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(1, 3)
	screen.Cursor.Row = 1
	screen.InsertLines(1)
	if screen.Display()[3] != "bar" {
		t.Fatalf("unexpected display")
	}

	screen.InsertLines(2)
	if screen.Display()[1] != "   " || screen.Display()[2] != "   " {
		t.Fatalf("unexpected insert lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 1
	screen.InsertLines(20)
	if screen.Display()[1] != "   " || screen.Display()[4] != "baz" {
		t.Fatalf("unexpected insert lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(2, 4)
	screen.InsertLines(5)
	if screen.Display()[0] != "sam" || screen.Display()[4] != "baz" {
		t.Fatalf("unexpected no-op insert")
	}
}

func TestPyteDeleteLinesFull(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(1))
	screen.DeleteLines(1)
	if screen.Display()[0] != "is " || screen.Display()[1] != "foo" {
		t.Fatalf("unexpected delete lines")
	}

	screen.DeleteLines(0)
	if screen.Display()[0] != "foo" {
		t.Fatalf("unexpected delete lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(1)
	if screen.Display()[1] != "foo" || screen.Display()[3] != "   " {
		t.Fatalf("unexpected delete lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(2)
	if screen.Display()[1] != "bar" {
		t.Fatalf("unexpected delete lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, nil)
	screen.SetMargins(1, 4)
	screen.Cursor.Row = 1
	screen.DeleteLines(5)
	if screen.Display()[1] != "   " || screen.Display()[4] != "baz" {
		t.Fatalf("unexpected delete lines")
	}

	screen = updateScreen(NewScreen(3, 5), []string{"sam", "is ", "foo", "bar", "baz"}, coloredLines(2, 3))
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 0
	screen.DeleteLines(5)
	if screen.Display()[0] != "sam" || screen.Display()[4] != "baz" {
		t.Fatalf("unexpected delete lines")
	}
}

func TestPyteInsertCharactersFull(t *testing.T) {
	screen := updateScreen(NewScreen(3, 4), []string{"sam", "is ", "foo", "bar"}, coloredLines(0))
	cursor := screen.Cursor
	screen.InsertCharacters(2)
	if screen.Cursor != cursor {
		t.Fatalf("expected cursor unchanged")
	}
	assertCellsEqual(t, screen.Buffer[:1], [][]Cell{{
		screen.defaultCell(),
		screen.defaultCell(),
		cellWith(screen, "s", withFg("red")),
	}})

	screen.Cursor.Row = 2
	screen.Cursor.Col = 1
	screen.InsertCharacters(1)
	assertCellsEqual(t, screen.Buffer[2:3], [][]Cell{{
		cellWith(screen, "f"),
		screen.defaultCell(),
		cellWith(screen, "o"),
	}})

	screen.Cursor.Row = 3
	screen.Cursor.Col = 1
	screen.InsertCharacters(10)
	assertCellsEqual(t, screen.Buffer[3:4], [][]Cell{{
		cellWith(screen, "b"),
		screen.defaultCell(),
		screen.defaultCell(),
	}})

	screen = updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.CursorPosition(0, 0)
	screen.InsertCharacters(0)
	assertCellsEqual(t, screen.Buffer[:1], [][]Cell{{
		screen.defaultCell(),
		cellWith(screen, "s", withFg("red")),
		cellWith(screen, "a", withFg("red")),
	}})
}

func TestPyteDeleteCharactersFull(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.DeleteCharacters(2)
	if screen.Display()[0] != "m  " {
		t.Fatalf("unexpected delete chars")
	}

	screen.Cursor.Row = 2
	screen.Cursor.Col = 2
	screen.DeleteCharacters(1)
	if screen.Display()[2] != "fo " {
		t.Fatalf("unexpected delete chars")
	}

	screen.Cursor.Row = 1
	screen.Cursor.Col = 1
	screen.DeleteCharacters(0)
	if screen.Display()[1] != "i  " {
		t.Fatalf("unexpected delete chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.Cursor.Col = 1
	screen.DeleteCharacters(3)
	if screen.Display()[0] != "15   " {
		t.Fatalf("unexpected delete chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.Cursor.Col = 2
	screen.DeleteCharacters(10)
	if screen.Display()[0] != "12   " {
		t.Fatalf("unexpected delete chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.DeleteCharacters(4)
	if screen.Display()[0] != "5    " {
		t.Fatalf("unexpected delete chars")
	}
}

func TestPyteEraseCharactersFull(t *testing.T) {
	screen := updateScreen(NewScreen(3, 3), []string{"sam", "is ", "foo"}, coloredLines(0))
	screen.EraseCharacters(2)
	if screen.Display()[0] != "  m" {
		t.Fatalf("unexpected erase chars")
	}

	screen.Cursor.Row = 2
	screen.Cursor.Col = 2
	screen.EraseCharacters(1)
	if screen.Display()[2] != "fo " {
		t.Fatalf("unexpected erase chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.Cursor.Col = 1
	screen.EraseCharacters(3)
	if screen.Display()[0] != "1   5" {
		t.Fatalf("unexpected erase chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.Cursor.Col = 2
	screen.EraseCharacters(10)
	if screen.Display()[0] != "12   " {
		t.Fatalf("unexpected erase chars")
	}

	screen = updateScreen(NewScreen(5, 1), []string{"12345"}, coloredLines(0))
	screen.EraseCharacters(4)
	if screen.Display()[0] != "    5" {
		t.Fatalf("unexpected erase chars")
	}
}

func TestPyteEraseInLineModes(t *testing.T) {
	screen := updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(0))
	screen.CursorPosition(1, 3)
	screen.EraseInLine(0, false)
	if screen.Display()[0] != "sa   " {
		t.Fatalf("unexpected erase line")
	}

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(0))
	screen.EraseInLine(1, false)
	if screen.Display()[0] != "    i" {
		t.Fatalf("unexpected erase line")
	}

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(0))
	screen.EraseInLine(2, false)
	if screen.Display()[0] != "     " {
		t.Fatalf("unexpected erase line")
	}
}

func TestPyteEraseInDisplayModes(t *testing.T) {
	screen := updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.CursorPosition(3, 3)
	screen.EraseInDisplay(0, false)
	if screen.Display()[2] != "bu   " {
		t.Fatalf("unexpected erase display")
	}

	screen = updateScreen(screen, []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.EraseInDisplay(1, false)
	if screen.Display()[0] != "     " {
		t.Fatalf("unexpected erase display")
	}

	screen.EraseInDisplay(2, false)
	if screen.Display()[0] != "     " {
		t.Fatalf("unexpected erase display")
	}

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.EraseInDisplay(3, false)
	if screen.Display()[0] != "     " {
		t.Fatalf("unexpected erase display")
	}

	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.EraseInDisplay(3, true)
	if screen.Display()[0] != "     " {
		t.Fatalf("expected private erase")
	}
	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.EraseInDisplay(3, false, 0)
	if screen.Display()[0] != "     " {
		t.Fatalf("unexpected erase display")
	}
	screen = updateScreen(NewScreen(5, 5), []string{"sam i", "s foo", "but a", "re yo", "u?   "}, coloredLines(2, 3))
	screen.EraseInDisplay(3, true, 0)
	if screen.Display()[0] != "     " {
		t.Fatalf("expected private erase")
	}
}
