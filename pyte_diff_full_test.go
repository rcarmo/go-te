package te

import "testing"

func TestPyteDiffDrawWrap(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.SetMode([]int{ModeDECAWM}, false)
	screen.Dirty = map[int]struct{}{}
	for i := 0; i < 80; i++ {
		screen.Draw("g")
	}
	screen.Dirty = map[int]struct{}{}
	screen.Draw("h")
	if len(screen.Dirty) != 2 {
		t.Fatalf("expected dirty lines 0 and 1")
	}
	if _, ok := screen.Dirty[0]; !ok {
		t.Fatalf("expected dirty line 0")
	}
	if _, ok := screen.Dirty[1]; !ok {
		t.Fatalf("expected dirty line 1")
	}
}

func TestPyteDiffDrawMultipleCharsWrap(t *testing.T) {
	screen := NewDiffScreen(5, 2)
	screen.Dirty = map[int]struct{}{}
	screen.Draw("1234567890")
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected cursor row 1")
	}
	if len(screen.Dirty) != 2 {
		t.Fatalf("expected dirty lines")
	}
}

func TestPyteDiffInsertDeleteLines(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Cursor.Row = screen.Lines / 2
	screen.Dirty = map[int]struct{}{}
	screen.InsertLines(1)
	if len(screen.Dirty) == 0 {
		t.Fatalf("expected dirty lines")
	}
	screen.Dirty = map[int]struct{}{}
	screen.DeleteLines(1)
	if len(screen.Dirty) == 0 {
		t.Fatalf("expected dirty lines")
	}
}

func TestPyteDiffEraseInDisplay(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Cursor.Row = screen.Lines / 2
	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(0, false)
	if len(screen.Dirty) == 0 {
		t.Fatalf("expected dirty lines")
	}
	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(1, false)
	if len(screen.Dirty) == 0 {
		t.Fatalf("expected dirty lines")
	}
	screen.Dirty = map[int]struct{}{}
	screen.EraseInDisplay(2, false)
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected full dirty")
	}
}
