package te

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
)

type Margins struct {
	Top    int
	Bottom int
}

type Savepoint struct {
	Cursor   Cursor
	G0       []rune
	G1       []rune
	Charset  int
	Origin   bool
	Wrap     bool
	SavedCol *int
}

type Screen struct {
	Columns           int
	Lines             int
	Buffer            [][]Cell
	Dirty             map[int]struct{}
	Cursor            Cursor
	Mode              map[int]struct{}
	Margins           *Margins
	Charset           int
	G0                []rune
	G1                []rune
	TabStops          map[int]struct{}
	Title             string
	IconName          string
	Savepoints        []Savepoint
	SavedColumns      *int
	WriteProcessInput func(string)
}

func NewScreen(cols, lines int) *Screen {
	s := &Screen{}
	s.Resize(lines, cols)
	s.Reset()
	return s
}

func (s *Screen) Reset() {
	s.Dirty = make(map[int]struct{})
	s.markDirtyRange(0, s.Lines-1)
	s.Buffer = makeBlankCells(s.Lines, s.Columns)
	s.Margins = nil
	s.Mode = map[int]struct{}{
		ModeDECAWM:  {},
		ModeDECTCEM: {},
	}
	s.Title = ""
	s.IconName = ""
	s.Charset = 0
	s.G0 = charsetLat1
	s.G1 = charsetVT100
	s.TabStops = make(map[int]struct{})
	for col := 8; col < s.Columns; col += 8 {
		s.TabStops[col] = struct{}{}
	}
	s.Cursor = Cursor{Row: 0, Col: 0, Attr: s.defaultAttr(), Hidden: false}
	s.Savepoints = nil
	s.SavedColumns = nil
	if s.WriteProcessInput == nil {
		s.WriteProcessInput = func(string) {}
	}
}

func (s *Screen) Resize(lines, columns int) {
	if lines <= 0 {
		lines = 1
	}
	if columns <= 0 {
		columns = 1
	}
	if s.Lines == lines && s.Columns == columns && s.Buffer != nil {
		return
	}
	if s.Buffer == nil {
		s.Lines = lines
		s.Columns = columns
		s.Buffer = makeBlankCells(lines, columns)
		return
	}

	s.markDirtyRange(0, lines-1)

	if lines < s.Lines {
		s.SaveCursor()
		s.CursorPosition(0, 0)
		s.DeleteLines(s.Lines - lines)
		s.RestoreCursor()
	}

	if columns < s.Columns {
		for row := range s.Buffer {
			if len(s.Buffer[row]) > columns {
				s.Buffer[row] = s.Buffer[row][:columns]
			}
		}
	}

	if columns > s.Columns {
		for row := range s.Buffer {
			if len(s.Buffer[row]) < columns {
				missing := make([]Cell, columns-len(s.Buffer[row]))
				for i := range missing {
					missing[i] = s.defaultCell()
				}
				s.Buffer[row] = append(s.Buffer[row], missing...)
			}
		}
	}

	if lines > s.Lines {
		for i := 0; i < lines-s.Lines; i++ {
			s.Buffer = append(s.Buffer, blankLine(columns, s.defaultCell()))
		}
	}

	s.Lines = lines
	s.Columns = columns
	s.SetMargins(0, 0)
}

func (s *Screen) Display() []string {
	lines := make([]string, s.Lines)
	for row := 0; row < s.Lines; row++ {
		var b strings.Builder
		col := 0
		for col < s.Columns {
			cell := s.Buffer[row][col]
			if cell.Data == "" {
				col++
				continue
			}
			width := runewidth.StringWidth(cell.Data)
			if width == 0 {
				b.WriteString(cell.Data)
				col++
				continue
			}
			b.WriteString(cell.Data)
			if width > 1 {
				col += width
			} else {
				col++
			}
		}
		lines[row] = b.String()
	}
	return lines
}

func (s *Screen) LinesCells() [][]Cell {
	lines := make([][]Cell, s.Lines)
	for row := range s.Buffer {
		lines[row] = append([]Cell(nil), s.Buffer[row]...)
	}
	return lines
}

