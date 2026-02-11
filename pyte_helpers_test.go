package te

import "testing"

type attrOpt func(*Attr)

type cellOpt func(*Cell)

func updateScreen(screen *Screen, lines []string, colored map[int]struct{}) *Screen {
	for y, line := range lines {
		for x, char := range line {
			attr := screen.defaultAttr()
			if colored != nil {
				if _, ok := colored[y]; ok {
					attr.Fg = colorFromName("red", ColorANSI16, 1)
				}
			}
			screen.Buffer[y][x] = Cell{Data: string(char), Attr: attr}
		}
	}
	return screen
}

func coloredLines(lines ...int) map[int]struct{} {
	set := make(map[int]struct{}, len(lines))
	for _, line := range lines {
		set[line] = struct{}{}
	}
	return set
}

func cellWith(screen *Screen, data string, opts ...attrOpt) Cell {
	attr := screen.defaultAttr()
	for _, opt := range opts {
		opt(&attr)
	}
	return Cell{Data: data, Attr: attr}
}

func withFg(name string) attrOpt {
	return func(attr *Attr) {
		attr.Fg = parseColor(name, true)
	}
}

func withBg(name string) attrOpt {
	return func(attr *Attr) {
		attr.Bg = parseColor(name, false)
	}
}

func withBold(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Bold = value
	}
}

func withItalics(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Italics = value
	}
}

func withUnderline(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Underline = value
	}
}

func withStrikethrough(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Strikethrough = value
	}
}

func withReverse(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Reverse = value
	}
}

func withBlink(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Blink = value
	}
}

func parseColor(name string, fg bool) Color {
	if name == "default" {
		return Color{Name: name, Mode: ColorDefault}
	}
	if fg {
		for code, value := range fgANSI {
			if value == name {
				return Color{Name: name, Mode: ColorANSI16, Index: uint8(code - 30)}
			}
		}
		for code, value := range fgAixterm {
			if value == name {
				return Color{Name: name, Mode: ColorANSI16, Index: uint8(code - 90 + 8)}
			}
		}
	} else {
		for code, value := range bgANSI {
			if value == name {
				return Color{Name: name, Mode: ColorANSI16, Index: uint8(code - 40)}
			}
		}
		for code, value := range bgAixterm {
			if value == name {
				return Color{Name: name, Mode: ColorANSI16, Index: uint8(code - 100 + 8)}
			}
		}
	}
	return Color{Name: name}
}

func lineCells(screen *Screen, text string, opts ...attrOpt) []Cell {
	cells := make([]Cell, screen.Columns)
	attr := screen.defaultAttr()
	for _, opt := range opts {
		opt(&attr)
	}
	for i := 0; i < screen.Columns; i++ {
		cells[i] = Cell{Data: " ", Attr: attr}
	}
	for i, r := range text {
		if i >= screen.Columns {
			break
		}
		cells[i] = Cell{Data: string(r), Attr: attr}
	}
	return cells
}

func assertDisplay(t *testing.T, screen *Screen, expected []string) {
	got := screen.Display()
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("line %d expected %q, got %q", i, expected[i], got[i])
		}
	}
}

func assertCellsEqual(t *testing.T, got [][]Cell, want [][]Cell) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d rows, got %d", len(want), len(got))
	}
	for row := range want {
		if len(got[row]) != len(want[row]) {
			t.Fatalf("row %d expected %d cols, got %d", row, len(want[row]), len(got[row]))
		}
		for col := range want[row] {
			assertCellEqual(t, got[row][col], want[row][col], row, col)
		}
	}
}

func assertCellEqual(t *testing.T, got Cell, want Cell, row, col int) {
	t.Helper()
	if got.Data != want.Data {
		t.Fatalf("cell(%d,%d) expected data %q, got %q", row, col, want.Data, got.Data)
	}
	if got.Attr.Fg.Name != want.Attr.Fg.Name {
		t.Fatalf("cell(%d,%d) expected fg %q, got %q", row, col, want.Attr.Fg.Name, got.Attr.Fg.Name)
	}
	if got.Attr.Bg.Name != want.Attr.Bg.Name {
		t.Fatalf("cell(%d,%d) expected bg %q, got %q", row, col, want.Attr.Bg.Name, got.Attr.Bg.Name)
	}
	if got.Attr.Bold != want.Attr.Bold || got.Attr.Italics != want.Attr.Italics ||
		got.Attr.Underline != want.Attr.Underline || got.Attr.Strikethrough != want.Attr.Strikethrough ||
		got.Attr.Reverse != want.Attr.Reverse || got.Attr.Blink != want.Attr.Blink {
		t.Fatalf("cell(%d,%d) attr mismatch", row, col)
	}
}
