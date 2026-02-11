package te

type DiffScreen struct {
	*Screen
	dirty map[int]struct{}
}

func NewDiffScreen(cols, lines int) *DiffScreen {
	return &DiffScreen{
		Screen: NewScreen(cols, lines),
		dirty:  make(map[int]struct{}),
	}
}

func (s *DiffScreen) PutRune(ch rune) {
	row := s.cursor.Row
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
		s.markDirtyRow(row)
		return
	}
	s.cursor.Col++
	if s.cursor.Col >= s.cols {
		s.cursor.Col = s.cols - 1
	}
	s.markDirtyRow(row)
}

func (s *DiffScreen) LineFeed() {
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

func (s *DiffScreen) EraseInDisplay(mode int) {
	s.Screen.EraseInDisplay(mode)
	s.markDisplayDirty(mode)
}

func (s *DiffScreen) EraseInLine(mode int) {
	s.Screen.EraseInLine(mode)
	s.markDirtyRow(s.cursor.Row)
}

func (s *DiffScreen) EraseChars(n int) {
	s.Screen.EraseChars(n)
	s.markDirtyRow(s.cursor.Row)
}

func (s *DiffScreen) DeleteChars(n int) {
	s.Screen.DeleteChars(n)
	s.markDirtyRow(s.cursor.Row)
}

func (s *DiffScreen) InsertChars(n int) {
	s.Screen.InsertChars(n)
	s.markDirtyRow(s.cursor.Row)
}

func (s *DiffScreen) DeleteLines(n int) {
	s.Screen.DeleteLines(n)
	s.markDirtyRange(s.cursor.Row, s.scrollBottom)
}

func (s *DiffScreen) InsertLines(n int) {
	s.Screen.InsertLines(n)
	s.markDirtyRange(s.cursor.Row, s.scrollBottom)
}

func (s *DiffScreen) ScrollUp(n int) {
	s.Screen.ScrollUp(n)
	s.markDirtyRange(s.scrollTop, s.scrollBottom)
}

func (s *DiffScreen) ScrollDown(n int) {
	s.Screen.ScrollDown(n)
	s.markDirtyRange(s.scrollTop, s.scrollBottom)
}

func (s *DiffScreen) Reset() {
	s.Screen.Reset()
	s.ClearDirty()
}

func (s *DiffScreen) DirtyLines() []int {
	lines := make([]int, 0, len(s.dirty))
	for line := range s.dirty {
		lines = append(lines, line)
	}
	return lines
}

func (s *DiffScreen) ClearDirty() {
	s.dirty = make(map[int]struct{})
}

func (s *DiffScreen) markDisplayDirty(mode int) {
	switch mode {
	case 0:
		s.markDirtyRange(s.cursor.Row, s.lines-1)
	case 1:
		s.markDirtyRange(0, s.cursor.Row)
	case 2:
		s.markDirtyRange(0, s.lines-1)
	}
}

func (s *DiffScreen) markDirtyRow(row int) {
	if row < 0 || row >= s.lines {
		return
	}
	s.dirty[row] = struct{}{}
}

func (s *DiffScreen) markDirtyRange(start, end int) {
	if start < 0 {
		start = 0
	}
	if end >= s.lines {
		end = s.lines - 1
	}
	for row := start; row <= end; row++ {
		s.dirty[row] = struct{}{}
	}
}
