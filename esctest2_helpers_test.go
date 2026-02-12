package te

import "testing"

func setCursor(screen *Screen, col, row int) {
	screen.CursorPosition(row, col)
}

func assertCursor(t *testing.T, screen *Screen, col, row int) {
	t.Helper()
	if screen.Cursor.Col != col-1 || screen.Cursor.Row != row-1 {
		t.Fatalf("expected cursor (%d,%d), got (%d,%d)", col, row, screen.Cursor.Col+1, screen.Cursor.Row+1)
	}
}

func assertCell(t *testing.T, screen *Screen, col, row int, expected string) {
	t.Helper()
	cell := screen.Buffer[row-1][col-1]
	if cell.Data != expected {
		t.Fatalf("expected cell (%d,%d) to be %q, got %q", col, row, expected, cell.Data)
	}
}
