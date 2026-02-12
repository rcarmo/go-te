package te

import "testing"

type lineFeedFunc func(*Screen)

func runLineFeedTests(t *testing.T, name string, feed lineFeedFunc) {
	t.Helper()
	t.Run(name+"_Basic", func(t *testing.T) {
		screen := NewScreen(80, 24)
		setCursor(screen, 5, 3)
		feed(screen)
		assertCursor(t, screen, 5, 4)
	})

	t.Run(name+"_Scrolls", func(t *testing.T) {
		screen := NewScreen(80, 24)
		height := screen.Lines
		setCursor(screen, 2, height-1)
		screen.Draw("a")
		setCursor(screen, 2, height)
		screen.Draw("b")
		setCursor(screen, 2, height-1)
		feed(screen)
		assertCursor(t, screen, 2, height)
		assertCell(t, screen, 2, height-1, "a")
		assertCell(t, screen, 2, height, "b")
		feed(screen)
		assertCursor(t, screen, 2, height)
		assertCell(t, screen, 2, height-1, "b")
		assertCell(t, screen, 2, height, " ")
	})

	t.Run(name+"_ScrollsInTopBottomRegionStartingAbove", func(t *testing.T) {
		screen := NewScreen(80, 24)
		screen.SetMargins(4, 5)
		setCursor(screen, 2, 5)
		screen.Draw("x")
		setCursor(screen, 2, 3)
		feed(screen)
		feed(screen)
		feed(screen)
		assertCursor(t, screen, 2, 5)
		assertCell(t, screen, 2, 4, "x")
		assertCell(t, screen, 2, 5, " ")
	})

	t.Run(name+"_ScrollsInTopBottomRegionStartingWithin", func(t *testing.T) {
		screen := NewScreen(80, 24)
		screen.SetMargins(4, 5)
		setCursor(screen, 2, 5)
		screen.Draw("x")
		setCursor(screen, 2, 4)
		feed(screen)
		feed(screen)
		assertCursor(t, screen, 2, 5)
		assertCell(t, screen, 2, 4, "x")
		assertCell(t, screen, 2, 5, " ")
	})

	t.Run(name+"_MovesDoesNotScrollOutsideLeftRight", func(t *testing.T) {
		screen := NewScreen(80, 24)
		screen.SetMargins(2, 5)
		screen.SetMode([]int{69}, true)
		screen.SetLeftRightMargins(2, 5)
		setCursor(screen, 3, 5)
		screen.Draw("x")
		setCursor(screen, 6, 5)
		feed(screen)
		assertCursor(t, screen, 6, 5)
		assertCell(t, screen, 3, 5, "x")
		setCursor(screen, 6, 4)
		feed(screen)
		assertCursor(t, screen, 6, 5)
		setCursor(screen, 6, screen.Lines)
		feed(screen)
		assertCursor(t, screen, 6, screen.Lines)
		assertCell(t, screen, 3, 5, "x")
		setCursor(screen, 1, 5)
		feed(screen)
		assertCursor(t, screen, 1, 5)
		assertCell(t, screen, 3, 5, "x")
		setCursor(screen, 1, screen.Lines)
		feed(screen)
		assertCursor(t, screen, 1, screen.Lines)
		assertCell(t, screen, 3, 5, "x")
	})

	t.Run(name+"_StopsAtBottomLineWhenBegunBelowScrollRegion", func(t *testing.T) {
		screen := NewScreen(80, 24)
		screen.SetMargins(4, 5)
		setCursor(screen, 1, 6)
		screen.Draw("x")
		for i := 0; i < screen.Lines; i++ {
			feed(screen)
		}
		assertCursor(t, screen, 2, screen.Lines)
		assertCell(t, screen, 1, 6, "x")
	})
}

func TestEsctest2LF(t *testing.T) {
	runLineFeedTests(t, "LF", func(screen *Screen) {
		screen.LineFeed()
	})
}

func TestEsctest2VT(t *testing.T) {
	runLineFeedTests(t, "VT", func(screen *Screen) {
		screen.LineFeed()
	})
}

func TestEsctest2FF(t *testing.T) {
	runLineFeedTests(t, "FF", func(screen *Screen) {
		screen.LineFeed()
	})
}
