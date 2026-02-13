package te

import (
	"fmt"
	"strconv"
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
	WrapNext bool
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
	lastDrawn         string
	leftMargin        int
	rightMargin       int
	lineWrapped       map[int]bool
	wrapNext          bool
	savedModes        map[int]bool
	selectionData     map[string]string
	colorPalette      map[int]string
	dynamicColors     map[int]string
	specialColors     map[int]string
	titleHexInput     bool
	titleHexOutput    bool
	conformanceLevel  int
	titleStack        []string
	iconStack         []string
	windowPosX        int
	windowPosY        int
	windowPixelWidth  int
	windowPixelHeight int
	screenPixelWidth  int
	screenPixelHeight int
	charPixelWidth    int
	charPixelHeight   int
	windowIconified   bool
}

func NewScreen(cols, lines int) *Screen {
	s := &Screen{}
	s.Resize(lines, cols)
	s.Reset()
	return s
}

func (s *Screen) Reset() {
	if s.SavedColumns != nil {
		s.Resize(s.Lines, *s.SavedColumns)
		s.SavedColumns = nil
	}
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
	s.SavedColumns = nil
	s.leftMargin = 0
	s.rightMargin = s.Columns - 1
	s.lineWrapped = make(map[int]bool)
	s.wrapNext = false
	s.savedModes = make(map[int]bool)
	s.selectionData = make(map[string]string)
	s.colorPalette = make(map[int]string)
	s.dynamicColors = make(map[int]string)
	s.specialColors = make(map[int]string)
	s.titleHexInput = false
	s.titleHexOutput = false
	s.conformanceLevel = 4
	s.titleStack = nil
	s.iconStack = nil
	if s.charPixelWidth == 0 {
		s.charPixelWidth = 8
	}
	if s.charPixelHeight == 0 {
		s.charPixelHeight = 16
	}
	s.windowPixelWidth = s.Columns * s.charPixelWidth
	s.windowPixelHeight = s.Lines * s.charPixelHeight
	s.screenPixelWidth = s.windowPixelWidth
	s.screenPixelHeight = s.windowPixelHeight
	s.windowPosX = 0
	s.windowPosY = 0
	s.windowIconified = false
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
		if s.charPixelWidth == 0 {
			s.charPixelWidth = 8
		}
		if s.charPixelHeight == 0 {
			s.charPixelHeight = 16
		}
		s.windowPixelWidth = s.Columns * s.charPixelWidth
		s.windowPixelHeight = s.Lines * s.charPixelHeight
		s.screenPixelWidth = s.windowPixelWidth
		s.screenPixelHeight = s.windowPixelHeight
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
	s.leftMargin = 0
	s.rightMargin = s.Columns - 1
	s.lineWrapped = make(map[int]bool)
	s.wrapNext = false
	s.savedModes = make(map[int]bool)
	s.selectionData = make(map[string]string)
	s.colorPalette = make(map[int]string)
	s.dynamicColors = make(map[int]string)
	s.specialColors = make(map[int]string)
	s.titleHexInput = false
	s.titleHexOutput = false
	s.conformanceLevel = 4
	s.titleStack = nil
	s.iconStack = nil
	if s.charPixelWidth == 0 {
		s.charPixelWidth = 8
	}
	if s.charPixelHeight == 0 {
		s.charPixelHeight = 16
	}
	s.windowPixelWidth = s.Columns * s.charPixelWidth
	s.windowPixelHeight = s.Lines * s.charPixelHeight
	s.screenPixelWidth = s.windowPixelWidth
	s.screenPixelHeight = s.windowPixelHeight
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
loop:
	for graphemes.Next() {
		cluster := graphemes.Str()
		width := runewidth.StringWidth(cluster)
		limit := s.Columns - 1
		if s.isModeSet(ModeDECLRMM) {
			limit = s.rightMargin
		}

		if s.wrapNext && width > 0 {
			s.wrapNext = false
			s.Dirty[s.Cursor.Row] = struct{}{}
			s.CarriageReturn()
			s.LineFeed()
		}
		if width > 0 && s.Cursor.Col > limit {
			s.Cursor.Col = limit
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
			break loop
		}

		if width > 0 {
			s.lastDrawn = cluster
			nextCol := s.Cursor.Col + width
			if nextCol > limit {
				if s.isModeSet(ModeDECAWM) {
					s.lineWrapped[s.Cursor.Row] = true
					s.wrapNext = true
					s.Cursor.Col = limit + 1
					continue
				}
				s.Cursor.Col = limit + 1
			} else {
				s.Cursor.Col = nextCol
			}
		}
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
}

func (s *Screen) SetTitle(param string) {
	s.Title = s.applyTitleModes(param)
}

func (s *Screen) SetIconName(param string) {
	s.IconName = s.applyTitleModes(param)
}

func (s *Screen) CarriageReturn() {
	s.wrapNext = false
	if s.isModeSet(ModeDECLRMM) {
		if s.isModeSet(ModeDECOM) {
			s.Cursor.Col = s.leftMargin
			return
		}
		if s.Cursor.Col < s.leftMargin {
			s.Cursor.Col = 0
			return
		}
		s.Cursor.Col = s.leftMargin
		return
	}
	s.Cursor.Col = 0
}

func (s *Screen) Index() {
	s.wrapNext = false
	top, bottom := s.scrollRegion()
	if s.Cursor.Row == bottom {
		if !s.canScrollHorizontal() {
			return
		}
		s.markDirtyRange(0, s.Lines-1)
		for row := top; row < bottom; row++ {
			s.Buffer[row] = s.Buffer[row+1]
		}
		s.Buffer[bottom] = blankLine(s.Columns, s.defaultCell())
		return
	}
	if s.Cursor.Row < bottom {
		s.Cursor.Row++
		return
	}
	if s.Cursor.Row < s.Lines-1 {
		s.Cursor.Row++
	}
}

func (s *Screen) ReverseIndex() {
	s.wrapNext = false
	top, bottom := s.scrollRegion()
	if s.Cursor.Row == top {
		if !s.canScrollHorizontal() {
			return
		}
		s.markDirtyRange(0, s.Lines-1)
		for row := bottom; row > top; row-- {
			s.Buffer[row] = s.Buffer[row-1]
		}
		s.Buffer[top] = blankLine(s.Columns, s.defaultCell())
		return
	}
	if s.Cursor.Row > top {
		s.Cursor.Row--
		return
	}
	if s.Cursor.Row > 0 {
		s.Cursor.Row--
	}
}

func (s *Screen) LineFeed() {
	s.wrapNext = false
	if !s.canScrollHorizontal() {
		_, bottom := s.scrollRegion()
		if s.Cursor.Row < bottom {
			s.Cursor.Row++
		} else if s.Cursor.Row > bottom {
			if s.Cursor.Row < s.Lines-1 {
				s.Cursor.Row++
			}
		}
	} else {
		s.Index()
	}
	if s.isModeSet(ModeLNM) {
		s.CarriageReturn()
	}
}

func (s *Screen) Tab() {
	s.wrapNext = false
	limit := s.Columns - 1
	if s.isModeSet(ModeDECLRMM) {
		limit = s.rightMargin
	}
	column := limit
	for _, stop := range sortedStops(s.TabStops) {
		if s.Cursor.Col < stop {
			column = stop
			break
		}
	}
	if column > limit {
		column = limit
	}
	s.Cursor.Col = column
}

func (s *Screen) Backspace() {
	s.moveLeft(1, true)
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
		WrapNext: s.wrapNext,
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
	s.wrapNext = last.WrapNext
	s.Cursor = last.Cursor
	s.ensureHB()
	s.ensureVB(true)
}

func (s *Screen) InsertLines(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	left, right := s.horizontalMargins()
	if s.isModeSet(ModeDECLRMM) && (s.Cursor.Col < left || s.Cursor.Col > right) {
		return
	}
	s.markDirtyRange(s.Cursor.Row, s.Lines-1)
	if count > bottom-s.Cursor.Row+1 {
		count = bottom - s.Cursor.Row + 1
	}
	for row := bottom; row >= s.Cursor.Row; row-- {
		if row+count <= bottom {
			s.copyRowSegment(row, row+count, left, right)
		}
		s.clearRowSegment(row, left, right)
	}
	s.CarriageReturn()
}

func (s *Screen) DeleteLines(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	left, right := s.horizontalMargins()
	if s.isModeSet(ModeDECLRMM) && (s.Cursor.Col < left || s.Cursor.Col > right) {
		return
	}
	s.markDirtyRange(s.Cursor.Row, s.Lines-1)
	if count > bottom-s.Cursor.Row+1 {
		count = bottom - s.Cursor.Row + 1
	}
	for row := s.Cursor.Row; row <= bottom; row++ {
		if row+count <= bottom {
			s.copyRowSegment(row+count, row, left, right)
		} else {
			s.clearRowSegment(row, left, right)
		}
	}
	s.CarriageReturn()
}

func (s *Screen) InsertCharacters(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	left, right := s.horizontalMargins()
	if s.Cursor.Col < left || s.Cursor.Col > right {
		return
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	for col := right; col >= s.Cursor.Col; col-- {
		if col+count <= right {
			line[col+count] = line[col]
		}
	}
	for col := s.Cursor.Col; col < s.Cursor.Col+count && col <= right; col++ {
		line[col] = s.defaultCell()
	}
}

func (s *Screen) DeleteCharacters(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	left, right := s.horizontalMargins()
	if s.Cursor.Col < left || s.Cursor.Col > right {
		return
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	for col := s.Cursor.Col; col <= right; col++ {
		if col+count <= right {
			line[col] = line[col+count]
		} else {
			line[col] = s.defaultCell()
		}
	}
}

func (s *Screen) EraseCharacters(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	left, right := s.horizontalMargins()
	if s.Cursor.Col < left || s.Cursor.Col > right {
		return
	}
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	end := minInt(s.Cursor.Col+count-1, right)
	for col := s.Cursor.Col; col <= end; col++ {
		line[col] = Cell{Data: " ", Attr: s.Cursor.Attr}
	}
}

func (s *Screen) EraseInLine(how int, private bool, _ ...int) {
	s.Dirty[s.Cursor.Row] = struct{}{}
	line := s.Buffer[s.Cursor.Row]
	left, right := s.horizontalMargins()
	var start, end int
	switch how {
	case 0:
		start = maxInt(s.Cursor.Col, left)
		end = right + 1
	case 1:
		start = left
		end = minInt(s.Cursor.Col, right) + 1
	case 2:
		start = left
		end = right + 1
	}
	for col := start; col < end; col++ {
		line[col] = Cell{Data: " ", Attr: s.Cursor.Attr}
	}
}

func (s *Screen) EraseInDisplay(how int, private bool, _ ...int) {
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

func (s *Screen) ClearTabStop(how ...int) {
	value := 0
	if len(how) > 0 {
		value = how[0]
	}
	switch value {
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

func (s *Screen) CursorUp(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.wrapNext = false
	if count <= 0 {
		count = 1
	}
	top, _ := s.scrollRegion()
	s.Cursor.Row = maxInt(s.Cursor.Row-count, top)
}

func (s *Screen) CursorUp1(params ...int) {
	s.wrapNext = false
	s.CursorUp(params...)
	s.CarriageReturn()
}

func (s *Screen) CursorDown(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.wrapNext = false
	if count <= 0 {
		count = 1
	}
	_, bottom := s.scrollRegion()
	s.Cursor.Row = minInt(s.Cursor.Row+count, bottom)
}

func (s *Screen) CursorDown1(params ...int) {
	s.wrapNext = false
	s.CursorDown(params...)
	s.CarriageReturn()
}

func (s *Screen) CursorBack(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.moveLeft(count, true)
}

func (s *Screen) CursorBackTab(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.wrapNext = false
	if count <= 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		prev := 0
		for stop := range s.TabStops {
			if stop < s.Cursor.Col && stop > prev {
				prev = stop
			}
		}
		s.Cursor.Col = prev
	}
}

func (s *Screen) CursorForwardTab(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.wrapNext = false
	if count <= 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		s.Tab()
	}
}

func (s *Screen) NextLine() {
	s.wrapNext = false
	origCol := s.Cursor.Col
	s.CarriageReturn()
	if s.isModeSet(ModeDECLRMM) {
		left, right := s.horizontalMargins()
		if origCol < left || origCol > right {
			_, bottom := s.scrollRegion()
			if s.Cursor.Row < bottom {
				s.Cursor.Row++
			} else if s.Cursor.Row > bottom {
				if s.Cursor.Row < s.Lines-1 {
					s.Cursor.Row++
				}
			}
			return
		}
	}
	s.Index()
}

func (s *Screen) moveLeft(count int, reverseWrap bool) {
	if count <= 0 {
		count = 1
	}
	left, right := s.horizontalMargins()
	reverseInline := s.isModeSet(ModeReverseWrapInline)
	reverseExtend := s.isModeSet(ModeReverseWrapExtend)
	for i := 0; i < count; i++ {
		if s.wrapNext {
			s.wrapNext = false
			if reverseWrap && s.isModeSet(ModeDECAWM) && (reverseInline || reverseExtend) {
				if s.Cursor.Col == right+1 {
					s.Cursor.Col = right
				}
				continue
			}
			if s.Cursor.Col == right+1 {
				s.Cursor.Col = maxInt(left, right-1)
				continue
			}
			if s.Cursor.Col > 0 {
				s.Cursor.Col--
			}
			continue
		}
		if s.Cursor.Col == right+1 {
			s.Cursor.Col = right
			continue
		}
		if s.Cursor.Col == 0 {
			if reverseWrap && s.isModeSet(ModeDECAWM) && (reverseInline || reverseExtend) {
				s.reverseWrapToPreviousLine(left, right, reverseInline, reverseExtend)
			}
			continue
		}
		if s.isModeSet(ModeDECLRMM) && s.Cursor.Col < left {
			s.Cursor.Col--
			continue
		}
		if s.Cursor.Col > left {
			s.Cursor.Col--
			continue
		}
		if s.Cursor.Col < left {
			s.Cursor.Col--
			continue
		}
		if !reverseWrap || !s.isModeSet(ModeDECAWM) || (!reverseInline && !reverseExtend) {
			continue
		}
		s.reverseWrapToPreviousLine(left, right, reverseInline, reverseExtend)
	}
}

func (s *Screen) reverseWrapToPreviousLine(left, right int, reverseInline, reverseExtend bool) {
	top, bottom := s.scrollRegion()
	targetRow := s.Cursor.Row - 1
	if reverseExtend && s.Cursor.Row == top {
		targetRow = bottom
	}
	if targetRow < 0 || targetRow >= s.Lines {
		return
	}
	if reverseInline && !reverseExtend {
		if !s.lineWrapped[targetRow] {
			return
		}
	}
	s.Cursor.Row = targetRow
	s.Cursor.Col = right
}

func (s *Screen) CursorForward(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	s.wrapNext = false
	if count <= 0 {
		count = 1
	}
	left, right := s.horizontalMargins()
	if s.isModeSet(ModeDECLRMM) && s.Cursor.Col >= left && s.Cursor.Col <= right {
		s.Cursor.Col = minInt(s.Cursor.Col+count, right)
		return
	}
	s.Cursor.Col = minInt(s.Cursor.Col+count, s.Columns-1)
}

func (s *Screen) ScrollUp(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if count > bottom-top+1 {
		count = bottom - top + 1
	}
	left, right := s.horizontalMargins()
	s.markDirtyRange(top, bottom)
	for i := 0; i < count; i++ {
		for row := top; row < bottom; row++ {
			s.copyRowSegment(row+1, row, left, right)
		}
		s.clearRowSegment(bottom, left, right)
	}
}

func (s *Screen) ScrollDown(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	if count > bottom-top+1 {
		count = bottom - top + 1
	}
	left, right := s.horizontalMargins()
	s.markDirtyRange(top, bottom)
	for i := 0; i < count; i++ {
		for row := bottom; row > top; row-- {
			s.copyRowSegment(row-1, row, left, right)
		}
		s.clearRowSegment(top, left, right)
	}
}

func (s *Screen) RepeatLast(params ...int) {
	count := 1
	if len(params) > 0 {
		count = params[0]
	}
	if count <= 0 {
		count = 1
	}
	if s.lastDrawn == "" {
		return
	}
	for i := 0; i < count; i++ {
		s.Draw(s.lastDrawn)
	}
}

func (s *Screen) CursorPosition(params ...int) {
	line := 0
	column := 0
	if len(params) > 0 {
		line = params[0]
	}
	if len(params) > 1 {
		column = params[1]
	}
	s.wrapNext = false
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
	if s.isModeSet(ModeDECOM) && s.isModeSet(ModeDECLRMM) {
		col += s.leftMargin
		if col < s.leftMargin || col > s.rightMargin {
			return
		}
	}
	s.Cursor.Row = row
	s.Cursor.Col = col
	s.ensureHB()
	s.ensureVB(false)
}

func (s *Screen) CursorToColumn(column ...int) {
	value := 1
	if len(column) > 0 {
		value = column[0]
	}
	s.wrapNext = false
	if value <= 0 {
		value = 1
	}
	col := value - 1
	if s.isModeSet(ModeDECOM) && s.isModeSet(ModeDECLRMM) {
		col += s.leftMargin
	}
	s.Cursor.Col = col
	s.ensureHB()
}

func (s *Screen) CursorToColumnAbsolute(column ...int) {
	value := 1
	if len(column) > 0 {
		value = column[0]
	}
	s.wrapNext = false
	if value <= 0 {
		value = 1
	}
	col := value - 1
	if col < 0 {
		col = 0
	}
	if col >= s.Columns {
		col = s.Columns - 1
	}
	s.Cursor.Col = col
}

func (s *Screen) CursorToLine(params ...int) {
	line := 1
	if len(params) > 0 {
		line = params[0]
	}
	s.wrapNext = false
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

func (s *Screen) ReportDeviceAttributes(mode int, private bool, prefix rune, _ ...int) {
	if mode != 0 {
		return
	}
	if prefix == '>' {
		s.WriteProcessInput(ControlCSI + ">0;0;0c")
		return
	}
	if !private {
		s.WriteProcessInput(ControlCSI + "?6c")
		return
	}
	if prefix == '?' {
		s.WriteProcessInput(ControlCSI + "?6c")
	}
}

func (s *Screen) ReportDeviceStatus(mode int, private bool, prefix rune, _ ...int) {
	if private && prefix == '?' {
		if mode == 6 {
			x := s.Cursor.Col + 1
			y := s.Cursor.Row + 1
			if s.isModeSet(ModeDECOM) && s.Margins != nil {
				y -= s.Margins.Top
			}
			if s.isModeSet(ModeDECOM) && s.isModeSet(ModeDECLRMM) {
				x -= s.leftMargin
			}
			s.WriteProcessInput(ControlCSI + fmt.Sprintf("?%d;%d;1R", y, x))
		}
		return
	}
	switch mode {
	case 5:
		s.WriteProcessInput(ControlCSI + "0n")
	case 6:
		x := s.Cursor.Col + 1
		y := s.Cursor.Row + 1
		if s.isModeSet(ModeDECOM) && s.Margins != nil {
			y -= s.Margins.Top
		}
		if s.isModeSet(ModeDECOM) && s.isModeSet(ModeDECLRMM) {
			x -= s.leftMargin
		}
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("%d;%dR", y, x))
	}
}

func (s *Screen) ReportMode(mode int, private bool) {
	check := mode
	prefix := ""
	if private {
		check = mode << 5
		prefix = "?"
	}
	status := 2
	if s.Mode != nil {
		if _, ok := s.Mode[check]; ok {
			status = 1
		}
	}
	s.WriteProcessInput(ControlCSI + fmt.Sprintf("%s%d;%d$y", prefix, mode, status))
}

func (s *Screen) RequestStatusString(query string) {
	if query == "" {
		return
	}
	response := ""
	switch query {
	case "m":
		codes := []string{"0"}
		if s.Cursor.Attr.Bold {
			codes = append(codes, "1")
		}
		if s.Cursor.Attr.Italics {
			codes = append(codes, "3")
		}
		if s.Cursor.Attr.Underline {
			codes = append(codes, "4")
		}
		if s.Cursor.Attr.Blink {
			codes = append(codes, "5")
		}
		if s.Cursor.Attr.Reverse {
			codes = append(codes, "7")
		}
		if s.Cursor.Attr.Conceal {
			codes = append(codes, "8")
		}
		if s.Cursor.Attr.Strikethrough {
			codes = append(codes, "9")
		}
		response = "1$r" + strings.Join(codes, ";") + "m"
	case "r":
		top, bottom := s.scrollRegion()
		response = fmt.Sprintf("1$r%d;%dr", top+1, bottom+1)
	case "s":
		if s.isModeSet(ModeDECLRMM) {
			response = fmt.Sprintf("1$r%d;%ds", s.leftMargin+1, s.rightMargin+1)
		} else {
			response = fmt.Sprintf("1$r1;%ds", s.Columns)
		}
	case " q":
		response = "1$r0 q"
	default:
		return
	}
	s.WriteProcessInput(ControlDCS + response + ControlST)
}

func (s *Screen) SoftReset() {
	s.wrapNext = false
	s.Cursor.Hidden = false
	s.Margins = nil
	s.leftMargin = 0
	s.rightMargin = s.Columns - 1
	s.Mode = map[int]struct{}{ModeDECAWM: {}, ModeDECTCEM: {}}
	s.Cursor.Attr = s.defaultAttr()
	s.G0 = charsetLat1
	s.G1 = charsetVT100
	s.Charset = 0
	s.Savepoints = nil
	s.savedModes = make(map[int]bool)
	s.selectionData = make(map[string]string)
	s.colorPalette = make(map[int]string)
	s.dynamicColors = make(map[int]string)
	s.specialColors = make(map[int]string)
	s.titleHexInput = false
	s.titleHexOutput = false
	s.conformanceLevel = 4
	s.titleStack = nil
	s.iconStack = nil
}

func (s *Screen) SaveModes(modes []int) {
	if s.savedModes == nil {
		s.savedModes = make(map[int]bool)
	}
	for _, mode := range modes {
		key := mode << 5
		_, ok := s.Mode[key]
		s.savedModes[key] = ok
	}
}

func (s *Screen) RestoreModes(modes []int) {
	for _, mode := range modes {
		key := mode << 5
		state, ok := s.savedModes[key]
		if !ok {
			continue
		}
		if state {
			s.applySetMode(mode, true)
		} else {
			s.applyResetMode(mode, true)
		}
	}
}

func (s *Screen) ForwardIndex() {
	left, right := s.horizontalMargins()
	top, bottom := s.scrollRegion()
	if s.isModeSet(ModeDECLRMM) {
		if s.Cursor.Col < left || s.Cursor.Col > right {
			if s.Cursor.Col < s.Columns-1 {
				s.Cursor.Col++
			}
			return
		}
	}
	if s.Cursor.Col < right {
		s.Cursor.Col++
		return
	}
	s.shiftHorizontal(left, right, top, bottom, -1)
}

func (s *Screen) BackIndex() {
	left, right := s.horizontalMargins()
	top, bottom := s.scrollRegion()
	if s.isModeSet(ModeDECLRMM) {
		if s.Cursor.Col < left || s.Cursor.Col > right {
			if s.Cursor.Col > 0 {
				s.Cursor.Col--
			}
			return
		}
	}
	if s.Cursor.Col > left {
		s.Cursor.Col--
		return
	}
	s.shiftHorizontal(left, right, top, bottom, 1)
}

func (s *Screen) InsertColumns(count int) {
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	left, right := s.horizontalMargins()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	if s.isModeSet(ModeDECLRMM) && (s.Cursor.Col < left || s.Cursor.Col > right) {
		return
	}
	if right < left {
		return
	}
	if count > right-s.Cursor.Col+1 {
		count = right - s.Cursor.Col + 1
	}
	s.markDirtyRange(top, bottom)
	for row := top; row <= bottom; row++ {
		for col := right; col >= s.Cursor.Col+count; col-- {
			s.Buffer[row][col] = s.Buffer[row][col-count]
		}
		for col := s.Cursor.Col; col < s.Cursor.Col+count && col <= right; col++ {
			s.Buffer[row][col] = s.defaultCell()
		}
	}
}

func (s *Screen) DeleteColumns(count int) {
	if count <= 0 {
		count = 1
	}
	top, bottom := s.scrollRegion()
	left, right := s.horizontalMargins()
	if s.Cursor.Row < top || s.Cursor.Row > bottom {
		return
	}
	if s.isModeSet(ModeDECLRMM) && (s.Cursor.Col < left || s.Cursor.Col > right) {
		return
	}
	if right < left {
		return
	}
	if count > right-s.Cursor.Col+1 {
		count = right - s.Cursor.Col + 1
	}
	s.markDirtyRange(top, bottom)
	for row := top; row <= bottom; row++ {
		for col := s.Cursor.Col; col <= right-count; col++ {
			s.Buffer[row][col] = s.Buffer[row][col+count]
		}
		for col := right - count + 1; col <= right; col++ {
			s.Buffer[row][col] = s.defaultCell()
		}
	}
}

func (s *Screen) EraseRectangle(top, left, bottom, right int) {
	s.fillRectangle(top, left, bottom, right, " ")
}

func (s *Screen) FillRectangle(ch rune, top, left, bottom, right int) {
	if ch == 0 {
		return
	}
	s.fillRectangle(top, left, bottom, right, string(ch))
}

func (s *Screen) fillRectangle(top, left, bottom, right int, data string) {
	if top <= 0 {
		top = 1
	}
	if left <= 0 {
		left = 1
	}
	if bottom <= 0 {
		bottom = s.Lines
	}
	if right <= 0 {
		right = s.Columns
	}
	if s.isModeSet(ModeDECOM) && s.Margins != nil {
		top += s.Margins.Top
		bottom += s.Margins.Top
		if s.isModeSet(ModeDECLRMM) {
			left += s.leftMargin
			right += s.leftMargin
		}
	}
	if top > bottom || left > right {
		return
	}
	if top < 1 {
		top = 1
	}
	if left < 1 {
		left = 1
	}
	if bottom > s.Lines {
		bottom = s.Lines
	}
	if right > s.Columns {
		right = s.Columns
	}
	for row := top - 1; row <= bottom-1; row++ {
		for col := left - 1; col <= right-1; col++ {
			s.Buffer[row][col] = Cell{Data: data, Attr: s.defaultAttr()}
		}
		s.Dirty[row] = struct{}{}
	}
}

func (s *Screen) CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft int) {
	if srcTop <= 0 {
		srcTop = 1
	}
	if srcLeft <= 0 {
		srcLeft = 1
	}
	if srcBottom <= 0 {
		srcBottom = s.Lines
	}
	if srcRight <= 0 {
		srcRight = s.Columns
	}
	if dstTop <= 0 {
		dstTop = 1
	}
	if dstLeft <= 0 {
		dstLeft = 1
	}
	if srcTop > srcBottom || srcLeft > srcRight {
		return
	}
	if srcTop < 1 {
		srcTop = 1
	}
	if srcLeft < 1 {
		srcLeft = 1
	}
	if srcBottom > s.Lines {
		srcBottom = s.Lines
	}
	if srcRight > s.Columns {
		srcRight = s.Columns
	}
	height := srcBottom - srcTop + 1
	width := srcRight - srcLeft + 1
	if height <= 0 || width <= 0 {
		return
	}
	temp := make([][]Cell, height)
	for r := 0; r < height; r++ {
		temp[r] = make([]Cell, width)
		copy(temp[r], s.Buffer[srcTop-1+r][srcLeft-1:srcLeft-1+width])
	}
	for r := 0; r < height; r++ {
		dstRow := dstTop - 1 + r
		if dstRow < 0 || dstRow >= s.Lines {
			continue
		}
		for c := 0; c < width; c++ {
			dstCol := dstLeft - 1 + c
			if dstCol < 0 || dstCol >= s.Columns {
				continue
			}
			s.Buffer[dstRow][dstCol] = temp[r][c]
		}
		s.Dirty[dstRow] = struct{}{}
	}
}

func (s *Screen) SetSelectionData(selection, data string) {
	if selection == "" {
		selection = "s0"
	}
	if s.selectionData == nil {
		s.selectionData = make(map[string]string)
	}
	s.selectionData[selection] = data
}

func (s *Screen) QuerySelectionData(selection string) {
	if selection == "" {
		selection = "s0"
	}
	data := ""
	if s.selectionData != nil {
		if value, ok := s.selectionData[selection]; ok {
			data = value
		}
	}
	s.WriteProcessInput(ControlOSC + "52;" + selection + ";" + data + ControlST)
}

func (s *Screen) SetColor(index int, value string) {
	normalized, ok := normalizeColorSpec(value)
	if !ok {
		return
	}
	if s.colorPalette == nil {
		s.colorPalette = make(map[int]string)
	}
	s.colorPalette[index] = normalized
}

func (s *Screen) QueryColor(index int) {
	value := "rgb:0000/0000/0000"
	if s.colorPalette != nil {
		if v, ok := s.colorPalette[index]; ok {
			value = v
		}
	}
	s.WriteProcessInput(ControlOSC + fmt.Sprintf("4;%d;%s", index, value) + ControlST)
}

func (s *Screen) ResetColor(index int, all bool) {
	if s.colorPalette == nil {
		return
	}
	if all {
		s.colorPalette = make(map[int]string)
		return
	}
	delete(s.colorPalette, index)
}

func (s *Screen) SetDynamicColor(index int, value string) {
	normalized, ok := normalizeColorSpec(value)
	if !ok {
		return
	}
	if s.dynamicColors == nil {
		s.dynamicColors = make(map[int]string)
	}
	s.dynamicColors[index] = normalized
}

func (s *Screen) QueryDynamicColor(index int) {
	value := "rgb:0000/0000/0000"
	if s.dynamicColors != nil {
		if v, ok := s.dynamicColors[index]; ok {
			value = v
		}
	}
	s.WriteProcessInput(ControlOSC + fmt.Sprintf("%d;%s", index, value) + ControlST)
}

func (s *Screen) SetSpecialColor(index int, value string) {
	normalized, ok := normalizeColorSpec(value)
	if !ok {
		return
	}
	if s.specialColors == nil {
		s.specialColors = make(map[int]string)
	}
	s.specialColors[index] = normalized
}

func (s *Screen) QuerySpecialColor(index int) {
	value := "rgb:0000/0000/0000"
	if s.specialColors != nil {
		if v, ok := s.specialColors[index]; ok {
			value = v
		}
	}
	s.WriteProcessInput(ControlOSC + fmt.Sprintf("5;%d;%s", index, value) + ControlST)
}

func (s *Screen) ResetSpecialColor(index int, all bool) {
	if s.specialColors == nil {
		return
	}
	if all {
		s.specialColors = make(map[int]string)
		return
	}
	delete(s.specialColors, index)
}

func normalizeColorSpec(spec string) (string, bool) {
	if strings.HasPrefix(spec, "rgb:") {
		parts := strings.Split(spec[4:], "/")
		if len(parts) != 3 {
			return "", false
		}
		return fmt.Sprintf("rgb:%s/%s/%s", normalizeHexComponent(parts[0]), normalizeHexComponent(parts[1]), normalizeHexComponent(parts[2])), true
	}
	if strings.HasPrefix(spec, "#") {
		hex := spec[1:]
		if len(hex)%3 != 0 {
			return "", false
		}
		size := len(hex) / 3
		if size < 1 || size > 4 {
			return "", false
		}
		return fmt.Sprintf("rgb:%s/%s/%s", normalizeHexComponent(hex[0:size]), normalizeHexComponent(hex[size:2*size]), normalizeHexComponent(hex[2*size:])), true
	}
	return "", false
}

func normalizeHexComponent(component string) string {
	value, err := strconv.ParseInt(component, 16, 32)
	if err != nil {
		return "0000"
	}
	byteValue := 0
	switch len(component) {
	case 1:
		byteValue = int(value) << 4
	case 2:
		byteValue = int(value)
	case 3:
		byteValue = int(value) >> 4
	default:
		byteValue = int(value) >> 8
	}
	if byteValue < 0 {
		byteValue = 0
	}
	if byteValue > 255 {
		byteValue = 255
	}
	full := (byteValue << 8) | byteValue
	return fmt.Sprintf("%04x", full)
}

func (s *Screen) applyTitleModes(param string) string {
	value := param
	if s.titleHexInput {
		decoded, ok := decodeHexString(param)
		if ok {
			value = decoded
		}
	}
	if s.titleHexOutput {
		value = encodeHexString(value)
	}
	return value
}

func decodeHexString(input string) (string, bool) {
	if len(input)%2 != 0 {
		return "", false
	}
	bytes := make([]byte, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		b, err := strconv.ParseUint(input[i:i+2], 16, 8)
		if err != nil {
			return "", false
		}
		bytes[i/2] = byte(b)
	}
	return string(bytes), true
}

func encodeHexString(input string) string {
	var b strings.Builder
	for _, r := range []byte(input) {
		fmt.Fprintf(&b, "%02x", r)
	}
	return b.String()
}

func (s *Screen) SetTitleMode(params []int, reset bool) {
	if len(params) == 0 {
		return
	}
	apply := func(param int) {
		switch param {
		case 0:
			if reset {
				s.titleHexInput = false
			} else {
				s.titleHexInput = true
			}
		case 2:
			s.titleHexInput = false
		case 1:
			if reset {
				s.titleHexOutput = false
			} else {
				s.titleHexOutput = true
			}
		case 3:
			s.titleHexOutput = false
		}
	}
	for _, param := range params {
		apply(param)
	}
}

func (s *Screen) SetConformance(level int, sevenBit int) {
	if level >= 60 {
		level -= 60
	}
	if level <= 0 {
		return
	}
	s.Reset()
	s.conformanceLevel = level
	_ = sevenBit
}

func (s *Screen) WindowOp(params []int) {
	if len(params) == 0 {
		return
	}
	switch params[0] {
	case 1:
		s.windowIconified = false
	case 2:
		s.windowIconified = true
	case 3:
		if len(params) > 1 && params[1] >= 0 {
			s.windowPosX = params[1]
		}
		if len(params) > 2 && params[2] >= 0 {
			s.windowPosY = params[2]
		}
	case 4:
		if len(params) > 1 {
			height := params[1]
			if height == 0 {
				height = s.screenPixelHeight
			}
			if height > 0 {
				s.windowPixelHeight = height
			}
		}
		if len(params) > 2 {
			width := params[2]
			if width == 0 {
				width = s.screenPixelWidth
			}
			if width > 0 {
				s.windowPixelWidth = width
			}
		}
	case 8:
		rows := s.Lines
		cols := s.Columns
		if len(params) > 1 && params[1] >= 0 {
			rows = params[1]
			if rows == 0 {
				rows = s.Lines
			}
		}
		if len(params) > 2 && params[2] >= 0 {
			cols = params[2]
			if cols == 0 {
				cols = s.Columns
			}
		}
		if rows > 0 && cols > 0 {
			s.Resize(rows, cols)
		}
	case 9:
		// maximize - no-op for now
	case 10:
		// fullscreen - no-op for now
	case 11:
		state := 1
		if s.windowIconified {
			state = 2
		}
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("%dt", state))
	case 13:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("3;%d;%dt", s.windowPosX, s.windowPosY))
	case 14:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("4;%d;%dt", s.windowPixelHeight, s.windowPixelWidth))
	case 15:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("5;%d;%dt", s.screenPixelHeight, s.screenPixelWidth))
	case 16:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("6;%d;%dt", s.charPixelHeight, s.charPixelWidth))
	case 18:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("8;%d;%dt", s.Lines, s.Columns))
	case 19:
		s.WriteProcessInput(ControlCSI + fmt.Sprintf("9;%d;%dt", s.Lines, s.Columns))
	case 20:
		s.WriteProcessInput(ControlOSC + "L" + s.IconName + ControlST)
	case 21:
		s.WriteProcessInput(ControlOSC + "l" + s.Title + ControlST)
	case 22:
		sub := 0
		if len(params) > 1 {
			sub = params[1]
		}
		switch sub {
		case 0:
			s.titleStack = append(s.titleStack, s.Title)
			s.iconStack = append(s.iconStack, s.IconName)
		case 1:
			s.iconStack = append(s.iconStack, s.IconName)
		case 2:
			s.titleStack = append(s.titleStack, s.Title)
		}
	case 23:
		sub := 0
		if len(params) > 1 {
			sub = params[1]
		}
		switch sub {
		case 0:
			if len(s.titleStack) > 0 {
				s.Title = s.titleStack[len(s.titleStack)-1]
				s.titleStack = s.titleStack[:len(s.titleStack)-1]
			}
			if len(s.iconStack) > 0 {
				s.IconName = s.iconStack[len(s.iconStack)-1]
				s.iconStack = s.iconStack[:len(s.iconStack)-1]
			}
		case 1:
			if len(s.iconStack) > 0 {
				s.IconName = s.iconStack[len(s.iconStack)-1]
				s.iconStack = s.iconStack[:len(s.iconStack)-1]
			}
		case 2:
			if len(s.titleStack) > 0 {
				s.Title = s.titleStack[len(s.titleStack)-1]
				s.titleStack = s.titleStack[:len(s.titleStack)-1]
			}
		}
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

