package te

import "math"

type historyDeque struct {
	items [][]Cell
	max   int
}

func (d *historyDeque) append(line []Cell) {
	if d.max == 0 {
		return
	}
	if len(d.items) == d.max {
		d.items = d.items[1:]
	}
	copyLine := append([]Cell(nil), line...)
	d.items = append(d.items, copyLine)
}

func (d *historyDeque) prepend(line []Cell) {
	if d.max == 0 {
		return
	}
	if len(d.items) == d.max {
		d.items = d.items[:len(d.items)-1]
	}
	copyLine := append([]Cell(nil), line...)
	d.items = append([][]Cell{copyLine}, d.items...)
}

func (d *historyDeque) pop() []Cell {
	if len(d.items) == 0 {
		return nil
	}
	last := d.items[len(d.items)-1]
	d.items = d.items[:len(d.items)-1]
	return last
}

func (d *historyDeque) popLeft() []Cell {
	if len(d.items) == 0 {
		return nil
	}
	first := d.items[0]
	d.items = d.items[1:]
	return first
}

type History struct {
	Top      historyDeque
	Bottom   historyDeque
	Ratio    float64
	Size     int
	Position int
}

type HistoryScreen struct {
	*Screen
	history History
}

func NewHistoryScreen(cols, lines, history int) *HistoryScreen {
	return NewHistoryScreenWithRatio(cols, lines, history, 0.5)
}

func NewHistoryScreenWithRatio(cols, lines, history int, ratio float64) *HistoryScreen {
	h := &HistoryScreen{}
	h.history = History{
		Top:      historyDeque{max: history},
		Bottom:   historyDeque{max: history},
		Ratio:    ratio,
		Size:     history,
		Position: history,
	}
	h.Screen = NewScreen(cols, lines)
	return h
}

func (h *HistoryScreen) beforeEvent(event string) {
	if event == "prev_page" || event == "next_page" {
		return
	}
	for h.history.Position < h.history.Size {
		h.NextPage()
	}
}

func (h *HistoryScreen) afterEvent(event string) {
	if event == "prev_page" || event == "next_page" {
		for row := range h.Buffer {
			if len(h.Buffer[row]) > h.Columns {
				h.Buffer[row] = h.Buffer[row][:h.Columns]
			}
		}
	}
	h.Cursor.Hidden = !(h.history.Position == h.history.Size && h.isModeSet(ModeDECTCEM))
}

func (h *HistoryScreen) Reset() {
	h.beforeEvent("reset")
	h.Screen.Reset()
	h.resetHistory()
	h.afterEvent("reset")
}

func (h *HistoryScreen) Bell() {
	h.beforeEvent("bell")
	h.Screen.Bell()
	h.afterEvent("bell")
}

func (h *HistoryScreen) Backspace() {
	h.beforeEvent("backspace")
	h.Screen.Backspace()
	h.afterEvent("backspace")
}

func (h *HistoryScreen) Tab() {
	h.beforeEvent("tab")
	h.Screen.Tab()
	h.afterEvent("tab")
}

func (h *HistoryScreen) CarriageReturn() {
	h.beforeEvent("carriage_return")
	h.Screen.CarriageReturn()
	h.afterEvent("carriage_return")
}

func (h *HistoryScreen) NextLine() {
	h.beforeEvent("next_line")
	h.Screen.NextLine()
	h.afterEvent("next_line")
}

func (h *HistoryScreen) ShiftOut() {
	h.beforeEvent("shift_out")
	h.Screen.ShiftOut()
	h.afterEvent("shift_out")
}

func (h *HistoryScreen) ShiftIn() {
	h.beforeEvent("shift_in")
	h.Screen.ShiftIn()
	h.afterEvent("shift_in")
}

func (h *HistoryScreen) SetTabStop() {
	h.beforeEvent("set_tab_stop")
	h.Screen.SetTabStop()
	h.afterEvent("set_tab_stop")
}

func (h *HistoryScreen) ClearTabStop(how int) {
	h.beforeEvent("clear_tab_stop")
	h.Screen.ClearTabStop(how)
	h.afterEvent("clear_tab_stop")
}

func (h *HistoryScreen) SaveCursor() {
	h.beforeEvent("save_cursor")
	h.Screen.SaveCursor()
	h.afterEvent("save_cursor")
}

func (h *HistoryScreen) RestoreCursor() {
	h.beforeEvent("restore_cursor")
	h.Screen.RestoreCursor()
	h.afterEvent("restore_cursor")
}

func (h *HistoryScreen) AlignmentDisplay() {
	h.beforeEvent("alignment_display")
	h.Screen.AlignmentDisplay()
	h.afterEvent("alignment_display")
}

func (h *HistoryScreen) InsertCharacters(count int) {
	h.beforeEvent("insert_characters")
	h.Screen.InsertCharacters(count)
	h.afterEvent("insert_characters")
}

func (h *HistoryScreen) CursorUp(count int) {
	h.beforeEvent("cursor_up")
	h.Screen.CursorUp(count)
	h.afterEvent("cursor_up")
}

