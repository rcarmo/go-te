package te

import "strings"

type Screen struct {
	cols         int
	lines        int
	cells        [][]Cell
	cursor       Cursor
	savedCursor  Cursor
	attr         Attr
	insertMode   bool
	originMode   bool
	autowrap     bool
	newlineMode  bool
	wrapPending  bool
	scrollTop    int
	scrollBottom int
	tabStops     map[int]struct{}
	usingAlt     bool
	primary      bufferState
}

type bufferState struct {
	cells        [][]Cell
	cursor       Cursor
	savedCursor  Cursor
	attr         Attr
	insertMode   bool
	originMode   bool
	autowrap     bool
	newlineMode  bool
	wrapPending  bool
	scrollTop    int
	scrollBottom int
	tabStops     map[int]struct{}
}

func NewScreen(cols, lines int) *Screen {
	s := &Screen{}
	s.Resize(cols, lines)
	return s
}

func (s *Screen) Resize(cols, lines int) {
	if cols <= 0 {
		cols = 1
	}
	if lines <= 0 {
		lines = 1
	}
	s.cols = cols
	s.lines = lines
	s.cells = make([][]Cell, lines)
	for row := 0; row < lines; row++ {
		s.cells[row] = make([]Cell, cols)
		for col := 0; col < cols; col++ {
			s.cells[row][col] = defaultCell()
		}
	}
	s.resetState()
}

func (s *Screen) Reset() {
	for row := 0; row < s.lines; row++ {
		for col := 0; col < s.cols; col++ {
			s.cells[row][col] = defaultCell()
		}
	}
	s.resetState()
}

func (s *Screen) resetState() {
	s.cursor = Cursor{}
	s.savedCursor = Cursor{}
	s.attr = Attr{}
	s.insertMode = false
	s.originMode = false
	s.autowrap = true
	s.newlineMode = false
	s.wrapPending = false
	s.scrollTop = 0
	s.scrollBottom = s.lines - 1
	s.tabStops = make(map[int]struct{})
	for col := 0; col < s.cols; col += 8 {
		s.tabStops[col] = struct{}{}
	}
}

func (s *Screen) Display() []string {
	lines := make([]string, s.lines)
	for row := 0; row < s.lines; row++ {
		var b strings.Builder
		b.Grow(s.cols)
		for col := 0; col < s.cols; col++ {
			b.WriteRune(s.cells[row][col].Ch)
		}
		lines[row] = b.String()
	}
	return lines
}

func (s *Screen) Lines() [][]Cell {
	lines := make([][]Cell, s.lines)
	for row := range s.cells {
		lines[row] = append([]Cell(nil), s.cells[row]...)
	}
	return lines
}

func (s *Screen) Cursor() Cursor {
	return s.cursor
}

func (s *Screen) Size() (int, int) {
	return s.cols, s.lines
}

func (s *Screen) EnableAlternateBuffer(clear bool) {
	if s.usingAlt {
		if clear {
			s.clearAll()
		}
		return
	}
	s.primary = s.snapshot()
	s.usingAlt = true
	s.cells = makeBlankCells(s.lines, s.cols)
	s.resetState()
	if clear {
		s.clearAll()
	}
}

func (s *Screen) DisableAlternateBuffer() {
	if !s.usingAlt {
		return
	}
	s.restore(s.primary)
	s.usingAlt = false
}

func (s *Screen) SetTabStop(col int) {
	if col < 0 || col >= s.cols {
		return
	}
	s.tabStops[col] = struct{}{}
}

func (s *Screen) ClearTabStop(col int) {
	delete(s.tabStops, col)
}

func (s *Screen) ClearAllTabStops() {
	s.tabStops = make(map[int]struct{})
}

func (s *Screen) SetInsertMode(enabled bool) {
	s.insertMode = enabled
}

func (s *Screen) SetAutowrap(enabled bool) {
	s.autowrap = enabled
}

