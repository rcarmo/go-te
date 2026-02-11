package te

import (
	"strconv"
	"testing"
)

func historyLines(lines [][]Cell, columns int) []string {
	out := make([]string, len(lines))
	for i := range lines {
		line := lines[i]
		text := make([]rune, 0, columns)
		for x := 0; x < columns && x < len(line); x++ {
			if line[x].Data == "" {
				text = append(text, ' ')
			} else {
				text = append(text, []rune(line[x].Data)...)
			}
		}
		out[i] = string(text)
	}
	return out
}

func TestPyteHistoryIndex(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(string(rune('0' + idx)))
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
	if historyLineString(screen.history.Top.items[len(screen.history.Top.items)-1]) != historyLineString(line) {
		t.Fatalf("unexpected history line")
	}

	line = screen.Buffer[0]
	screen.Index()
	if len(screen.history.Top.items) != 2 {
		t.Fatalf("expected 2 history entries")
	}
	if historyLineString(screen.history.Top.items[len(screen.history.Top.items)-1]) != historyLineString(line) {
		t.Fatalf("unexpected history line")
	}
}

func TestPyteHistoryReverseIndex(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(string(rune('0' + idx)))
		if idx != screen.Lines-1 {
			screen.LineFeed()
		}
	}
	screen.CursorPosition(0, 0)
	line := screen.Buffer[screen.Lines-1]
	screen.ReverseIndex()
	if len(screen.history.Bottom.items) == 0 {
		t.Fatalf("expected bottom history")
	}
	if historyLineString(screen.history.Bottom.items[0]) != historyLineString(line) {
		t.Fatalf("unexpected history line")
	}

	line = screen.Buffer[screen.Lines-1]
	screen.ReverseIndex()
	if len(screen.history.Bottom.items) != 2 {
		t.Fatalf("expected 2 history entries")
	}
	if historyLineString(screen.history.Bottom.items[1]) != historyLineString(line) {
		t.Fatalf("unexpected history line")
	}
}

func TestPyteHistoryPrevPage(t *testing.T) {
	screen := NewHistoryScreen(4, 4, 40)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(strconv.Itoa(idx))
		screen.LineFeed()
	}
	assertDisplay(t, screen.Screen, []string{"37  ", "38  ", "39  ", "    "})

	screen.PrevPage()
	if screen.history.Position != 38 {
		t.Fatalf("expected position 38")
	}
	assertDisplay(t, screen.Screen, []string{"35  ", "36  ", "37  ", "38  "})
	if len(screen.history.Bottom.items) != 2 {
		t.Fatalf("expected bottom history")
	}

	screen.PrevPage()
	if screen.history.Position != 36 {
		t.Fatalf("expected position 36")
	}
	assertDisplay(t, screen.Screen, []string{"33  ", "34  ", "35  ", "36  "})
}

func TestPyteHistoryPrevPageOdd(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(strconv.Itoa(idx))
		screen.LineFeed()
	}
	assertDisplay(t, screen.Screen, []string{"46   ", "47   ", "48   ", "49   ", "     "})
	screen.PrevPage()
	assertDisplay(t, screen.Screen, []string{"43   ", "44   ", "45   ", "46   ", "47   "})
}

func TestPyteHistoryPrevPageRatio(t *testing.T) {
	screen := NewHistoryScreenWithRatio(4, 4, 40, 0.75)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(strconv.Itoa(idx))
		screen.LineFeed()
	}
	assertDisplay(t, screen.Screen, []string{"37  ", "38  ", "39  ", "    "})
	screen.PrevPage()
	assertDisplay(t, screen.Screen, []string{"34  ", "35  ", "36  ", "37  "})
}

func TestPyteHistoryNextPage(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*5; idx++ {
		screen.Draw(strconv.Itoa(idx))
		screen.LineFeed()
	}
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "     "})
	screen.PrevPage()
	screen.NextPage()
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "     "})
}

func TestPyteHistoryEnsureWidth(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	stream := NewStream(screen, false)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed("0000\n")
	}
	screen.Resize(5, 3)
	screen.PrevPage()
	for _, line := range screen.history.Bottom.items {
		if len(line) > 3 {
			t.Fatalf("expected truncated history")
		}
	}
}

func TestPyteHistoryNotEnoughLines(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 6)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(string(rune('0' + idx)))
		screen.LineFeed()
	}
	screen.PrevPage()
	if len(screen.history.Bottom.items) != 1 {
		t.Fatalf("expected bottom history")
	}
	assertDisplay(t, screen.Screen, []string{"0    ", "1    ", "2    ", "3    ", "4    "})
	screen.NextPage()
	assertDisplay(t, screen.Screen, []string{"1    ", "2    ", "3    ", "4    ", "     "})
}

func TestPyteHistoryDraw(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	screen.SetMode([]int{ModeLNM}, false)
	stream := NewStream(screen, false)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed(strconv.Itoa(idx) + "\n")
	}
	screen.PrevPage()
	screen.PrevPage()
	screen.NextPage()
	screen.Draw("x")
	assertDisplay(t, screen.Screen, []string{"21   ", "22   ", "23   ", "24   ", "x    "})
}

func TestPyteHistoryCursorHidden(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	stream := NewStream(screen, false)
	for idx := 0; idx < screen.Lines*5; idx++ {
		stream.Feed(strconv.Itoa(idx) + "\n")
	}
	screen.PrevPage()
	if !screen.Cursor.Hidden {
		t.Fatalf("expected cursor hidden")
	}
	screen.NextPage()
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
	}
}

func TestPyteHistoryEraseInDisplay(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 6)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(string(rune('0' + idx)))
		screen.LineFeed()
	}
	screen.PrevPage()
	screen.EraseInDisplay(3, false)
	if len(screen.history.Top.items) != 0 || len(screen.history.Bottom.items) != 0 {
		t.Fatalf("expected history reset")
	}
}
