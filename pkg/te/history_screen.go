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

// History holds scrollback buffers and paging state.
type History struct {
	Top      historyDeque
	Bottom   historyDeque
	Ratio    float64
	Size     int
	Position int
}

// HistoryScreen wraps Screen to maintain scrollback history.
type HistoryScreen struct {
	*Screen
	history History
}

// NewHistoryScreen creates a HistoryScreen with the given scrollback size.
func NewHistoryScreen(cols, lines, history int) *HistoryScreen {
	return NewHistoryScreenWithRatio(cols, lines, history, 0.5)
}

// NewHistoryScreenWithRatio creates a HistoryScreen with a custom top/bottom ratio.
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

// Reset resets the screen and history state.
func (h *HistoryScreen) Reset() {
	h.beforeEvent("reset")
	h.Screen.Reset()
	h.resetHistory()
	h.afterEvent("reset")
}

// Bell triggers the terminal bell.
func (h *HistoryScreen) Bell() {
	h.beforeEvent("bell")
	h.Screen.Bell()
	h.afterEvent("bell")
}

// Backspace moves the cursor left.
func (h *HistoryScreen) Backspace() {
	h.beforeEvent("backspace")
	h.Screen.Backspace()
	h.afterEvent("backspace")
}

// Tab moves the cursor to the next tab stop.
func (h *HistoryScreen) Tab() {
	h.beforeEvent("tab")
	h.Screen.Tab()
	h.afterEvent("tab")
}

// CarriageReturn moves the cursor to column zero.
func (h *HistoryScreen) CarriageReturn() {
	h.beforeEvent("carriage_return")
	h.Screen.CarriageReturn()
	h.afterEvent("carriage_return")
}

// NextLine performs a line feed followed by carriage return.
func (h *HistoryScreen) NextLine() {
	h.beforeEvent("next_line")
	h.Screen.NextLine()
	h.afterEvent("next_line")
}

// ShiftOut selects the G1 character set.
func (h *HistoryScreen) ShiftOut() {
	h.beforeEvent("shift_out")
	h.Screen.ShiftOut()
	h.afterEvent("shift_out")
}

// ShiftIn selects the G0 character set.
func (h *HistoryScreen) ShiftIn() {
	h.beforeEvent("shift_in")
	h.Screen.ShiftIn()
	h.afterEvent("shift_in")
}

// SetTabStop adds a tab stop at the cursor column.
func (h *HistoryScreen) SetTabStop() {
	h.beforeEvent("set_tab_stop")
	h.Screen.SetTabStop()
	h.afterEvent("set_tab_stop")
}

// StartProtectedArea enables ISO protected mode.
func (h *HistoryScreen) StartProtectedArea() {
	h.beforeEvent("start_protected_area")
	h.Screen.StartProtectedArea()
	h.afterEvent("start_protected_area")
}

// EndProtectedArea disables ISO protected mode.
func (h *HistoryScreen) EndProtectedArea() {
	h.beforeEvent("end_protected_area")
	h.Screen.EndProtectedArea()
	h.afterEvent("end_protected_area")
}

// SetCharacterProtection updates character protection mode.
func (h *HistoryScreen) SetCharacterProtection(mode int) {
	h.beforeEvent("set_character_protection")
	h.Screen.SetCharacterProtection(mode)
	h.afterEvent("set_character_protection")
}

// ClearTabStop clears tab stops based on the given mode.
func (h *HistoryScreen) ClearTabStop(how ...int) {
	h.beforeEvent("clear_tab_stop")
	h.Screen.ClearTabStop(how...)
	h.afterEvent("clear_tab_stop")
}

// SaveCursor saves cursor.
func (h *HistoryScreen) SaveCursor() {
	h.beforeEvent("save_cursor")
	h.Screen.SaveCursor()
	h.afterEvent("save_cursor")
}

// RestoreCursor restores cursor.
func (h *HistoryScreen) RestoreCursor() {
	h.beforeEvent("restore_cursor")
	h.Screen.RestoreCursor()
	h.afterEvent("restore_cursor")
}

