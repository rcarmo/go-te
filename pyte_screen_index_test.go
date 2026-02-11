package te

import "testing"

func TestPyteIndexWithMargins(t *testing.T) {
	screen := updateScreen(NewScreen(2, 5), []string{"bo", "sh", "th", "er", "oh"}, coloredLines(1, 2))
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 3

	screen.Index()
	if screen.Display()[0] != "bo" || screen.Display()[1] != "th" || screen.Display()[2] != "er" || screen.Display()[3] != "  " || screen.Display()[4] != "oh" {
		t.Fatalf("unexpected display after index")
	}

	screen.Index()
	if screen.Display()[1] != "er" || screen.Display()[2] != "  " {
		t.Fatalf("unexpected display after second index")
	}

	screen.Index()
	if screen.Display()[1] != "  " || screen.Display()[2] != "  " {
		t.Fatalf("unexpected display after third index")
	}
}

func TestPyteReverseIndexWithMargins(t *testing.T) {
	screen := updateScreen(NewScreen(2, 5), []string{"bo", "sh", "th", "er", "oh"}, coloredLines(2, 3))
	screen.SetMargins(2, 4)
	screen.Cursor.Row = 1

	screen.ReverseIndex()
	if screen.Display()[0] != "bo" || screen.Display()[1] != "  " || screen.Display()[2] != "sh" || screen.Display()[3] != "th" || screen.Display()[4] != "oh" {
		t.Fatalf("unexpected display after reverse index")
	}

	screen.ReverseIndex()
	if screen.Display()[2] != "  " || screen.Display()[3] != "sh" {
		t.Fatalf("unexpected display after second reverse index")
	}

	screen.ReverseIndex()
	if screen.Display()[2] != "  " || screen.Display()[3] != "  " {
		t.Fatalf("unexpected display after third reverse index")
	}
}
