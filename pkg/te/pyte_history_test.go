package te

import (
	"fmt"
	"testing"
)

func pyteHistoryStream(screen *HistoryScreen) *Stream {
	stream := NewStream(screen, false)
	stream.escapeOverrides = map[rune]func(){
		'N': func() { screen.NextPage() },
		'P': func() { screen.PrevPage() },
	}
	return stream
}

func pyteHistoryChars(lines [][]Cell, columns int) []string {
	out := make([]string, len(lines))
	for y := range lines {
		row := lines[y]
		buf := make([]rune, columns)
		for x := 0; x < columns; x++ {
			if x < len(row) && row[x].Data != "" {
				buf[x] = []rune(row[x].Data)[0]
			} else {
				buf[x] = ' '
			}
		}
		out[y] = string(buf)
	}
	return out
}

// From pyte/tests/test_history.py::test_index
func TestPyteTestHistoryIndex(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		if idx != screen.Lines-1 {
			screen.LineFeed()
		}
	}
	if len(screen.history.Top.items) != 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected empty history")
	}
	line := screen.Buffer[0]
	screen.Index()
	if len(screen.history.Top.items) == 0 {
		t.Fatalf("expected top history")
	}
	if pyteHistoryChars([][]Cell{screen.history.Top.items[len(screen.history.Top.items)-1]}, screen.Columns)[0] != pyteHistoryChars([][]Cell{line}, screen.Columns)[0] {
		t.Fatalf("unexpected history line")
	}

	line = screen.Buffer[0]
	screen.Index()
	if len(screen.history.Top.items) != 2 {
		t.Fatalf("expected 2 history entries")
	}
	if pyteHistoryChars([][]Cell{screen.history.Top.items[len(screen.history.Top.items)-1]}, screen.Columns)[0] != pyteHistoryChars([][]Cell{line}, screen.Columns)[0] {
		t.Fatalf("unexpected history line")
	}

	for i := 0; i < screen.history.Size*2; i++ {
		screen.Index()
	}
	if len(screen.history.Top.items) != 50 {
		t.Fatalf("expected history rotation")
	}
}

// From pyte/tests/test_history.py::test_reverse_index
func TestPyteTestHistoryReverseIndex(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < len(screen.Buffer); idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		if idx != len(screen.Buffer)-1 {
			screen.LineFeed()
		}
	}
	if len(screen.history.Top.items) != 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected empty history")
	}
	screen.CursorPosition(0, 0)
	line := screen.Buffer[screen.Lines-1]
	screen.ReverseIndex()
	if len(screen.history.Bottom.items) == 0 {
		t.Fatalf("expected bottom history")
	}
	if pyteHistoryChars([][]Cell{screen.history.Bottom.items[0]}, screen.Columns)[0] != pyteHistoryChars([][]Cell{line}, screen.Columns)[0] {
		t.Fatalf("unexpected history line")
	}

	line = screen.Buffer[screen.Lines-1]
	screen.ReverseIndex()
	if len(screen.history.Bottom.items) != 2 {
		t.Fatalf("expected 2 history entries")
	}
	if pyteHistoryChars([][]Cell{screen.history.Bottom.items[1]}, screen.Columns)[0] != pyteHistoryChars([][]Cell{line}, screen.Columns)[0] {
		t.Fatalf("unexpected history line")
	}

	iterations := 1
	for i := 0; i < screen.Lines; i++ {
		iterations *= len(screen.Buffer)
	}
	for i := 0; i < iterations; i++ {
		screen.ReverseIndex()
	}
	if len(screen.history.Bottom.items) != 50 {
		t.Fatalf("expected history rotation")
	}
}