func (s *Screen) Draw(data string) {
	data = s.translate(data)
	graphemes := uniseg.NewGraphemes(data)
	for graphemes.Next() {
		cluster := graphemes.Str()
		width := runewidth.StringWidth(cluster)

		if s.Cursor.Col == s.Columns {
			if s.isModeSet(ModeDECAWM) {
				s.Dirty[s.Cursor.Row] = struct{}{}
				s.CarriageReturn()
				s.LineFeed()
			} else if width > 0 {
				s.Cursor.Col -= width
				if s.Cursor.Col < 0 {
					s.Cursor.Col = 0
				}
			}
		}

		if s.isModeSet(ModeIRM) && width > 0 {
			s.InsertCharacters(width)
		}

		line := s.Buffer[s.Cursor.Row]
		switch {
		case width == 1:
			line[s.Cursor.Col] = Cell{Data: cluster, Attr: s.Cursor.Attr}
		case width == 2:
			line[s.Cursor.Col] = Cell{Data: cluster, Attr: s.Cursor.Attr}
			if s.Cursor.Col+1 < s.Columns {
				line[s.Cursor.Col+1] = Cell{Data: "", Attr: s.Cursor.Attr}
			}
		case width == 0 && isCombiningCluster(cluster):
			if s.Cursor.Col > 0 {
				prev := line[s.Cursor.Col-1]
				line[s.Cursor.Col-1] = Cell{Data: norm.NFC.String(prev.Data + cluster), Attr: prev.Attr}
			} else if s.Cursor.Row > 0 {
				prevLine := s.Buffer[s.Cursor.Row-1]
				prev := prevLine[s.Columns-1]
				prevLine[s.Columns-1] = Cell{Data: norm.NFC.String(prev.Data + cluster), Attr: prev.Attr}
			}
		default:
			continue
		}

		if width > 0 {
			s.Cursor.Col += width
			if s.Cursor.Col > s.Columns {
				s.Cursor.Col = s.Columns
			}
		}
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
}

func (s *Screen) SetTitle(param string) {
	s.Title = param
}

func (s *Screen) SetIconName(param string) {
	s.IconName = param
}

func (s *Screen) CarriageReturn() {
	s.Cursor.Col = 0
}

func (s *Screen) Index() {
	top, bottom := s.scrollRegion()
	if s.Cursor.Row == bottom {
		s.markDirtyRange(0, s.Lines-1)
		for row := top; row < bottom; row++ {
			s.Buffer[row] = s.Buffer[row+1]
		}
		s.Buffer[bottom] = blankLine(s.Columns, s.defaultCell())
	} else {
		s.CursorDown(1)
	}
}

func (s *Screen) ReverseIndex() {
	top, bottom := s.scrollRegion()
	if s.Cursor.Row == top {
		s.markDirtyRange(0, s.Lines-1)
		for row := bottom; row > top; row-- {
			s.Buffer[row] = s.Buffer[row-1]
		}
		s.Buffer[top] = blankLine(s.Columns, s.defaultCell())
	} else {
		s.CursorUp(1)
	}
}

func (s *Screen) LineFeed() {
	s.Index()
	if s.isModeSet(ModeLNM) {
		s.CarriageReturn()
	}
}

func (s *Screen) Tab() {
	column := s.Columns - 1
	for _, stop := range sortedStops(s.TabStops) {
		if s.Cursor.Col < stop {
			column = stop
			break
		}
	}
	s.Cursor.Col = column
}

func (s *Screen) Backspace() {
	s.CursorBack(1)
}

func (s *Screen) SaveCursor() {
	modeOrigin := s.isModeSet(ModeDECOM)
	modeWrap := s.isModeSet(ModeDECAWM)
	s.Savepoints = append(s.Savepoints, Savepoint{
		Cursor:   s.Cursor,
		G0:       s.G0,
		G1:       s.G1,
		Charset:  s.Charset,
		Origin:   modeOrigin,
		Wrap:     modeWrap,
		SavedCol: s.SavedColumns,
	})
}

func (s *Screen) RestoreCursor() {
	if len(s.Savepoints) == 0 {
		s.ResetMode([]int{ModeDECOM}, false)
		s.CursorPosition(0, 0)
		return
	}

	last := s.Savepoints[len(s.Savepoints)-1]
	s.Savepoints = s.Savepoints[:len(s.Savepoints)-1]
	s.G0 = last.G0
	s.G1 = last.G1
	s.Charset = last.Charset
	if last.Origin {
		s.SetMode([]int{ModeDECOM}, false)
	} else {
		s.ResetMode([]int{ModeDECOM}, false)
	}
	if last.Wrap {
		s.SetMode([]int{ModeDECAWM}, false)
	} else {
		s.ResetMode([]int{ModeDECAWM}, false)
	}
	s.Cursor = last.Cursor
	s.ensureHB()
	s.ensureVB(true)
}

func (s *Screen) InsertLines(count int) {
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	s.markDirtyRange(s.Cursor.Row, s.Lines-1)
	for row := bottom; row >= s.Cursor.Row; row-- {
		if row+count <= bottom {
			s.Buffer[row+count] = s.Buffer[row]
		}
		s.Buffer[row] = blankLine(s.Columns, s.defaultCell())
	}
	s.CarriageReturn()
}

func (s *Screen) DeleteLines(count int) {
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	s.markDirtyRange(s.Cursor.Row, s.Lines-1)
	for row := s.Cursor.Row; row <= bottom; row++ {
		if row+count <= bottom {
			s.Buffer[row] = s.Buffer[row+count]
		} else {
			s.Buffer[row] = blankLine(s.Columns, s.defaultCell())
		}
	}
	s.CarriageReturn()
}

func (s *Screen) InsertCharacters(count int) {
	if count <= 0 {
		count = 1
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	for col := s.Columns - 1; col >= s.Cursor.Col; col-- {
		if col+count < s.Columns {
			line[col+count] = line[col]
		}
	}
	for col := s.Cursor.Col; col < s.Cursor.Col+count && col < s.Columns; col++ {
		line[col] = s.defaultCell()
	}
}

func (s *Screen) DeleteCharacters(count int) {
	if count <= 0 {
		count = 1
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	for col := s.Cursor.Col; col < s.Columns; col++ {
		if col+count < s.Columns {
			line[col] = line[col+count]
		} else {
			line[col] = s.defaultCell()
		}
	}
}

func (s *Screen) EraseCharacters(count int) {
	if count <= 0 {
		count = 1
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	for col := s.Cursor.Col; col < s.Columns && col < s.Cursor.Col+count; col++ {
		line[col] = Cell{Data: " ", Attr: s.Cursor.Attr}
	}
}

func (s *Screen) EraseInLine(how int, private bool) {
	if private {
		return
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	var start, end int
	switch how {
	case 0:
		start = s.Cursor.Col
		end = s.Columns
	case 1:
		start = 0
		end = s.Cursor.Col + 1
	case 2:
		start = 0
		end = s.Columns
	}
	for col := start; col < end; col++ {
		line[col] = Cell{Data: " ", Attr: s.Cursor.Attr}
	}
}

func (s *Screen) EraseInDisplay(how int, private bool, _ ...int) {
	if private {
		return
	}
	var start, end int
	switch how {
	case 0:
		start = s.Cursor.Row + 1
		end = s.Lines
	case 1:
		start = 0
		end = s.Cursor.Row
	case 2, 3:
		start = 0
		end = s.Lines
	}
	if how == 0 || how == 1 {
		s.EraseInLine(how, false)
	}
	for row := start; row < end; row++ {
		line := s.Buffer[row]
		for col := 0; col < s.Columns; col++ {
			line[col] = Cell{Data: " ", Attr: s.Cursor.Attr}
		}
		s.Dirty[row] = struct{}{}
	}
	if how == 2 || how == 3 {
		s.markDirtyRange(0, s.Lines-1)
	}
}

func (s *Screen) SetTabStop() {
	s.TabStops[s.Cursor.Col] = struct{}{}
}

func (s *Screen) ClearTabStop(how int) {
	switch how {
	case 0:
		delete(s.TabStops, s.Cursor.Col)
	case 3:
		s.TabStops = make(map[int]struct{})
	}
}

func (s *Screen) EnsureCursor() {
	s.ensureHB()
	s.ensureVB(false)
}

func (s *Screen) CursorUp(count int) {
	if count <= 0 {
		count = 1
	}
	top, _ := s.scrollRegion()
	s.Cursor.Row = maxInt(s.Cursor.Row-count, top)
}

func (s *Screen) CursorUp1(count int) {
	s.CursorUp(count)
	s.CarriageReturn()
}

func (s *Screen) CursorDown(count int) {
	if count <= 0 {
		count = 1
	}
	_, bottom := s.scrollRegion()
	s.Cursor.Row = minInt(s.Cursor.Row+count, bottom)
}

func (s *Screen) CursorDown1(count int) {
	s.CursorDown(count)
	s.CarriageReturn()
}

func (s *Screen) CursorBack(count int) {
	if s.Cursor.Col == s.Columns {
		s.Cursor.Col--
	}
	if count <= 0 {
		count = 1
	}
	s.Cursor.Col -= count
	s.ensureHB()
}

func (s *Screen) CursorForward(count int) {
	if count <= 0 {
		count = 1
	}
	s.Cursor.Col += count
	s.ensureHB()
}

func (s *Screen) CursorPosition(line, column int) {
	if column <= 0 {
		column = 1
	}
	if line <= 0 {
		line = 1
	}
	row := line - 1
	col := column - 1
	if s.Margins != nil && s.isModeSet(ModeDECOM) {
		row += s.Margins.Top
		if row < s.Margins.Top || row > s.Margins.Bottom {
			return
		}
	}
	s.Cursor.Row = row
	s.Cursor.Col = col
	s.ensureHB()
	s.ensureVB(false)
}

func (s *Screen) CursorToColumn(column int) {
	if column <= 0 {
		column = 1
	}
	s.Cursor.Col = column - 1
	s.ensureHB()
}

func (s *Screen) CursorToLine(line int) {
	if line <= 0 {
		line = 1
	}
	row := line - 1
	if s.isModeSet(ModeDECOM) && s.Margins != nil {
		row += s.Margins.Top
	}
	s.Cursor.Row = row
	s.ensureVB(false)
}

func (s *Screen) Bell() {
}

func (s *Screen) AlignmentDisplay() {
	s.markDirtyRange(0, s.Lines-1)
	for row := 0; row < s.Lines; row++ {
		for col := 0; col < s.Columns; col++ {
			cell := s.Buffer[row][col]
			cell.Data = "E"
			s.Buffer[row][col] = cell
		}
	}
}

func (s *Screen) SelectGraphicRendition(attrs []int, private bool) {
	if private {
		return
	}
	if len(attrs) == 0 || (len(attrs) == 1 && attrs[0] == 0) {
		s.Cursor.Attr = s.defaultAttr()
		return
	}

	replace := s.Cursor.Attr
	queue := append([]int(nil), attrs...)

	for len(queue) > 0 {
		attr := queue[0]
		queue = queue[1:]
		switch {
		case attr == 0:
			replace = s.defaultAttr()
		case fgANSI[attr] != "":
			replace.Fg = colorFromName(fgANSI[attr], ColorANSI16, uint8(attr-30))
		case bgANSI[attr] != "":
			replace.Bg = colorFromName(bgANSI[attr], ColorANSI16, uint8(attr-40))
		case fgAixterm[attr] != "":
			replace.Fg = colorFromName(fgAixterm[attr], ColorANSI16, uint8(attr-90+8))
		case bgAixterm[attr] != "":
			replace.Bg = colorFromName(bgAixterm[attr], ColorANSI16, uint8(attr-100+8))
		case textAttributes[attr] != "":
			flag := textAttributes[attr]
			switch flag {
			case "+bold":
				replace.Bold = true
			case "-bold":
				replace.Bold = false
			case "+italics":
				replace.Italics = true
			case "-italics":
				replace.Italics = false
			case "+underline":
				replace.Underline = true
			case "-underline":
				replace.Underline = false
			case "+blink":
				replace.Blink = true
			case "-blink":
				replace.Blink = false
			case "+reverse":
				replace.Reverse = true
			case "-reverse":
				replace.Reverse = false
			case "+strikethrough":
				replace.Strikethrough = true
			case "-strikethrough":
				replace.Strikethrough = false
			}
		case attr == SgrFg256 || attr == SgrBg256:
			keyIsFg := attr == SgrFg256
			if len(queue) < 1 {
				break
			}
			mode := queue[0]
			queue = queue[1:]
			switch mode {
			case 5:
				if len(queue) < 1 {
					break
				}
				index := queue[0]
				queue = queue[1:]
				if index >= 0 && index < len(fgBg256) {
					color := Color{Mode: ColorANSI256, Index: uint8(index), Name: fgBg256[index]}
					if keyIsFg {
						replace.Fg = color
					} else {
						replace.Bg = color
					}
				}
			case 2:
				if len(queue) < 3 {
					break
				}
				r := queue[0]
				g := queue[1]
				b := queue[2]
				queue = queue[3:]
				color := Color{Mode: ColorTrueColor, Name: rgbHex(r, g, b)}
				if keyIsFg {
					replace.Fg = color
				} else {
					replace.Bg = color
				}
			}
		}
	}
	s.Cursor.Attr = replace
}

func (s *Screen) ReportDeviceAttributes(mode int, private bool) {
	if mode == 0 && !private {
		s.WriteProcessInput(ControlCSI + "?6c")
	}
}

func (s *Screen) ReportDeviceStatus(mode int) {
	switch mode {
	case 5:
		s.WriteProcessInput(ControlCSI + "0n")
	case 6:
		x := s.Cursor.Col + 1
		y := s.Cursor.Row + 1
		if s.isModeSet(ModeDECOM) && s.Margins != nil {
			y -= s.Margins.Top
		}
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("%d;%dR", y, x))
	}
}

func (s *Screen) Debug(_ ...interface{}) {
}

func (s *Screen) DefineCharset(code, mode string) {
	mapping, ok := charsetMaps[code]
	if !ok {
		return
	}
	if mode == "(" {
		s.G0 = mapping
	} else if mode == ")" {
		s.G1 = mapping
	}
}

func (s *Screen) ShiftIn() {
	s.Charset = 0
}

func (s *Screen) ShiftOut() {
	s.Charset = 1
}

func (s *Screen) SetMargins(top, bottom int) {
	if (top == 0 || top == -1) && bottom == 0 {
		s.Margins = nil
		return
	}
	margins := s.Margins
	if margins == nil {
		margins = &Margins{Top: 0, Bottom: s.Lines - 1}
	}
	if top == 0 {
		top = margins.Top + 1
	}
	if bottom == 0 {
		bottom = margins.Bottom + 1
	}

	top = maxInt(0, minInt(top-1, s.Lines-1))
	bottom = maxInt(0, minInt(bottom-1, s.Lines-1))

	if bottom-top >= 1 {
		s.Margins = &Margins{Top: top, Bottom: bottom}
		s.CursorPosition(0, 0)
	}
}

func (s *Screen) SetMode(modes []int, private bool) {
	for _, mode := range modes {
		s.applySetMode(mode, private)
	}
}

func (s *Screen) ResetMode(modes []int, private bool) {
	for _, mode := range modes {
		s.applyResetMode(mode, private)
	}
}

func (s *Screen) applySetMode(mode int, private bool) {
	if private {
		mode = mode << 5
	}
	if s.Mode == nil {
		s.Mode = make(map[int]struct{})
	}
	s.Mode[mode] = struct{}{}

	if mode == ModeDECCOLM {
		saved := s.Columns
		s.SavedColumns = &saved
		s.Resize(s.Lines, 132)
		s.EraseInDisplay(2, false)
		s.CursorPosition(0, 0)
	}

	if mode == ModeDECOM {
		s.CursorPosition(0, 0)
	}

	if mode == ModeDECSCNM {
		s.markDirtyRange(0, s.Lines-1)
		for row := 0; row < s.Lines; row++ {
			for col := 0; col < s.Columns; col++ {
				cell := s.Buffer[row][col]
				cell.Attr.Reverse = true
				s.Buffer[row][col] = cell
			}
		}
		s.SelectGraphicRendition([]int{7}, false)
	}

	if mode == ModeDECTCEM {
		s.Cursor.Hidden = false
	}
}

func (s *Screen) applyResetMode(mode int, private bool) {
	if private {
		mode = mode << 5
	}
	delete(s.Mode, mode)

	if mode == ModeDECCOLM {
		if s.Columns == 132 && s.SavedColumns != nil {
			s.Resize(s.Lines, *s.SavedColumns)
			s.SavedColumns = nil
		}
		s.EraseInDisplay(2, false)
		s.CursorPosition(0, 0)
	}

	if mode == ModeDECOM {
		s.CursorPosition(0, 0)
	}

	if mode == ModeDECSCNM {
		s.markDirtyRange(0, s.Lines-1)
		for row := 0; row < s.Lines; row++ {
			for col := 0; col < s.Columns; col++ {
				cell := s.Buffer[row][col]
				cell.Attr.Reverse = false
				s.Buffer[row][col] = cell
			}
		}
		s.SelectGraphicRendition([]int{27}, false)
	}

	if mode == ModeDECTCEM {
		s.Cursor.Hidden = true
	}
}

func (s *Screen) defaultAttr() Attr {
	reverse := s.isModeSet(ModeDECSCNM)
	return Attr{
		Fg:      Color{Name: "default", Mode: ColorDefault},
		Bg:      Color{Name: "default", Mode: ColorDefault},
		Reverse: reverse,
	}
}

func (s *Screen) defaultCell() Cell {
	return Cell{Data: " ", Attr: s.defaultAttr()}
}

func (s *Screen) scrollRegion() (int, int) {
	if s.Margins != nil {
		return s.Margins.Top, s.Margins.Bottom
	}
	return 0, s.Lines - 1
}

func (s *Screen) ensureHB() {
	if s.Cursor.Col < 0 {
		s.Cursor.Col = 0
	}
	if s.Cursor.Col >= s.Columns {
		s.Cursor.Col = s.Columns - 1
	}
}

func (s *Screen) ensureVB(useMargins bool) {
	if (useMargins || s.isModeSet(ModeDECOM)) && s.Margins != nil {
		if s.Cursor.Row < s.Margins.Top {
			s.Cursor.Row = s.Margins.Top
		}
		if s.Cursor.Row > s.Margins.Bottom {
			s.Cursor.Row = s.Margins.Bottom
		}
		return
	}
	if s.Cursor.Row < 0 {
		s.Cursor.Row = 0
	}
	if s.Cursor.Row >= s.Lines {
		s.Cursor.Row = s.Lines - 1
	}
}

func (s *Screen) isModeSet(mode int) bool {
	_, ok := s.Mode[mode]
	return ok
}

func (s *Screen) translate(data string) string {
	mapping := s.G0
	if s.Charset == 1 {
		mapping = s.G1
	}
	var b strings.Builder
	for _, r := range data {
		if r >= 0 && r < 256 {
			b.WriteRune(mapping[r])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func sortedStops(stops map[int]struct{}) []int {
	values := make([]int, 0, len(stops))
	for value := range stops {
		values = append(values, value)
	}
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
	return values
}

func colorFromName(name string, mode ColorMode, index uint8) Color {
	if name == "default" {
		return Color{Name: name, Mode: ColorDefault}
	}
	return Color{Name: name, Mode: mode, Index: index}
}

func blankLine(cols int, cell Cell) []Cell {
	line := make([]Cell, cols)
	for col := 0; col < cols; col++ {
		line[col] = cell
	}
	return line
}

func makeBlankCells(lines, cols int) [][]Cell {
	cells := make([][]Cell, lines)
	cell := Cell{Data: " ", Attr: Attr{Fg: Color{Name: "default", Mode: ColorDefault}, Bg: Color{Name: "default", Mode: ColorDefault}}}
	for row := 0; row < lines; row++ {
		cells[row] = blankLine(cols, cell)
	}
	return cells
}

func (s *Screen) markDirtyRange(start, end int) {
	if start < 0 {
		start = 0
	}
	if end >= s.Lines {
		end = s.Lines - 1
	}
	for row := start; row <= end; row++ {
		s.Dirty[row] = struct{}{}
	}
}

func isCombiningCluster(cluster string) bool {
	if cluster == "" {
		return false
	}
	for _, r := range cluster {
		if !unicode.Is(unicode.Mn, r) && !unicode.Is(unicode.Me, r) {
			return false
		}
	}
	return true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
