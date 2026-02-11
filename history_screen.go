package te

type HistoryScreen struct {
	*Screen
	history [][]Cell
	max     int
}

func NewHistoryScreen(cols, lines, history int) *HistoryScreen {
	return &HistoryScreen{
		Screen:  NewScreen(cols, lines),
		history: make([][]Cell, 0, history),
		max:     history,
	}
}

func (s *HistoryScreen) PutRune(ch rune) {
	if s.wrapPending {
		s.LineFeed()
		s.cursor.Col = 0
		s.wrapPending = false
	}
	if s.cursor.Row < 0 || s.cursor.Row >= s.lines || s.cursor.Col < 0 || s.cursor.Col >= s.cols {
		return
	}
	if s.insertMode {
		for col := s.cols - 1; col > s.cursor.Col; col-- {
			s.cells[s.cursor.Row][col] = s.cells[s.cursor.Row][col-1]
		}
		s.cells[s.cursor.Row][s.cursor.Col] = defaultCell()
	}
	s.cells[s.cursor.Row][s.cursor.Col] = Cell{Ch: ch, Attr: s.attr}
	if s.autowrap && s.cursor.Col == s.cols-1 {
		s.wrapPending = true
		return
	}
	s.cursor.Col++
	if s.cursor.Col >= s.cols {
		s.cursor.Col = s.cols - 1
	}
}

func (s *HistoryScreen) LineFeed() {
	if s.newlineMode {
		s.cursor.Col = 0
	}
	if s.cursor.Row == s.scrollBottom {
		s.ScrollUp(1)
	} else if s.cursor.Row < s.lines-1 {
		s.cursor.Row++
	}
	s.wrapPending = false
}

func (s *HistoryScreen) ScrollUp(n int) {
	if n <= 0 {
		return
	}
	if n > s.scrollBottom-s.scrollTop+1 {
		n = s.scrollBottom - s.scrollTop + 1
	}
	for i := 0; i < n; i++ {
		s.captureLine(s.scrollTop)
		s.Screen.ScrollUp(1)
	}
}

func (s *HistoryScreen) Reset() {
	s.Screen.Reset()
	s.history = s.history[:0]
}

func (s *HistoryScreen) History() [][]Cell {
	return append([][]Cell(nil), s.history...)
}

func (s *HistoryScreen) Scrollback() int {
	return len(s.history)
}

func (s *HistoryScreen) captureLine(row int) {
	if row < 0 || row >= s.lines {
		return
	}
	lineCopy := append([]Cell(nil), s.cells[row]...)
	s.history = append(s.history, lineCopy)
	if s.max > 0 && len(s.history) > s.max {
		s.history = s.history[len(s.history)-s.max:]
	}
}