// From pyte/tests/test_history.py::test_prev_page
func TestPyteTestHistoryPrevPage(t *testing.T) {
	screen := NewHistoryScreen(4, 4, 40)
	screen.SetMode([]int{ModeLNM}, false)
	if screen.history.Position != 40 {
		t.Fatalf("expected position 40")
	}
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 40 {
		t.Fatalf("expected position 40")
	}
	assertDisplay(t, screen.Screen, []string{"37  ", "38  ", "39  ", "    "})
	if got := pyteHistoryChars(screen.history.Top.items, screen.Columns); len(got) < 4 || got[len(got)-4] != "33  " || got[len(got)-1] != "36  " {
		t.Fatalf("unexpected top history")
	}

	screen.PrevPage()
	if screen.history.Position != 38 {
		t.Fatalf("expected position 38")
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"35  ", "36  ", "37  ", "38  "})
	if got := pyteHistoryChars(screen.history.Top.items, screen.Columns); len(got) < 4 || got[len(got)-4] != "31  " || got[len(got)-1] != "34  " {
		t.Fatalf("unexpected top history after prev")
	}
	if len(screen.history.Bottom.items) != 2 {
		t.Fatalf("expected bottom history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "39  " || got[1] != "    " {
		t.Fatalf("unexpected bottom history")
	}

	screen.PrevPage()
	if screen.history.Position != 36 {
		t.Fatalf("expected position 36")
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"33  ", "34  ", "35  ", "36  "})
	if len(screen.history.Bottom.items) != 4 {
		t.Fatalf("expected bottom history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "37  " || got[3] != "    " {
		t.Fatalf("unexpected bottom history")
	}

	screen = NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 50 {
		t.Fatalf("expected position 50")
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})
	screen.PrevPage()
	if screen.history.Position != 47 {
		t.Fatalf("expected position 47")
	}
	assertDisplay(t, screen.Screen, []string{"43   ", "44   ", "45   ", "46   ", "47   "})
	if len(screen.history.Bottom.items) != 3 {
		t.Fatalf("expected bottom history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "48   " || got[2] != "     " {
		t.Fatalf("unexpected bottom history")
	}

	screen = NewHistoryScreenWithRatio(4, 4, 40, 0.75)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 40 {
		t.Fatalf("expected position 40")
	}
	assertDisplay(t, screen.Screen, []string{"37  ", "38  ", "39  ", "    "})
	screen.PrevPage()
	if screen.history.Position != 37 {
		t.Fatalf("expected position 37")
	}
	assertDisplay(t, screen.Screen, []string{"34  ", "35  ", "36  ", "37  "})
	if len(screen.history.Bottom.items) != 3 {
		t.Fatalf("expected bottom history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "38  " || got[2] != "    " {
		t.Fatalf("unexpected bottom history")
	}

	screen = NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 50 {
		t.Fatalf("expected position 50")
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})
	screen.CursorToLine(screen.Lines / 2)
	for screen.history.Position > screen.Lines {
		screen.PrevPage()
	}
	if screen.history.Position != screen.Lines {
		t.Fatalf("expected position %d", screen.Lines)
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"1    ", "2    ", "3    ", "4    ", "5    "})
	for screen.history.Position < screen.history.Size {
		screen.NextPage()
	}
	if screen.history.Position != screen.history.Size {
		t.Fatalf("expected position %d", screen.history.Size)
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})

	screen = NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 50 {
		t.Fatalf("expected position 50")
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})
	screen.CursorToLine(screen.Lines/2 - 2)
	for screen.history.Position > screen.Lines {
		screen.PrevPage()
	}
	if screen.history.Position != screen.Lines {
		t.Fatalf("expected position %d", screen.Lines)
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"1    ", "2    ", "3    ", "4    ", "5    "})
	for screen.history.Position < screen.history.Size {
		screen.NextPage()
	}
	if screen.history.Position != screen.history.Size {
		t.Fatalf("expected position %d", screen.history.Size)
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})
}