// AlignmentDisplay fills the screen with the alignment pattern.
func (h *HistoryScreen) AlignmentDisplay() {
	h.beforeEvent("alignment_display")
	h.Screen.AlignmentDisplay()
	h.afterEvent("alignment_display")
}

// InsertCharacters inserts blank characters.
func (h *HistoryScreen) InsertCharacters(count ...int) {
	h.beforeEvent("insert_characters")
	h.Screen.InsertCharacters(count...)
	h.afterEvent("insert_characters")
}

// CursorUp moves the cursor up by a count.
func (h *HistoryScreen) CursorUp(count ...int) {
	h.beforeEvent("cursor_up")
	h.Screen.CursorUp(count...)
	h.afterEvent("cursor_up")
}

// CursorDown moves the cursor down by a count.
func (h *HistoryScreen) CursorDown(count ...int) {
	h.beforeEvent("cursor_down")
	h.Screen.CursorDown(count...)
	h.afterEvent("cursor_down")
}

// CursorForward moves the cursor forward by a count.
func (h *HistoryScreen) CursorForward(count ...int) {
	h.beforeEvent("cursor_forward")
	h.Screen.CursorForward(count...)
	h.afterEvent("cursor_forward")
}

// CursorBack moves the cursor backward by a count.
func (h *HistoryScreen) CursorBack(count ...int) {
	h.beforeEvent("cursor_back")
	h.Screen.CursorBack(count...)
	h.afterEvent("cursor_back")
}

// CursorDown1 moves the cursor down and to column one.
func (h *HistoryScreen) CursorDown1(count ...int) {
	h.beforeEvent("cursor_down1")
	h.Screen.CursorDown1(count...)
	h.afterEvent("cursor_down1")
}

// CursorUp1 moves the cursor up and to column one.
func (h *HistoryScreen) CursorUp1(count ...int) {
	h.beforeEvent("cursor_up1")
	h.Screen.CursorUp1(count...)
	h.afterEvent("cursor_up1")
}

// CursorToColumn moves the cursor to a column.
func (h *HistoryScreen) CursorToColumn(column ...int) {
	h.beforeEvent("cursor_to_column")
	h.Screen.CursorToColumn(column...)
	h.afterEvent("cursor_to_column")
}

// CursorToColumnAbsolute moves the cursor to an absolute column.
func (h *HistoryScreen) CursorToColumnAbsolute(column ...int) {
	h.beforeEvent("cursor_to_column_absolute")
	h.Screen.CursorToColumnAbsolute(column...)
	h.afterEvent("cursor_to_column_absolute")
}

// CursorPosition moves the cursor to a line and column.
func (h *HistoryScreen) CursorPosition(params ...int) {
	h.beforeEvent("cursor_position")
	h.Screen.CursorPosition(params...)
	h.afterEvent("cursor_position")
}

// EraseInDisplay erases portions of the display.
func (h *HistoryScreen) EraseInDisplay(how int, private bool, rest ...int) {
	h.beforeEvent("erase_in_display")
	h.Screen.EraseInDisplay(how, private, rest...)
	if how == 3 {
		h.resetHistory()
	}
	h.afterEvent("erase_in_display")
}

// EraseInLine erases portions of the current line.
func (h *HistoryScreen) EraseInLine(how int, private bool, rest ...int) {
	h.beforeEvent("erase_in_line")
	h.Screen.EraseInLine(how, private, rest...)
	h.afterEvent("erase_in_line")
}

// InsertLines inserts blank lines at the cursor.
func (h *HistoryScreen) InsertLines(count ...int) {
	h.beforeEvent("insert_lines")
	h.Screen.InsertLines(count...)
	h.afterEvent("insert_lines")
}

// DeleteLines deletes lines at the cursor.
func (h *HistoryScreen) DeleteLines(count ...int) {
	h.beforeEvent("delete_lines")
	h.Screen.DeleteLines(count...)
	h.afterEvent("delete_lines")
}