func (h *HistoryScreen) CursorDown(count int) {
	h.beforeEvent("cursor_down")
	h.Screen.CursorDown(count)
	h.afterEvent("cursor_down")
}

func (h *HistoryScreen) CursorForward(count int) {
	h.beforeEvent("cursor_forward")
	h.Screen.CursorForward(count)
	h.afterEvent("cursor_forward")
}

func (h *HistoryScreen) CursorBack(count int) {
	h.beforeEvent("cursor_back")
	h.Screen.CursorBack(count)
	h.afterEvent("cursor_back")
}

func (h *HistoryScreen) CursorDown1(count int) {
	h.beforeEvent("cursor_down1")
	h.Screen.CursorDown1(count)
	h.afterEvent("cursor_down1")
}

func (h *HistoryScreen) CursorUp1(count int) {
	h.beforeEvent("cursor_up1")
	h.Screen.CursorUp1(count)
	h.afterEvent("cursor_up1")
}

func (h *HistoryScreen) CursorToColumn(column int) {
	h.beforeEvent("cursor_to_column")
	h.Screen.CursorToColumn(column)
	h.afterEvent("cursor_to_column")
}

func (h *HistoryScreen) CursorToColumnAbsolute(column int) {
	h.beforeEvent("cursor_to_column_absolute")
	h.Screen.CursorToColumnAbsolute(column)
	h.afterEvent("cursor_to_column_absolute")
}

func (h *HistoryScreen) CursorPosition(line, column int) {
	h.beforeEvent("cursor_position")
	h.Screen.CursorPosition(line, column)
	h.afterEvent("cursor_position")
}

func (h *HistoryScreen) EraseInDisplay(how int, private bool, rest ...int) {
	h.beforeEvent("erase_in_display")
	h.Screen.EraseInDisplay(how, private, rest...)
	if how == 3 {
		h.resetHistory()
	}
	h.afterEvent("erase_in_display")
}

func (h *HistoryScreen) EraseInLine(how int, private bool) {
	h.beforeEvent("erase_in_line")
	h.Screen.EraseInLine(how, private)
	h.afterEvent("erase_in_line")
}

func (h *HistoryScreen) InsertLines(count int) {
	h.beforeEvent("insert_lines")
	h.Screen.InsertLines(count)
	h.afterEvent("insert_lines")
}

func (h *HistoryScreen) DeleteLines(count int) {
	h.beforeEvent("delete_lines")
	h.Screen.DeleteLines(count)
	h.afterEvent("delete_lines")
}

func (h *HistoryScreen) DeleteCharacters(count int) {
	h.beforeEvent("delete_characters")
	h.Screen.DeleteCharacters(count)
	h.afterEvent("delete_characters")
}

func (h *HistoryScreen) EraseCharacters(count int) {
	h.beforeEvent("erase_characters")
	h.Screen.EraseCharacters(count)
	h.afterEvent("erase_characters")
}

func (h *HistoryScreen) ReportDeviceAttributes(mode int, private bool, prefix rune) {
	h.beforeEvent("report_device_attributes")
	h.Screen.ReportDeviceAttributes(mode, private, prefix)
	h.afterEvent("report_device_attributes")
}

func (h *HistoryScreen) CursorToLine(line int) {
	h.beforeEvent("cursor_to_line")
	h.Screen.CursorToLine(line)
	h.afterEvent("cursor_to_line")
}

func (h *HistoryScreen) CursorBackTab(count int) {
	h.beforeEvent("cursor_back_tab")
	h.Screen.CursorBackTab(count)
	h.afterEvent("cursor_back_tab")
}

func (h *HistoryScreen) CursorForwardTab(count int) {
	h.beforeEvent("cursor_forward_tab")
	h.Screen.CursorForwardTab(count)
	h.afterEvent("cursor_forward_tab")
}

func (h *HistoryScreen) ScrollUp(count int) {
	h.beforeEvent("scroll_up")
	h.Screen.ScrollUp(count)
	h.afterEvent("scroll_up")
}

func (h *HistoryScreen) ScrollDown(count int) {
	h.beforeEvent("scroll_down")
	h.Screen.ScrollDown(count)
	h.afterEvent("scroll_down")
}

func (h *HistoryScreen) RepeatLast(count int) {
	h.beforeEvent("repeat_last")
	h.Screen.RepeatLast(count)
	h.afterEvent("repeat_last")
}

func (h *HistoryScreen) ReportDeviceStatus(mode int, private bool, prefix rune) {
	h.beforeEvent("report_device_status")
	h.Screen.ReportDeviceStatus(mode, private, prefix)
	h.afterEvent("report_device_status")
}

func (h *HistoryScreen) ReportMode(mode int, private bool) {
	h.beforeEvent("report_mode")
	h.Screen.ReportMode(mode, private)
	h.afterEvent("report_mode")
}

func (h *HistoryScreen) RequestStatusString(query string) {
	h.beforeEvent("request_status_string")
	h.Screen.RequestStatusString(query)
	h.afterEvent("request_status_string")
}

