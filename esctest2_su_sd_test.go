package te

import "testing"

func rowString(screen *Screen, row int) string {
	cells := screen.Buffer[row]
	out := make([]rune, len(cells))
	for i, cell := range cells {
		if cell.Data == "" {
			out[i] = ' '
		} else {
			r := []rune(cell.Data)
			if len(r) > 0 {
				out[i] = r[0]
			} else {
				out[i] = ' '
			}
		}
	}
	return string(out)
}

func fillScreenLines(screen *Screen, lines []string) {
	for i, line := range lines {
		setCursor(screen, 1, i+1)
		screen.Draw(line)
	}
}

func TestEsctest2SUDefault(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.ScrollUp(1)
	if rowString(screen, 0) != "fghij" || rowString(screen, 4) != "     " {
		t.Fatalf("unexpected rows: %q %q", rowString(screen, 0), rowString(screen, 4))
	}
}

func TestEsctest2SDDefault(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.ScrollDown(1)
	if rowString(screen, 0) != "     " || rowString(screen, 1) != "abcde" {
		t.Fatalf("unexpected rows: %q %q", rowString(screen, 0), rowString(screen, 1))
	}
}

func TestEsctest2SULeftRightMargins(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	screen.ScrollUp(2)
	if rowString(screen, 0) != "almne" {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 3) != "p   t" {
		t.Fatalf("unexpected row: %q", rowString(screen, 3))
	}
}

func TestEsctest2SDLeftRightMargins(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	screen.ScrollDown(2)
	if rowString(screen, 0) != "a   e" {
		t.Fatalf("unexpected row: %q", rowString(screen, 0))
	}
	if rowString(screen, 2) != "kbcdo" {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}