// DeleteCharacters deletes characters at the cursor.
func (h *HistoryScreen) DeleteCharacters(count ...int) {
	h.beforeEvent("delete_characters")
	h.Screen.DeleteCharacters(count...)
	h.afterEvent("delete_characters")
}

// EraseCharacters erases characters at the cursor.
func (h *HistoryScreen) EraseCharacters(count ...int) {
	h.beforeEvent("erase_characters")
	h.Screen.EraseCharacters(count...)
	h.afterEvent("erase_characters")
}

// ReportDeviceAttributes emits a device attributes response.
func (h *HistoryScreen) ReportDeviceAttributes(mode int, private bool, prefix rune, rest ...int) {
	h.beforeEvent("report_device_attributes")
	h.Screen.ReportDeviceAttributes(mode, private, prefix, rest...)
	h.afterEvent("report_device_attributes")
}

// CursorToLine moves the cursor to a line.
func (h *HistoryScreen) CursorToLine(line ...int) {
	h.beforeEvent("cursor_to_line")
	h.Screen.CursorToLine(line...)
	h.afterEvent("cursor_to_line")
}

// CursorBackTab moves the cursor to the previous tab stop.
func (h *HistoryScreen) CursorBackTab(count ...int) {
	h.beforeEvent("cursor_back_tab")
	h.Screen.CursorBackTab(count...)
	h.afterEvent("cursor_back_tab")
}

// CursorForwardTab moves the cursor to the next tab stop.
func (h *HistoryScreen) CursorForwardTab(count ...int) {
	h.beforeEvent("cursor_forward_tab")
	h.Screen.CursorForwardTab(count...)
	h.afterEvent("cursor_forward_tab")
}

// ScrollUp scrolls the screen up.
func (h *HistoryScreen) ScrollUp(count ...int) {
	h.beforeEvent("scroll_up")
	h.Screen.ScrollUp(count...)
	h.afterEvent("scroll_up")
}

// ScrollDown scrolls the screen down.
func (h *HistoryScreen) ScrollDown(count ...int) {
	h.beforeEvent("scroll_down")
	h.Screen.ScrollDown(count...)
	h.afterEvent("scroll_down")
}

// RepeatLast repeats the last drawn character.
func (h *HistoryScreen) RepeatLast(count ...int) {
	h.beforeEvent("repeat_last")
	h.Screen.RepeatLast(count...)
	h.afterEvent("repeat_last")
}

// ReportDeviceStatus emits a device status response.
func (h *HistoryScreen) ReportDeviceStatus(mode int, private bool, prefix rune, rest ...int) {
	h.beforeEvent("report_device_status")
	h.Screen.ReportDeviceStatus(mode, private, prefix, rest...)
	h.afterEvent("report_device_status")
}

// ReportMode emits a mode status report.
func (h *HistoryScreen) ReportMode(mode int, private bool) {
	h.beforeEvent("report_mode")
	h.Screen.ReportMode(mode, private)
	h.afterEvent("report_mode")
}

// RequestStatusString responds to DECRQSS queries.
func (h *HistoryScreen) RequestStatusString(query string) {
	h.beforeEvent("request_status_string")
	h.Screen.RequestStatusString(query)
	h.afterEvent("request_status_string")
}

// SoftReset resets state without a full reset.
func (h *HistoryScreen) SoftReset() {
	h.beforeEvent("soft_reset")
	h.Screen.SoftReset()
	h.afterEvent("soft_reset")
}

// SaveModes saves the current mode settings.
func (h *HistoryScreen) SaveModes(modes []int) {
	h.beforeEvent("save_modes")
	h.Screen.SaveModes(modes)
	h.afterEvent("save_modes")
}

// RestoreModes restores saved mode settings.
func (h *HistoryScreen) RestoreModes(modes []int) {
	h.beforeEvent("restore_modes")
	h.Screen.RestoreModes(modes)
	h.afterEvent("restore_modes")
}

// ForwardIndex scrolls forward within the scroll region.
func (h *HistoryScreen) ForwardIndex() {
	h.beforeEvent("forward_index")
	h.Screen.ForwardIndex()
	h.afterEvent("forward_index")
}