func (h *HistoryScreen) SoftReset() {
	h.beforeEvent("soft_reset")
	h.Screen.SoftReset()
	h.afterEvent("soft_reset")
}

func (h *HistoryScreen) SetMargins(top, bottom int) {
	h.beforeEvent("set_margins")
	h.Screen.SetMargins(top, bottom)
	h.afterEvent("set_margins")
}

func (h *HistoryScreen) SetLeftRightMargins(left, right int) {
	h.beforeEvent("set_left_right_margins")
	h.Screen.SetLeftRightMargins(left, right)
	h.afterEvent("set_left_right_margins")
}

func (h *HistoryScreen) SelectGraphicRendition(attrs []int, private bool) {
	h.beforeEvent("select_graphic_rendition")
	h.Screen.SelectGraphicRendition(attrs, private)
	h.afterEvent("select_graphic_rendition")
}

func (h *HistoryScreen) SetMode(modes []int, private bool) {
	h.beforeEvent("set_mode")
	h.Screen.SetMode(modes, private)
	h.afterEvent("set_mode")
}

func (h *HistoryScreen) ResetMode(modes []int, private bool) {
	h.beforeEvent("reset_mode")
	h.Screen.ResetMode(modes, private)
	h.afterEvent("reset_mode")
}

func (h *HistoryScreen) DefineCharset(code, mode string) {
	h.beforeEvent("define_charset")
	h.Screen.DefineCharset(code, mode)
	h.afterEvent("define_charset")
}

func (h *HistoryScreen) SetTitle(param string) {
	h.beforeEvent("set_title")
	h.Screen.SetTitle(param)
	h.afterEvent("set_title")
}

func (h *HistoryScreen) SetIconName(param string) {
	h.beforeEvent("set_icon_name")
	h.Screen.SetIconName(param)
	h.afterEvent("set_icon_name")
}

func (h *HistoryScreen) Index() {
	h.beforeEvent("index")
	h.indexInternal()
	h.afterEvent("index")
}

func (h *HistoryScreen) ReverseIndex() {
	h.beforeEvent("reverse_index")
	h.reverseIndexInternal()
	h.afterEvent("reverse_index")
}

func (h *HistoryScreen) indexInternal() {
	top, bottom := h.scrollRegion()
	if h.Cursor.Row == bottom {
		h.history.Top.append(h.Buffer[top])
	}
	h.Screen.Index()
}

func (h *HistoryScreen) reverseIndexInternal() {
	top, bottom := h.scrollRegion()
	if h.Cursor.Row == top {
		h.history.Bottom.append(h.Buffer[bottom])
	}
	h.Screen.ReverseIndex()
}

func (h *HistoryScreen) Draw(data string) {
	h.beforeEvent("draw")
	h.Screen.Draw(data)
	h.afterEvent("draw")
}

func (h *HistoryScreen) LineFeed() {
	h.beforeEvent("linefeed")
	h.indexInternal()
	if h.isModeSet(ModeLNM) {
		h.CarriageReturn()
	}
	h.afterEvent("linefeed")
}

func (h *HistoryScreen) PrevPage() {
	if h.history.Position > h.Lines && len(h.history.Top.items) > 0 {
		mid := minInt(len(h.history.Top.items), int(math.Ceil(float64(h.Lines)*h.history.Ratio)))
		for row := h.Lines - 1; row >= h.Lines-mid; row-- {
			h.history.Bottom.prepend(h.Buffer[row])
		}
		h.history.Position -= mid
		for row := h.Lines - 1; row >= mid; row-- {
			h.Buffer[row] = h.Buffer[row-mid]
		}
		for row := mid - 1; row >= 0; row-- {
			line := h.history.Top.pop()
			if line == nil {
				line = blankLine(h.Columns, h.defaultCell())
			}
			h.Buffer[row] = line
		}
		h.markDirtyRange(0, h.Lines-1)
	}
	h.afterEvent("prev_page")
}

func (h *HistoryScreen) NextPage() {
	if h.history.Position < h.history.Size && len(h.history.Bottom.items) > 0 {
		mid := minInt(len(h.history.Bottom.items), int(math.Ceil(float64(h.Lines)*h.history.Ratio)))
		for row := 0; row < mid; row++ {
			h.history.Top.append(h.Buffer[row])
		}
		h.history.Position += mid
		for row := 0; row < h.Lines-mid; row++ {
			h.Buffer[row] = h.Buffer[row+mid]
		}
		for row := h.Lines - mid; row < h.Lines; row++ {
			line := h.history.Bottom.popLeft()
			if line == nil {
				line = blankLine(h.Columns, h.defaultCell())
			}
			h.Buffer[row] = line
		}
		h.markDirtyRange(0, h.Lines-1)
	}
	h.afterEvent("next_page")
}

func (h *HistoryScreen) History() [][]Cell {
	return append([][]Cell(nil), h.history.Top.items...)
}

func (h *HistoryScreen) Scrollback() int {
	return len(h.history.Top.items)
}

func (h *HistoryScreen) resetHistory() {
	h.history.Top.items = nil
	h.history.Bottom.items = nil
	h.history.Position = h.history.Size
}