// From pyte/tests/test_history.py::test_next_page
func TestPyteTestHistoryNextPage(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*5; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 50 {
		t.Fatalf("expected position 50")
	}
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "     "})

	screen.PrevPage()
	screen.NextPage()
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 50 {
		t.Fatalf("expected position 50")
	}
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "     "})

	screen.PrevPage()
	screen.PrevPage()
	screen.NextPage()
	if screen.history.Position != 47 {
		t.Fatalf("expected position 47")
	}
	if len(screen.history.Top.items) == 0 {
		t.Fatalf("expected top history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "23   " || got[2] != "     " {
		t.Fatalf("unexpected bottom history")
	}
	assertDisplay(t, screen.Screen, []string{"18   ", "19   ", "20   ", "21   ", "22   "})

	screen.PrevPage()
	screen.PrevPage()
	screen.NextPage()
	screen.NextPage()
	if screen.history.Position != 47 {
		t.Fatalf("expected position 47")
	}
	if len(screen.Buffer) != screen.Lines {
		t.Fatalf("expected buffer lines")
	}
	assertDisplay(t, screen.Screen, []string{"18   ", "19   ", "20   ", "21   ", "22   "})
}

// From pyte/tests/test_history.py::test_ensure_width
func TestPyteTestHistoryEnsureWidth(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	stream := pyteHistoryStream(screen)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed(fmt.Sprintf("%04d\n", idx))
	}
	assertDisplay(t, screen.Screen, []string{"0021 ", "0022 ", "0023 ", "0024 ", "     "})

	screen.Resize(5, 3)
	if err := stream.Feed(ControlESC + "P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	for _, line := range screen.history.Bottom.items {
		if len(line) > 3 {
			t.Fatalf("expected truncated history")
		}
	}
	assertDisplay(t, screen.Screen, []string{"001", "001", "002", "002", "002"})
}

// From pyte/tests/test_history.py::test_not_enough_lines
func TestPyteTestHistoryNotEnoughLines(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 6)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	if screen.history.Position != 6 {
		t.Fatalf("expected position 6")
	}
	assertDisplay(t, screen.Screen, []string{"1    ", "2    ", "3    ", "4    ", "     "})

	screen.PrevPage()
	if len(screen.history.Top.items) != 0 {
		t.Fatalf("expected no top history")
	}
	if len(screen.history.Bottom.items) != 1 {
		t.Fatalf("expected one bottom history")
	}
	if got := pyteHistoryChars(screen.history.Bottom.items, screen.Columns); got[0] != "     " {
		t.Fatalf("unexpected bottom history")
	}
	assertDisplay(t, screen.Screen, []string{"0    ", "1    ", "2    ", "3    ", "4    "})

	screen.NextPage()
	if len(screen.history.Top.items) == 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected top history only")
	}
	assertDisplay(t, screen.Screen, []string{"1    ", "2    ", "3    ", "4    ", "     "})
}

// From pyte/tests/test_history.py::test_draw
func TestPyteTestHistoryDraw(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	stream := pyteHistoryStream(screen)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed(fmt.Sprintf("%d\n", idx))
	}
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "     "})

	if err := stream.Feed(ControlESC + "P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlESC + "P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlESC + "N"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("x"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "x    "})
}

// From pyte/tests/test_history.py::test_cursor_is_hidden
func TestPyteTestHistoryCursorIsHidden(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	stream := pyteHistoryStream(screen)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed(fmt.Sprintf("%d\n", idx))
	}
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
	if err := stream.Feed(ControlESC + "P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
	if err := stream.Feed(ControlESC + "P"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
	if err := stream.Feed(ControlESC + "N"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
	if err := stream.Feed(ControlESC + "N"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
}

// From pyte/tests/test_history.py::test_erase_in_display
func TestPyteTestHistoryEraseInDisplay(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 6)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(fmt.Sprintf("%d", idx))
		screen.LineFeed()
	}
	screen.PrevPage()
	screen.EraseInDisplay(3, false)
	if len(screen.history.Top.items) != 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected history reset")
	}
}