// BackIndex scrolls backward within the scroll region.
func (h *HistoryScreen) BackIndex() {
	h.beforeEvent("back_index")
	h.Screen.BackIndex()
	h.afterEvent("back_index")
}

// InsertColumns inserts columns.
func (h *HistoryScreen) InsertColumns(count int) {
	h.beforeEvent("insert_columns")
	h.Screen.InsertColumns(count)
	h.afterEvent("insert_columns")
}

// DeleteColumns deletes columns.
func (h *HistoryScreen) DeleteColumns(count int) {
	h.beforeEvent("delete_columns")
	h.Screen.DeleteColumns(count)
	h.afterEvent("delete_columns")
}

// EraseRectangle erases a rectangular area.
func (h *HistoryScreen) EraseRectangle(top, left, bottom, right int) {
	h.beforeEvent("erase_rectangle")
	h.Screen.EraseRectangle(top, left, bottom, right)
	h.afterEvent("erase_rectangle")
}

// FillRectangle fills a rectangle with the given rune.
func (h *HistoryScreen) FillRectangle(ch rune, top, left, bottom, right int) {
	h.beforeEvent("fill_rectangle")
	h.Screen.FillRectangle(ch, top, left, bottom, right)
	h.afterEvent("fill_rectangle")
}

// CopyRectangle copies a rectangle within the screen.
func (h *HistoryScreen) CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft int) {
	h.beforeEvent("copy_rectangle")
	h.Screen.CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft)
	h.afterEvent("copy_rectangle")
}

// SetSelectionData stores OSC 52 selection data.
func (h *HistoryScreen) SetSelectionData(selection, data string) {
	h.beforeEvent("set_selection_data")
	h.Screen.SetSelectionData(selection, data)
	h.afterEvent("set_selection_data")
}

// QuerySelectionData replies with OSC 52 selection data.
func (h *HistoryScreen) QuerySelectionData(selection string) {
	h.beforeEvent("query_selection_data")
	h.Screen.QuerySelectionData(selection)
	h.afterEvent("query_selection_data")
}

// SetColor updates a palette color.
func (h *HistoryScreen) SetColor(index int, value string) {
	h.beforeEvent("set_color")
	h.Screen.SetColor(index, value)
	h.afterEvent("set_color")
}

// QueryColor queries a palette color.
func (h *HistoryScreen) QueryColor(index int) {
	h.beforeEvent("query_color")
	h.Screen.QueryColor(index)
	h.afterEvent("query_color")
}

// ResetColor resets palette colors.
func (h *HistoryScreen) ResetColor(index int, all bool) {
	h.beforeEvent("reset_color")
	h.Screen.ResetColor(index, all)
	h.afterEvent("reset_color")
}

// SetDynamicColor updates a dynamic color slot.
func (h *HistoryScreen) SetDynamicColor(index int, value string) {
	h.beforeEvent("set_dynamic_color")
	h.Screen.SetDynamicColor(index, value)
	h.afterEvent("set_dynamic_color")
}

// QueryDynamicColor queries a dynamic color slot.
func (h *HistoryScreen) QueryDynamicColor(index int) {
	h.beforeEvent("query_dynamic_color")
	h.Screen.QueryDynamicColor(index)
	h.afterEvent("query_dynamic_color")
}

// SetSpecialColor updates a special color slot.
func (h *HistoryScreen) SetSpecialColor(index int, value string) {
	h.beforeEvent("set_special_color")
	h.Screen.SetSpecialColor(index, value)
	h.afterEvent("set_special_color")
}

// QuerySpecialColor queries a special color slot.
func (h *HistoryScreen) QuerySpecialColor(index int) {
	h.beforeEvent("query_special_color")
	h.Screen.QuerySpecialColor(index)
	h.afterEvent("query_special_color")
}

// ResetSpecialColor resets special color slots.
func (h *HistoryScreen) ResetSpecialColor(index int, all bool) {
	h.beforeEvent("reset_special_color")
	h.Screen.ResetSpecialColor(index, all)
	h.afterEvent("reset_special_color")
}

