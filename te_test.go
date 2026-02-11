package te

import "testing"

func TestStreamAutowrap(t *testing.T) {
	screen := NewScreen(4, 2)
	stream := NewStream(screen, false)
	if err := stream.FeedString("abcdE"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "abcd")
	assertEqual(t, lines[1], "E   ")
}

func TestStreamLineFeedAndCarriageReturn(t *testing.T) {
	screen := NewScreen(5, 2)
	stream := NewStream(screen, false)
	if err := stream.FeedString("ab\ncd"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "ab   ")
	assertEqual(t, lines[1], "  cd ")
}

func TestStreamNewlineMode(t *testing.T) {
	screen := NewScreen(5, 2)
	stream := NewStream(screen, false)
	if err := stream.FeedString("\x1b[20hab\ncd"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "ab   ")
	assertEqual(t, lines[1], "cd   ")
}

func TestStreamCursorMovementAndErase(t *testing.T) {
	screen := NewScreen(6, 1)
	stream := NewStream(screen, false)
	if err := stream.FeedString("hello\x1b[2DX"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "helXo ")

	screen = NewScreen(6, 1)
	stream = NewStream(screen, false)
	if err := stream.FeedString("hello\x1b[3G\x1b[K"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines = screen.Display()
	assertEqual(t, lines[0], "he    ")
}

func TestStreamScrollAndHistory(t *testing.T) {
	history := NewHistoryScreen(6, 2, 4)
	stream := NewStream(history, false)
	if err := stream.FeedString("one\r\nTwo\r\nThree"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := history.Display()
	assertEqual(t, lines[0], "Two   ")
	assertEqual(t, lines[1], "Three ")
	if history.Scrollback() != 1 {
		t.Fatalf("expected history length 1, got %d", history.Scrollback())
	}
	if len(history.History()) != 1 || stringLine(history.History()[0]) != "one   " {
		t.Fatalf("unexpected history content")
	}
}

func TestDiffScreenDirtyTracking(t *testing.T) {
	diff := NewDiffScreen(4, 2)
	stream := NewStream(diff, false)
	if err := stream.FeedString("hi"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	if len(diff.DirtyLines()) == 0 {
		t.Fatalf("expected dirty lines")
	}
	diff.ClearDirty()
	if err := stream.FeedString("\x1b[2K"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	if len(diff.DirtyLines()) != 1 || diff.DirtyLines()[0] != 0 {
		t.Fatalf("expected dirty line 0")
	}
}

func TestSGRAttributes(t *testing.T) {
	screen := NewScreen(3, 1)
	stream := NewStream(screen, false)
	if err := stream.FeedString("A\x1b[31mB\x1b[0mC"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Lines()
	if lines[0][0].Attr.Fg.Mode != ColorDefault {
		t.Fatalf("expected default foreground for first cell")
	}
	if lines[0][1].Attr.Fg.Mode != ColorANSI16 || lines[0][1].Attr.Fg.Index != 1 {
		t.Fatalf("expected red foreground for second cell")
	}
	if lines[0][2].Attr.Fg.Mode != ColorDefault {
		t.Fatalf("expected reset foreground for third cell")
	}
}

func TestScrollRegion(t *testing.T) {
	screen := NewScreen(4, 4)
	stream := NewStream(screen, false)
	if err := stream.FeedString("\x1b[2;3rA\r\nB\r\nC"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "    ")
	assertEqual(t, lines[1], "B   ")
	assertEqual(t, lines[2], "C   ")
	assertEqual(t, lines[3], "    ")
}

func TestInsertAndDeleteChars(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewStream(screen, false)
	if err := stream.FeedString("abcde\x1b[3G\x1b[@Z"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "abZcd")

	screen = NewScreen(5, 1)
	stream = NewStream(screen, false)
	if err := stream.FeedString("abcde\x1b[2G\x1b[2P"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines = screen.Display()
	assertEqual(t, lines[0], "ade  ")
}

func TestTabStops(t *testing.T) {
	screen := NewScreen(8, 1)
	stream := NewStream(screen, false)
	if err := stream.FeedString("\x1b[4G\x1bH\x1b[1G\tX"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "   X    ")

	screen = NewScreen(8, 1)
	stream = NewStream(screen, false)
	if err := stream.FeedString("\x1b[4G\x1bH\x1b[4G\x1b[g\x1b[1G\tX"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines = screen.Display()
	assertEqual(t, lines[0], "       X")
}

func TestAlternateBuffer(t *testing.T) {
	screen := NewScreen(4, 2)
	stream := NewStream(screen, false)
	if err := stream.FeedString("AB\x1b[?1049hZZ\x1b[?1049l"); err != nil {
		t.Fatalf("FeedString: %v", err)
	}
	lines := screen.Display()
	assertEqual(t, lines[0], "AB  ")
	assertEqual(t, lines[1], "    ")
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func stringLine(line []Cell) string {
	out := make([]rune, len(line))
	for i, cell := range line {
		out[i] = cell.Ch
	}
	return string(out)
}
