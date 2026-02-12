package te

import "testing"

func TestEsctest2DLDefault(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	setCursor(screen, 1, 2)
	screen.DeleteLines(1)
	if rowString(screen, 1) != "klmno" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
	if rowString(screen, 4) != "     " {
		t.Fatalf("unexpected row: %q", rowString(screen, 4))
	}
}

func TestEsctest2DLLeftRightMargins(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	setCursor(screen, 3, 2)
	screen.DeleteLines(1)
	if rowString(screen, 1) != "flmnj" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
}

func TestEsctest2ILDefault(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	setCursor(screen, 1, 2)
	screen.InsertLines(1)
	if rowString(screen, 1) != "     " {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
	if rowString(screen, 2) != "fghij" {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}

func TestEsctest2ILLeftRightMargins(t *testing.T) {
	screen := NewScreen(5, 5)
	fillScreenLines(screen, []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"})
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(2, 4)
	setCursor(screen, 3, 2)
	screen.InsertLines(1)
	if rowString(screen, 1) != "f   j" {
		t.Fatalf("unexpected row: %q", rowString(screen, 1))
	}
	if rowString(screen, 2) != "kghio" {
		t.Fatalf("unexpected row: %q", rowString(screen, 2))
	}
}