// SetTitleMode updates title query/set behavior.
func (h *HistoryScreen) SetTitleMode(params []int, reset bool) {
	h.beforeEvent("set_title_mode")
	h.Screen.SetTitleMode(params, reset)
	h.afterEvent("set_title_mode")
}

// SetConformance updates the VT conformance level.
func (h *HistoryScreen) SetConformance(level int, sevenBit int) {
	h.beforeEvent("set_conformance")
	h.Screen.SetConformance(level, sevenBit)
	h.afterEvent("set_conformance")
}

// WindowOp handles window operations.
func (h *HistoryScreen) WindowOp(params []int) {
	h.beforeEvent("window_op")
	h.Screen.WindowOp(params)
	h.afterEvent("window_op")
}

// SetMargins sets the top and bottom scrolling margins.
func (h *HistoryScreen) SetMargins(params ...int) {
	h.beforeEvent("set_margins")
	h.Screen.SetMargins(params...)
	h.afterEvent("set_margins")
}

// SetLeftRightMargins sets left/right scrolling margins.
func (h *HistoryScreen) SetLeftRightMargins(left, right int) {
	h.beforeEvent("set_left_right_margins")
	h.Screen.SetLeftRightMargins(left, right)
	h.afterEvent("set_left_right_margins")
}

// SelectGraphicRendition selects graphic rendition.
func (h *HistoryScreen) SelectGraphicRendition(attrs []int, private bool) {
	h.beforeEvent("select_graphic_rendition")
	h.Screen.SelectGraphicRendition(attrs, private)
	h.afterEvent("select_graphic_rendition")
}

// SetMode enables one or more terminal modes.
func (h *HistoryScreen) SetMode(modes []int, private bool) {
	h.beforeEvent("set_mode")
	h.Screen.SetMode(modes, private)
	h.afterEvent("set_mode")
}

// ResetMode disables one or more terminal modes.
func (h *HistoryScreen) ResetMode(modes []int, private bool) {
	h.beforeEvent("reset_mode")
	h.Screen.ResetMode(modes, private)
	h.afterEvent("reset_mode")
}

// DefineCharset executes the define charset operation.
func (h *HistoryScreen) DefineCharset(code, mode string) {
	h.beforeEvent("define_charset")
	h.Screen.DefineCharset(code, mode)
	h.afterEvent("define_charset")
}

// SetTitle updates the window title.
func (h *HistoryScreen) SetTitle(param string) {
	h.beforeEvent("set_title")
	h.Screen.SetTitle(param)
	h.afterEvent("set_title")
}

// SetIconName updates the icon name.
func (h *HistoryScreen) SetIconName(param string) {
	h.beforeEvent("set_icon_name")
	h.Screen.SetIconName(param)
	h.afterEvent("set_icon_name")
}

// Index performs an index (scroll up) operation.
func (h *HistoryScreen) Index() {
	h.beforeEvent("index")
	h.indexInternal()
	h.afterEvent("index")
}

// ReverseIndex performs a reverse index operation.
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

// Draw renders text at the current cursor position.
func (h *HistoryScreen) Draw(data string) {
	h.beforeEvent("draw")
	h.Screen.Draw(data)
	h.afterEvent("draw")
}

// LineFeed moves the cursor down and optionally scrolls.
func (h *HistoryScreen) LineFeed() {
	h.beforeEvent("linefeed")
	h.indexInternal()
	if h.isModeSet(ModeLNM) {
		h.CarriageReturn()
	}
	h.afterEvent("linefeed")
}

// PrevPage moves to the previous line page.
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

// NextPage moves to the next line page.
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

// History executes the history operation.
func (h *HistoryScreen) History() [][]Cell {
	return append([][]Cell(nil), h.history.Top.items...)
}

// Scrollback scrolls the screen back.
func (h *HistoryScreen) Scrollback() int {
	return len(h.history.Top.items)
}

func (h *HistoryScreen) resetHistory() {
	h.history.Top.items = nil
	h.history.Bottom.items = nil
	h.history.Position = h.history.Size
}
