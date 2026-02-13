package te

import "testing"

// From pyte/tests/test_diff.py::test_mark_whole_screen
func TestPyteTestDiffMarkWholeScreen(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	if screen.Dirty == nil {
		t.Fatalf("expected dirty set")
	}
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.Reset()
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.Resize(24, 130)
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.AlignmentDisplay()
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)
}

// From pyte/tests/test_diff.py::test_mark_single_line
func TestPyteTestDiffMarkSingleLine(t *testing.T) {
	screen := NewDiffScreen(80, 24)

	screen.Dirty = map[int]struct{}{}
	screen.Draw("f")
	assertDirtySet(t, screen.Dirty, []int{screen.Cursor.Row})

	for _, method := range []func(){
		func() { screen.InsertCharacters(1) },
		func() { screen.DeleteCharacters(1) },
		func() { screen.EraseCharacters(1) },
		func() { screen.EraseInLine(0, false) },
	} {
		screen.Dirty = map[int]struct{}{}
		method()
		assertDirtySet(t, screen.Dirty, []int{screen.Cursor.Row})
	}
}

// From pyte/tests/test_diff.py::test_modes
func TestPyteTestDiffModes(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.SetMode([]int{ModeDECSCNM >> 5}, true)
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.ResetMode([]int{ModeDECSCNM >> 5}, true)
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)
}

// From pyte/tests/test_diff.py::test_index
func TestPyteTestDiffIndex(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.Index()
	if len(screen.Dirty) != 0 {
		t.Fatalf("expected no dirty rows")
	}

	screen.CursorToLine(24)
	screen.Index()
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)
}

// From pyte/tests/test_diff.py::test_reverse_index
func TestPyteTestDiffReverseIndex(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.ReverseIndex()
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.CursorToLine(screen.Lines / 2)
	screen.ReverseIndex()
	if len(screen.Dirty) != 0 {
		t.Fatalf("expected no dirty rows")
	}
}

// From pyte/tests/test_diff.py::test_insert_delete_lines
func TestPyteTestDiffInsertDeleteLines(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.CursorToLine(screen.Lines / 2)

	screen.Dirty = map[int]struct{}{}
	screen.InsertLines(1)
	assertDirtyRange(t, screen.Dirty, screen.Cursor.Row, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.DeleteLines(1)
	assertDirtyRange(t, screen.Dirty, screen.Cursor.Row, screen.Lines)
}

// From pyte/tests/test_diff.py::test_erase_in_display
func TestPyteTestDiffEraseInDisplay(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.CursorToLine(screen.Lines / 2)

	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(0, false)
	assertDirtyRange(t, screen.Dirty, screen.Cursor.Row, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(1, false)
	assertDirtyRange(t, screen.Dirty, 0, screen.Cursor.Row+1)

	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(2, false)
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)

	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(3, false)
	assertDirtyRange(t, screen.Dirty, 0, screen.Lines)
}

// From pyte/tests/test_diff.py::test_draw_wrap
func TestPyteTestDiffDrawWrap(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.SetMode([]int{ModeDECAWM}, false)
	for i := 0; i < 80; i++ {
		screen.Draw("g")
	}
	if screen.Cursor.Row != 0 {
		t.Fatalf("expected cursor row 0")
	}
	screen.Dirty = map[int]struct{}{}
	screen.Draw("h")
	assertDirtySet(t, screen.Dirty, []int{0, 1})
}

// From pyte/tests/test_diff.py::test_draw_multiple_chars_wrap
func TestPyteTestDiffDrawMultipleCharsWrap(t *testing.T) {
	screen := NewScreen(5, 2)
	screen.Dirty = map[int]struct{}{}
	screen.Draw("1234567890")
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected cursor row 1")
	}
	assertDirtySet(t, screen.Dirty, []int{0, 1})
}
