package te

import "testing"

func assertDirtySet(t *testing.T, dirty map[int]struct{}, expected []int) {
	t.Helper()
	if len(dirty) != len(expected) {
		t.Fatalf("expected %d dirty rows, got %d", len(expected), len(dirty))
	}
	for _, row := range expected {
		if _, ok := dirty[row]; !ok {
			t.Fatalf("expected dirty row %d", row)
		}
	}
}

func assertDirtyRange(t *testing.T, dirty map[int]struct{}, start, end int) {
	t.Helper()
	if len(dirty) != end-start {
		t.Fatalf("expected %d dirty rows, got %d", end-start, len(dirty))
	}
	for row := start; row < end; row++ {
		if _, ok := dirty[row]; !ok {
			t.Fatalf("expected dirty row %d", row)
		}
	}
}

func assertDisplay(t *testing.T, screen *Screen, expected []string) {
	t.Helper()
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

func updateScreen(screen *Screen, lines []string, colored []int) *Screen {
	coloredSet := map[int]struct{}{}
	for _, row := range colored {
		coloredSet[row] = struct{}{}
	}
	for y, line := range lines {
		for x, char := range line {
			attr := screen.defaultAttr()
			if _, ok := coloredSet[y]; ok {
				attr.Fg = colorFromName("red", ColorANSI16, 1)
			}
			screen.Buffer[y][x] = Cell{Data: string(char), Attr: attr}
		}
	}
	return screen
}

func tolist(screen *Screen) [][]Cell {
	out := make([][]Cell, screen.Lines)
	for y := 0; y < screen.Lines; y++ {
		row := make([]Cell, screen.Columns)
		copy(row, screen.Buffer[y])
		out[y] = row
	}
	return out
}

type attrOpt func(*Attr)

func cellWith(screen *Screen, data string, opts ...attrOpt) Cell {
	attr := screen.defaultAttr()
	for _, opt := range opts {
		opt(&attr)
	}
	return Cell{Data: data, Attr: attr}
}

func withFg(name string) attrOpt {
	return func(attr *Attr) {
		attr.Fg = colorFromName(name, ColorANSI16, 0)
		if name == "default" {
			attr.Fg = Color{Name: name, Mode: ColorDefault}
		}
	}
}

func withBg(name string) attrOpt {
	return func(attr *Attr) {
		attr.Bg = colorFromName(name, ColorANSI16, 0)
		if name == "default" {
			attr.Bg = Color{Name: name, Mode: ColorDefault}
		}
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

func withConceal(value bool) attrOpt {
	return func(attr *Attr) {
		attr.Conceal = value
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
		got.Attr.Reverse != want.Attr.Reverse || got.Attr.Blink != want.Attr.Blink || got.Attr.Conceal != want.Attr.Conceal {
		t.Fatalf("cell(%d,%d) attr mismatch", row, col)
	}
}
