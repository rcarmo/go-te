package te

import "testing"

func TestDiffMarkWholeScreen(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty set on init")
	}
	screen.Dirty = map[int]struct{}{}
	screen.Reset()
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty set on reset")
	}
	screen.Dirty = map[int]struct{}{}
	screen.Resize(24, 130)
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty set on resize")
	}
	screen.Dirty = map[int]struct{}{}
	screen.AlignmentDisplay()
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty set on alignment")
	}
}

func TestDiffMarkSingleLine(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.Draw("f")
	if len(screen.Dirty) != 1 {
		t.Fatalf("expected single dirty line")
	}
}

func TestDiffModes(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.SetMode([]int{ModeDECSCNM >> 5}, true)
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty on reverse mode")
	}
	screen.Dirty = map[int]struct{}{}
	screen.ResetMode([]int{ModeDECSCNM >> 5}, true)
	if len(screen.Dirty) != screen.Lines {
		t.Fatalf("expected dirty on reverse reset")
	}
}

func TestDiffIndexReverseIndex(t *testing.T) {
	screen := NewDiffScreen(80, 24)
	screen.Dirty = map[int]struct{}{}
	screen.Cursor.Row = screen.Lines - 1
	screen.Index()
	if len(screen.Dirty) == 0 {
		t.Fatalf("expected dirty on index at bottom")
	}
	screen.Dirty = map[int]struct{}{}
	screen.Cursor.Row = screen.Lines / 2
	screen.ReverseIndex()
	if len(screen.Dirty) != 0 {
		t.Fatalf("expected no dirty outside top margin")
	}
}