func (s *Screen) SetMargins(params ...int) {
	top := 0
	bottom := 0
	if len(params) > 0 {
		top = params[0]
	}
	if len(params) > 1 {
		bottom = params[1]
	}
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

func (s *Screen) SetLeftRightMargins(left, right int) {
	if (left == 0 || left == -1) && right == 0 {
		s.leftMargin = 0
		s.rightMargin = s.Columns - 1
		s.CursorPosition(0, 0)
		return
	}
	if left == 0 {
		left = 1
	}
	if right == 0 {
		right = s.Columns
	}
	left = maxInt(1, minInt(left, s.Columns))
	right = maxInt(1, minInt(right, s.Columns))
	if right-left >= 1 {
		s.leftMargin = left - 1
		s.rightMargin = right - 1
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

	if mode == ModeDECSaveCursor {
		s.SaveCursor()
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

	if mode == ModeDECLRMM {
		s.leftMargin = 0
		s.rightMargin = s.Columns - 1
		s.CursorPosition(0, 0)
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

	if mode == ModeDECSaveCursor {
		s.RestoreCursor()
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

	if mode == ModeDECLRMM {
		s.leftMargin = 0
		s.rightMargin = s.Columns - 1
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

func (s *Screen) copyRowSegment(srcRow, dstRow, left, right int) {
	if srcRow < 0 || srcRow >= s.Lines || dstRow < 0 || dstRow >= s.Lines {
		return
	}
	for col := left; col <= right && col < s.Columns; col++ {
		if col < 0 {
			continue
		}
		s.Buffer[dstRow][col] = s.Buffer[srcRow][col]
	}
}

func (s *Screen) clearRowSegment(row, left, right int) {
	if row < 0 || row >= s.Lines {
		return
	}
	for col := left; col <= right && col < s.Columns; col++ {
		if col < 0 {
			continue
		}
		s.Buffer[row][col] = s.defaultCell()
	}
}

func (s *Screen) shiftHorizontal(left, right, top, bottom, direction int) {
	if left < 0 {
		left = 0
	}
	if right >= s.Columns {
		right = s.Columns - 1
	}
	if left > right {
		return
	}
	for row := top; row <= bottom && row < s.Lines; row++ {
		if row < 0 {
			continue
		}
		if direction < 0 {
			for col := left; col < right; col++ {
				s.Buffer[row][col] = s.Buffer[row][col+1]
			}
			s.Buffer[row][right] = s.defaultCell()
		} else {
			for col := right; col > left; col-- {
				s.Buffer[row][col] = s.Buffer[row][col-1]
			}
			s.Buffer[row][left] = s.defaultCell()
		}
	}
}

func (s *Screen) scrollRegion() (int, int) {
	if s.Margins != nil {
		return s.Margins.Top, s.Margins.Bottom
	}
	return 0, s.Lines - 1
}

func (s *Screen) horizontalMargins() (int, int) {
	if s.isModeSet(ModeDECLRMM) {
		return s.leftMargin, s.rightMargin
	}
	return 0, s.Columns - 1
}

func (s *Screen) canScrollHorizontal() bool {
	if !s.isModeSet(ModeDECLRMM) {
		return true
	}
	left, right := s.horizontalMargins()
	return s.Cursor.Col >= left && s.Cursor.Col <= right
}

func (s *Screen) ensureHB() {
	if s.isModeSet(ModeDECLRMM) {
		left, right := s.horizontalMargins()
		switch {
		case s.Cursor.Col < left:
			s.Cursor.Col = maxInt(0, s.Cursor.Col)
		case s.Cursor.Col > right:
			s.Cursor.Col = minInt(s.Cursor.Col, s.Columns-1)
		default:
			s.Cursor.Col = minInt(maxInt(s.Cursor.Col, left), right)
		}
		return
	}
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