func (s *Screen) SetOriginMode(enabled bool) {
	s.originMode = enabled
	if enabled {
		s.cursor.Row = s.scrollTop
		s.cursor.Col = 0
	} else {
		s.cursor.Row = 0
		s.cursor.Col = 0
	}
	s.wrapPending = false
}

func (s *Screen) SetNewlineMode(enabled bool) {
	s.newlineMode = enabled
}

func (s *Screen) SaveCursor() {
	s.savedCursor = s.cursor
}

func (s *Screen) RestoreCursor() {
	s.cursor = s.savedCursor
	s.wrapPending = false
}

func (s *Screen) PutRune(ch rune) {
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

func (s *Screen) Backspace() {
	if s.cursor.Col > 0 {
		s.cursor.Col--
	}
	s.wrapPending = false
}

func (s *Screen) Tab() {
	next := s.cols - 1
	for col := s.cursor.Col + 1; col < s.cols; col++ {
		if _, ok := s.tabStops[col]; ok {
			next = col
			break
		}
	}
	s.cursor.Col = next
	s.wrapPending = false
}

func (s *Screen) LineFeed() {
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

func (s *Screen) CarriageReturn() {
	s.cursor.Col = 0
	s.wrapPending = false
}

func (s *Screen) MoveCursor(rowDelta, colDelta int) {
	row := s.cursor.Row + rowDelta
	col := s.cursor.Col + colDelta
	if s.originMode {
		row = s.clampRowInRegion(row)
	} else {
		row = s.clampRow(row)
	}
	col = s.clampCol(col)
	s.cursor.Row = row
	s.cursor.Col = col
	s.wrapPending = false
}

func (s *Screen) SetCursor(row, col int) {
	if s.originMode {
		row += s.scrollTop
		row = s.clampRowInRegion(row)
	} else {
		row = s.clampRow(row)
	}
	col = s.clampCol(col)
	s.cursor.Row = row
	s.cursor.Col = col
	s.wrapPending = false
}

func (s *Screen) SetScrollRegion(top, bottom int) {
	if top < 0 || bottom >= s.lines || top >= bottom {
		s.scrollTop = 0
		s.scrollBottom = s.lines - 1
	} else {
		s.scrollTop = top
		s.scrollBottom = bottom
	}
	s.cursor.Row = s.scrollTop
	s.cursor.Col = 0
	s.wrapPending = false
}

func (s *Screen) ScrollUp(n int) {
	if n <= 0 {
		return
	}
	if n > s.scrollBottom-s.scrollTop+1 {
		n = s.scrollBottom - s.scrollTop + 1
	}
	for i := 0; i < n; i++ {
		for row := s.scrollTop; row < s.scrollBottom; row++ {
			s.cells[row] = s.cells[row+1]
		}
		s.cells[s.scrollBottom] = blankLine(s.cols)
	}
}

func (s *Screen) ScrollDown(n int) {
	if n <= 0 {
		return
	}
	if n > s.scrollBottom-s.scrollTop+1 {
		n = s.scrollBottom - s.scrollTop + 1
	}
	for i := 0; i < n; i++ {
		for row := s.scrollBottom; row > s.scrollTop; row-- {
			s.cells[row] = s.cells[row-1]
		}
		s.cells[s.scrollTop] = blankLine(s.cols)
	}
}

func (s *Screen) EraseInDisplay(mode int) {
	switch mode {
	case 0:
		s.EraseInLine(0)
		for row := s.cursor.Row + 1; row < s.lines; row++ {
			s.clearLine(row)
		}
	case 1:
		s.EraseInLine(1)
		for row := 0; row < s.cursor.Row; row++ {
			s.clearLine(row)
		}
	case 2:
		for row := 0; row < s.lines; row++ {
			s.clearLine(row)
		}
	}
}

func (s *Screen) EraseInLine(mode int) {
	row := s.cursor.Row
	if row < 0 || row >= s.lines {
		return
	}
	switch mode {
	case 0:
		for col := s.cursor.Col; col < s.cols; col++ {
			s.cells[row][col] = defaultCell()
		}
	case 1:
		for col := 0; col <= s.cursor.Col && col < s.cols; col++ {
			s.cells[row][col] = defaultCell()
		}
	case 2:
		s.clearLine(row)
	}
}

func (s *Screen) EraseChars(n int) {
	if n <= 0 {
		n = 1
	}
	row := s.cursor.Row
	if row < 0 || row >= s.lines {
		return
	}
	for col := s.cursor.Col; col < s.cols && n > 0; col++ {
		s.cells[row][col] = defaultCell()
		n--
	}
}

func (s *Screen) DeleteChars(n int) {
	if n <= 0 {
		n = 1
	}
	row := s.cursor.Row
	if row < 0 || row >= s.lines {
		return
	}
	if n > s.cols-s.cursor.Col {
		n = s.cols - s.cursor.Col
	}
	for col := s.cursor.Col; col < s.cols-n; col++ {
		s.cells[row][col] = s.cells[row][col+n]
	}
	for col := s.cols - n; col < s.cols; col++ {
		s.cells[row][col] = defaultCell()
	}
}

func (s *Screen) InsertChars(n int) {
	if n <= 0 {
		n = 1
	}
	row := s.cursor.Row
	if row < 0 || row >= s.lines {
		return
	}
	if n > s.cols-s.cursor.Col {
		n = s.cols - s.cursor.Col
	}
	for col := s.cols - 1; col >= s.cursor.Col+n; col-- {
		s.cells[row][col] = s.cells[row][col-n]
	}
	for col := s.cursor.Col; col < s.cursor.Col+n; col++ {
		s.cells[row][col] = defaultCell()
	}
}

func (s *Screen) DeleteLines(n int) {
	if n <= 0 {
		n = 1
	}
	if s.cursor.Row < s.scrollTop || s.cursor.Row > s.scrollBottom {
		return
	}
	max := s.scrollBottom - s.cursor.Row + 1
	if n > max {
		n = max
	}
	for i := 0; i < n; i++ {
		for row := s.cursor.Row; row < s.scrollBottom; row++ {
			s.cells[row] = s.cells[row+1]
		}
		s.cells[s.scrollBottom] = blankLine(s.cols)
	}
}

func (s *Screen) InsertLines(n int) {
	if n <= 0 {
		n = 1
	}
	if s.cursor.Row < s.scrollTop || s.cursor.Row > s.scrollBottom {
		return
	}
	max := s.scrollBottom - s.cursor.Row + 1
	if n > max {
		n = max
	}
	for i := 0; i < n; i++ {
		for row := s.scrollBottom; row > s.cursor.Row; row-- {
			s.cells[row] = s.cells[row-1]
		}
		s.cells[s.cursor.Row] = blankLine(s.cols)
	}
}

func (s *Screen) ApplySGR(params []int) {
	if len(params) == 0 {
		params = []int{0}
	}
	for i := 0; i < len(params); i++ {
		p := params[i]
		switch {
		case p == 0:
			s.attr = Attr{}
		case p == 1:
			s.attr.Bold = true
		case p == 4:
			s.attr.Underline = true
		case p == 5:
			s.attr.Blink = true
		case p == 7:
			s.attr.Reverse = true
		case p == 8:
			s.attr.Conceal = true
		case p == 22:
			s.attr.Bold = false
		case p == 24:
			s.attr.Underline = false
		case p == 25:
			s.attr.Blink = false
		case p == 27:
			s.attr.Reverse = false
		case p == 28:
			s.attr.Conceal = false
		case p == 39:
			s.attr.Fg = Color{}
		case p == 49:
			s.attr.Bg = Color{}
		case p >= 30 && p <= 37:
			s.attr.Fg = Color{Mode: ColorANSI16, Index: uint8(p - 30)}
		case p >= 90 && p <= 97:
			s.attr.Fg = Color{Mode: ColorANSI16, Index: uint8(p - 90 + 8)}
		case p >= 40 && p <= 47:
			s.attr.Bg = Color{Mode: ColorANSI16, Index: uint8(p - 40)}
		case p >= 100 && p <= 107:
			s.attr.Bg = Color{Mode: ColorANSI16, Index: uint8(p - 100 + 8)}
		case p == 38 || p == 48:
			if i+2 < len(params) && params[i+1] == 5 {
				index := params[i+2]
				if index >= 0 && index <= 255 {
					color := Color{Mode: ColorANSI256, Index: uint8(index)}
					if p == 38 {
						s.attr.Fg = color
					} else {
						s.attr.Bg = color
					}
				}
				i += 2
			}
		}
	}
}

func (s *Screen) snapshot() bufferState {
	return bufferState{
		cells:        cloneCells(s.cells),
		cursor:       s.cursor,
		savedCursor:  s.savedCursor,
		attr:         s.attr,
		insertMode:   s.insertMode,
		originMode:   s.originMode,
		autowrap:     s.autowrap,
		newlineMode:  s.newlineMode,
		wrapPending:  s.wrapPending,
		scrollTop:    s.scrollTop,
		scrollBottom: s.scrollBottom,
		tabStops:     cloneTabStops(s.tabStops),
	}
}

func (s *Screen) restore(state bufferState) {
	s.cells = cloneCells(state.cells)
	s.cursor = state.cursor
	s.savedCursor = state.savedCursor
	s.attr = state.attr
	s.insertMode = state.insertMode
	s.originMode = state.originMode
	s.autowrap = state.autowrap
	s.newlineMode = state.newlineMode
	s.wrapPending = state.wrapPending
	s.scrollTop = state.scrollTop
	s.scrollBottom = state.scrollBottom
	s.tabStops = cloneTabStops(state.tabStops)
}

func (s *Screen) clearAll() {
	for row := 0; row < s.lines; row++ {
		for col := 0; col < s.cols; col++ {
			s.cells[row][col] = defaultCell()
		}
	}
}

func (s *Screen) clampRow(row int) int {
	if row < 0 {
		return 0
	}
	if row >= s.lines {
		return s.lines - 1
	}
	return row
}

func (s *Screen) clampRowInRegion(row int) int {
	if row < s.scrollTop {
		return s.scrollTop
	}
	if row > s.scrollBottom {
		return s.scrollBottom
	}
	return row
}

func (s *Screen) clampCol(col int) int {
	if col < 0 {
		return 0
	}
	if col >= s.cols {
		return s.cols - 1
	}
	return col
}

func (s *Screen) clearLine(row int) {
	if row < 0 || row >= s.lines {
		return
	}
	for col := 0; col < s.cols; col++ {
		s.cells[row][col] = defaultCell()
	}
}

func blankLine(cols int) []Cell {
	line := make([]Cell, cols)
	for col := 0; col < cols; col++ {
		line[col] = defaultCell()
	}
	return line
}

func makeBlankCells(lines, cols int) [][]Cell {
	cells := make([][]Cell, lines)
	for row := 0; row < lines; row++ {
		cells[row] = blankLine(cols)
	}
	return cells
}

func cloneCells(source [][]Cell) [][]Cell {
	clone := make([][]Cell, len(source))
	for row := range source {
		clone[row] = append([]Cell(nil), source[row]...)
	}
	return clone
}

func cloneTabStops(source map[int]struct{}) map[int]struct{} {
	clone := make(map[int]struct{}, len(source))
	for key := range source {
		clone[key] = struct{}{}
	}
	return clone
}

func defaultCell() Cell {
	return Cell{Ch: ' ', Attr: Attr{}}
}
