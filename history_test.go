package te

import (
	"strings"
	"testing"
)

func historyChars(lines [][]Cell, columns int) []string {
	out := make([]string, len(lines))
	for i := range lines {
		var b strings.Builder
		for x := 0; x < columns; x++ {
			b.WriteString(lines[i][x].Data)
		}
		out[i] = b.String()
	}
	return out
}

func historyLineString(line []Cell) string {
	var b strings.Builder
	for _, cell := range line {
		b.WriteString(cell.Data)
	}
	return b.String()
}

func TestHistoryIndex(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < screen.Lines; idx++ {
		screen.Draw(string(rune('0' + idx)))
		if idx != screen.Lines-1 {
			screen.LineFeed()
		}
	}
	line := screen.Buffer[0]
	screen.Index()
	if len(screen.history.Top.items) == 0 {
		t.Fatalf("expected top history updated")
	}
	if historyLineString(screen.history.Top.items[len(screen.history.Top.items)-1]) != historyLineString(line) {
		t.Fatalf("expected top history line match")
	}
}

func TestHistoryReverseIndex(t *testing.T) {
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
		t.Fatalf("expected bottom history updated")
	}
	if historyLineString(screen.history.Bottom.items[0]) != historyLineString(line) {
		t.Fatalf("expected bottom history line match")
	}
}

func TestHistoryPrevNextPage(t *testing.T) {
	screen := NewHistoryScreen(4, 4, 40)
	screen.SetMode([]int{ModeLNM}, false)
	for idx := 0; idx < screen.Lines*10; idx++ {
		screen.Draw(string(rune('0' + idx%10)))
		screen.LineFeed()
	}
	if screen.history.Position != 40 {
		t.Fatalf("expected position 40")
	}
	screen.PrevPage()
	if screen.history.Position >= 40 {
		t.Fatalf("expected position decreased")
	}
	screen.NextPage()
	if screen.history.Position != 40 {
		t.Fatalf("expected position restored")
	}
}

func TestHistoryCursorHidden(t *testing.T) {
	screen := NewHistoryScreen(5, 5, 50)
	for idx := 0; idx < screen.Lines*5; idx++ {
		screen.Draw(string(rune('0' + idx%10)))
		screen.LineFeed()
	}
	if screen.Cursor.Hidden {
		t.Fatalf("expected cursor visible")
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

func TestHistoryEraseInDisplay(t *testing.T) {
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
