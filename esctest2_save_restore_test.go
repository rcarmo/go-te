package te

import "testing"

type saveRestoreCase struct {
	name       string
	save       func(*Stream) error
	restore    func(*Stream) error
	lrmWorks   bool
}

func runSaveRestoreCursorCases(t *testing.T, tc saveRestoreCase) {
	t.Helper()
	base := func(t *testing.T, lrmWorks bool) {
		screen := NewScreen(80, 24)
		stream := NewStream(screen, false)
		setCursor(screen, 5, 6)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		setCursor(screen, 1, 1)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		assertCursor(t, screen, 5, 6)

		// Restore with no save.
		screen = NewScreen(80, 24)
		stream = NewStream(screen, false)
		setCursor(screen, 5, 6)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		assertCursor(t, screen, 1, 1)

		// Reset origin mode on restore.
		screen = NewScreen(80, 24)
		stream = NewStream(screen, false)
		setCursor(screen, 5, 6)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		screen.SetMargins(5, 7)
		screen.SetMode([]int{69}, true)
		screen.SetLeftRightMargins(5, 7)
		screen.SetMode([]int{6}, true)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		screen.CursorPosition(1, 1)
		assertCursor(t, screen, 1, 1)

		// LRM behavior.
		screen = NewScreen(80, 24)
		stream = NewStream(screen, false)
		setCursor(screen, 2, 3)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		screen.SetMode([]int{69}, true)
		screen.SetLeftRightMargins(1, 10)
		setCursor(screen, 5, 6)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		setCursor(screen, 4, 5)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if lrmWorks {
			assertCursor(t, screen, 5, 6)
		} else {
			assertCursor(t, screen, 2, 3)
		}

		// Reverse wrap unaffected.
		screen = NewScreen(80, 24)
		stream = NewStream(screen, false)
		screen.SetMode([]int{7, 1045}, true)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		screen.ResetMode([]int{1045}, true)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		setCursor(screen, 1, 2)
		screen.Backspace()
		assertCursor(t, screen, 1, 2)

		// Insert mode unaffected.
		screen = NewScreen(80, 24)
		stream = NewStream(screen, false)
		screen.SetMode([]int{4}, false)
		if err := tc.save(stream); err != nil {
			t.Fatalf("save: %v", err)
		}
		screen.ResetMode([]int{4}, false)
		if err := tc.restore(stream); err != nil {
			t.Fatalf("restore: %v", err)
		}
		setCursor(screen, 1, 1)
		screen.Draw("a")
		setCursor(screen, 1, 1)
		screen.Draw("b")
		assertCell(t, screen, 1, 1, "b")
		assertCell(t, screen, 2, 1, " ")
	}

	t.Run(tc.name, func(t *testing.T) {
		base(t, tc.lrmWorks)
	})
}

func TestEsctest2DECRC(t *testing.T) {
	caseInfo := saveRestoreCase{
		name:     "DECRC",
		lrmWorks: true,
		save: func(stream *Stream) error {
			return stream.Feed(ControlESC + "7")
		},
		restore: func(stream *Stream) error {
			return stream.Feed(ControlESC + "8")
		},
	}
	runSaveRestoreCursorCases(t, caseInfo)
}

func TestEsctest2DECSETTiteInhibit(t *testing.T) {
	caseInfo := saveRestoreCase{
		name:     "DECSET_TITE_INHIBIT",
		lrmWorks: true,
		save: func(stream *Stream) error {
			return stream.Feed(ControlCSI + "?1048h")
		},
		restore: func(stream *Stream) error {
			return stream.Feed(ControlCSI + "?1048l")
		},
	}
	runSaveRestoreCursorCases(t, caseInfo)
}
