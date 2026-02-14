package te

import (
	"encoding/json"
	"io"
)

// DebugScreen records screen events as JSON payloads.
type DebugScreen struct {
	To   io.Writer
	Only map[string]struct{}
}

// NewDebugScreen creates a DebugScreen that writes to the provided writer.
func NewDebugScreen(to io.Writer) *DebugScreen {
	return &DebugScreen{To: to}
}

func (d *DebugScreen) record(name string, args []interface{}, kwargs map[string]interface{}) {
	if len(d.Only) > 0 {
		if _, ok := d.Only[name]; !ok {
			return
		}
	}
	payload := []interface{}{name, args, kwargs}
	data, _ := json.Marshal(payload)
	d.To.Write(data)
	d.To.Write([]byte("\n"))
}

func intsToInterfaces(values []int) []interface{} {
	if len(values) == 0 {
		return nil
	}
	out := make([]interface{}, len(values))
	for i, value := range values {
		out[i] = value
	}
	return out
}

// Bell triggers the terminal bell.
func (d *DebugScreen) Bell() { d.record("bell", nil, map[string]interface{}{}) }

// Backspace moves the cursor left.
func (d *DebugScreen) Backspace() { d.record("backspace", nil, map[string]interface{}{}) }

// Tab moves the cursor to the next tab stop.
func (d *DebugScreen) Tab() { d.record("tab", nil, map[string]interface{}{}) }

// LineFeed moves the cursor down and optionally scrolls.
func (d *DebugScreen) LineFeed() { d.record("linefeed", nil, map[string]interface{}{}) }

// NextLine performs a line feed followed by carriage return.
func (d *DebugScreen) NextLine() { d.record("next_line", nil, map[string]interface{}{}) }

// CarriageReturn moves the cursor to column zero.
func (d *DebugScreen) CarriageReturn() { d.record("carriage_return", nil, map[string]interface{}{}) }

// ShiftOut selects the G1 character set.
func (d *DebugScreen) ShiftOut() { d.record("shift_out", nil, map[string]interface{}{}) }

// ShiftIn selects the G0 character set.
func (d *DebugScreen) ShiftIn() { d.record("shift_in", nil, map[string]interface{}{}) }

// Reset resets state to defaults.
func (d *DebugScreen) Reset() { d.record("reset", nil, map[string]interface{}{}) }

// Index performs an index (scroll up) operation.
func (d *DebugScreen) Index() { d.record("index", nil, map[string]interface{}{}) }

// ReverseIndex performs a reverse index operation.
func (d *DebugScreen) ReverseIndex() { d.record("reverse_index", nil, map[string]interface{}{}) }

// SetTabStop adds a tab stop at the cursor column.
func (d *DebugScreen) SetTabStop() { d.record("set_tab_stop", nil, map[string]interface{}{}) }

// StartProtectedArea enables ISO protected mode.
func (d *DebugScreen) StartProtectedArea() {
	d.record("start_protected_area", nil, map[string]interface{}{})
}

// EndProtectedArea disables ISO protected mode.
func (d *DebugScreen) EndProtectedArea() {
	d.record("end_protected_area", nil, map[string]interface{}{})
}

// SetCharacterProtection updates character protection mode.
func (d *DebugScreen) SetCharacterProtection(mode int) {
	d.record("set_character_protection", []interface{}{mode}, map[string]interface{}{})
}

// ClearTabStop clears tab stops based on the given mode.
func (d *DebugScreen) ClearTabStop(how ...int) {
	d.record("clear_tab_stop", intsToInterfaces(how), map[string]interface{}{})
}

// SaveCursor saves cursor.
func (d *DebugScreen) SaveCursor() { d.record("save_cursor", nil, map[string]interface{}{}) }

// RestoreCursor restores cursor.
func (d *DebugScreen) RestoreCursor() { d.record("restore_cursor", nil, map[string]interface{}{}) }

// SaveCursorDEC saves cursor dec.
func (d *DebugScreen) SaveCursorDEC() { d.record("save_cursor_dec", nil, map[string]interface{}{}) }

// RestoreCursorDEC restores cursor dec.
func (d *DebugScreen) RestoreCursorDEC() {
	d.record("restore_cursor_dec", nil, map[string]interface{}{})
}

// AlignmentDisplay fills the screen with the alignment pattern.
func (d *DebugScreen) AlignmentDisplay() {
	d.record("alignment_display", nil, map[string]interface{}{})
}

// InsertCharacters inserts blank characters.
func (d *DebugScreen) InsertCharacters(count ...int) {
	d.record("insert_characters", intsToInterfaces(count), map[string]interface{}{})
}

// CursorUp moves the cursor up by a count.
func (d *DebugScreen) CursorUp(count ...int) {
	d.record("cursor_up", intsToInterfaces(count), map[string]interface{}{})
}

// CursorDown moves the cursor down by a count.
func (d *DebugScreen) CursorDown(count ...int) {
	d.record("cursor_down", intsToInterfaces(count), map[string]interface{}{})
}

// CursorForward moves the cursor forward by a count.
func (d *DebugScreen) CursorForward(count ...int) {
	d.record("cursor_forward", intsToInterfaces(count), map[string]interface{}{})
}

// CursorBack moves the cursor backward by a count.
func (d *DebugScreen) CursorBack(count ...int) {
	d.record("cursor_back", intsToInterfaces(count), map[string]interface{}{})
}

// CursorDown1 moves the cursor down and to column one.
func (d *DebugScreen) CursorDown1(count ...int) {
	d.record("cursor_down1", intsToInterfaces(count), map[string]interface{}{})
}

// CursorUp1 moves the cursor up and to column one.
func (d *DebugScreen) CursorUp1(count ...int) {
	d.record("cursor_up1", intsToInterfaces(count), map[string]interface{}{})
}

// CursorToColumn moves the cursor to a column.
func (d *DebugScreen) CursorToColumn(column ...int) {
	d.record("cursor_to_column", intsToInterfaces(column), map[string]interface{}{})
}

// CursorToColumnAbsolute moves the cursor to an absolute column.
func (d *DebugScreen) CursorToColumnAbsolute(column ...int) {
	d.record("cursor_to_column_absolute", intsToInterfaces(column), map[string]interface{}{})
}

// CursorPosition moves the cursor to a line and column.
func (d *DebugScreen) CursorPosition(params ...int) {
	d.record("cursor_position", intsToInterfaces(params), map[string]interface{}{})
}

// CursorBackTab moves the cursor to the previous tab stop.
func (d *DebugScreen) CursorBackTab(count ...int) {
	d.record("cursor_back_tab", intsToInterfaces(count), map[string]interface{}{})
}

// CursorForwardTab moves the cursor to the next tab stop.
func (d *DebugScreen) CursorForwardTab(count ...int) {
	d.record("cursor_forward_tab", intsToInterfaces(count), map[string]interface{}{})
}

// ScrollUp scrolls the screen up.
func (d *DebugScreen) ScrollUp(count ...int) {
	d.record("scroll_up", intsToInterfaces(count), map[string]interface{}{})
}

// ScrollDown scrolls the screen down.
func (d *DebugScreen) ScrollDown(count ...int) {
	d.record("scroll_down", intsToInterfaces(count), map[string]interface{}{})
}

// RepeatLast repeats the last drawn character.
func (d *DebugScreen) RepeatLast(count ...int) {
	d.record("repeat_last", intsToInterfaces(count), map[string]interface{}{})
}

// EraseInDisplay erases portions of the display.
func (d *DebugScreen) EraseInDisplay(how int, private bool, rest ...int) {
	args := append([]interface{}{how}, intsToInterfaces(rest)...)
	d.record("erase_in_display", args, map[string]interface{}{"private": private})
}

// EraseInLine erases portions of the current line.
func (d *DebugScreen) EraseInLine(how int, private bool, rest ...int) {
	args := append([]interface{}{how}, intsToInterfaces(rest)...)
	d.record("erase_in_line", args, map[string]interface{}{"private": private})
}

// InsertLines inserts blank lines at the cursor.
func (d *DebugScreen) InsertLines(count ...int) {
	d.record("insert_lines", intsToInterfaces(count), map[string]interface{}{})
}

// DeleteLines deletes lines at the cursor.
func (d *DebugScreen) DeleteLines(count ...int) {
	d.record("delete_lines", intsToInterfaces(count), map[string]interface{}{})
}

// DeleteCharacters deletes characters at the cursor.
func (d *DebugScreen) DeleteCharacters(count ...int) {
	d.record("delete_characters", intsToInterfaces(count), map[string]interface{}{})
}

// EraseCharacters erases characters at the cursor.
func (d *DebugScreen) EraseCharacters(count ...int) {
	d.record("erase_characters", intsToInterfaces(count), map[string]interface{}{})
}

// ReportDeviceAttributes emits a device attributes response.
func (d *DebugScreen) ReportDeviceAttributes(mode int, private bool, prefix rune, rest ...int) {
	args := append([]interface{}{mode}, intsToInterfaces(rest)...)
	d.record("report_device_attributes", args, map[string]interface{}{"private": private, "prefix": string(prefix)})
}

// CursorToLine moves the cursor to a line.
func (d *DebugScreen) CursorToLine(line ...int) {
	d.record("cursor_to_line", intsToInterfaces(line), map[string]interface{}{})
}

// ReportDeviceStatus emits a device status response.
func (d *DebugScreen) ReportDeviceStatus(mode int, private bool, prefix rune, rest ...int) {
	args := append([]interface{}{mode}, intsToInterfaces(rest)...)
	d.record("report_device_status", args, map[string]interface{}{"private": private, "prefix": string(prefix)})
}

// ReportMode emits a mode status report.
func (d *DebugScreen) ReportMode(mode int, private bool) {
	d.record("report_mode", []interface{}{mode}, map[string]interface{}{"private": private})
}

// RequestStatusString responds to DECRQSS queries.
func (d *DebugScreen) RequestStatusString(query string) {
	d.record("request_status_string", []interface{}{query}, map[string]interface{}{})
}

// SoftReset resets state without a full reset.
func (d *DebugScreen) SoftReset() {
	d.record("soft_reset", nil, map[string]interface{}{})
}

// SaveModes saves the current mode settings.
func (d *DebugScreen) SaveModes(modes []int) {
	d.record("save_modes", []interface{}{modes}, map[string]interface{}{})
}

// RestoreModes restores saved mode settings.
func (d *DebugScreen) RestoreModes(modes []int) {
	d.record("restore_modes", []interface{}{modes}, map[string]interface{}{})
}

// ForwardIndex scrolls forward within the scroll region.
func (d *DebugScreen) ForwardIndex() {
	d.record("forward_index", nil, map[string]interface{}{})
}

// BackIndex scrolls backward within the scroll region.
func (d *DebugScreen) BackIndex() {
	d.record("back_index", nil, map[string]interface{}{})
}

// InsertColumns inserts columns.
func (d *DebugScreen) InsertColumns(count int) {
	d.record("insert_columns", []interface{}{count}, map[string]interface{}{})
}

// DeleteColumns deletes columns.
func (d *DebugScreen) DeleteColumns(count int) {
	d.record("delete_columns", []interface{}{count}, map[string]interface{}{})
}

// EraseRectangle erases a rectangular area.
func (d *DebugScreen) EraseRectangle(top, left, bottom, right int) {
	d.record("erase_rectangle", []interface{}{top, left, bottom, right}, map[string]interface{}{})
}

// SelectiveEraseRectangle selectively erases a rectangular area.
func (d *DebugScreen) SelectiveEraseRectangle(top, left, bottom, right int) {
	d.record("selective_erase_rectangle", []interface{}{top, left, bottom, right}, map[string]interface{}{})
}

// FillRectangle fills a rectangle with the given rune.
func (d *DebugScreen) FillRectangle(ch rune, top, left, bottom, right int) {
	d.record("fill_rectangle", []interface{}{string(ch), top, left, bottom, right}, map[string]interface{}{})
}

// CopyRectangle copies a rectangle within the screen.
func (d *DebugScreen) CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft int) {
	d.record("copy_rectangle", []interface{}{srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft}, map[string]interface{}{})
}

// SetSelectionData stores OSC 52 selection data.
func (d *DebugScreen) SetSelectionData(selection, data string) {
	d.record("set_selection_data", []interface{}{selection, data}, map[string]interface{}{})
}

// QuerySelectionData replies with OSC 52 selection data.
func (d *DebugScreen) QuerySelectionData(selection string) {
	d.record("query_selection_data", []interface{}{selection}, map[string]interface{}{})
}

// SetColor updates a palette color.
func (d *DebugScreen) SetColor(index int, value string) {
	d.record("set_color", []interface{}{index, value}, map[string]interface{}{})
}

// QueryColor queries a palette color.
func (d *DebugScreen) QueryColor(index int) {
	d.record("query_color", []interface{}{index}, map[string]interface{}{})
}

// ResetColor resets palette colors.
func (d *DebugScreen) ResetColor(index int, all bool) {
	d.record("reset_color", []interface{}{index}, map[string]interface{}{"all": all})
}

// SetDynamicColor updates a dynamic color slot.
func (d *DebugScreen) SetDynamicColor(index int, value string) {
	d.record("set_dynamic_color", []interface{}{index, value}, map[string]interface{}{})
}

// QueryDynamicColor queries a dynamic color slot.
func (d *DebugScreen) QueryDynamicColor(index int) {
	d.record("query_dynamic_color", []interface{}{index}, map[string]interface{}{})
}

// SetSpecialColor updates a special color slot.
func (d *DebugScreen) SetSpecialColor(index int, value string) {
	d.record("set_special_color", []interface{}{index, value}, map[string]interface{}{})
}

// QuerySpecialColor queries a special color slot.
func (d *DebugScreen) QuerySpecialColor(index int) {
	d.record("query_special_color", []interface{}{index}, map[string]interface{}{})
}

// ResetSpecialColor resets special color slots.
func (d *DebugScreen) ResetSpecialColor(index int, all bool) {
	d.record("reset_special_color", []interface{}{index}, map[string]interface{}{"all": all})
}

// ResetDynamicColor resets dynamic color slots.
func (d *DebugScreen) ResetDynamicColor(index int, all bool) {
	d.record("reset_dynamic_color", []interface{}{index}, map[string]interface{}{"all": all})
}

// SetTitleMode updates title query/set behavior.
func (d *DebugScreen) SetTitleMode(params []int, reset bool) {
	d.record("set_title_mode", []interface{}{params}, map[string]interface{}{"reset": reset})
}

// SetConformance updates the VT conformance level.
func (d *DebugScreen) SetConformance(level int, sevenBit int) {
	d.record("set_conformance", []interface{}{level, sevenBit}, map[string]interface{}{})
}

// WindowOp handles window operations.
func (d *DebugScreen) WindowOp(params []int) {
	d.record("window_op", []interface{}{params}, map[string]interface{}{})
}

// SetMargins sets the top and bottom scrolling margins.
func (d *DebugScreen) SetMargins(params ...int) {
	d.record("set_margins", intsToInterfaces(params), map[string]interface{}{})
}

// SetLeftRightMargins sets left/right scrolling margins.
func (d *DebugScreen) SetLeftRightMargins(left, right int) {
	d.record("set_left_right_margins", []interface{}{left, right}, map[string]interface{}{})
}

// SelectGraphicRendition selects graphic rendition.
func (d *DebugScreen) SelectGraphicRendition(attrs []int, private bool) {
	args := make([]interface{}, len(attrs))
	for i, v := range attrs {
		args[i] = v
	}
	d.record("select_graphic_rendition", args, map[string]interface{}{"private": private})
}

// Draw renders text at the current cursor position.
func (d *DebugScreen) Draw(data string) {
	d.record("draw", []interface{}{data}, map[string]interface{}{})
}

// Debug handles debug control sequences.
func (d *DebugScreen) Debug(params ...interface{}) {
	d.record("debug", params, map[string]interface{}{})
}

// SetMode enables one or more terminal modes.
func (d *DebugScreen) SetMode(modes []int, private bool) {
	args := make([]interface{}, len(modes))
	for i, v := range modes {
		args[i] = v
	}
	kwargs := map[string]interface{}{}
	if private {
		kwargs["private"] = private
	}
	d.record("set_mode", args, kwargs)
}

// ResetMode disables one or more terminal modes.
func (d *DebugScreen) ResetMode(modes []int, private bool) {
	args := make([]interface{}, len(modes))
	for i, v := range modes {
		args[i] = v
	}
	kwargs := map[string]interface{}{}
	if private {
		kwargs["private"] = private
	}
	d.record("reset_mode", args, kwargs)
}

// DefineCharset executes the define charset operation.
func (d *DebugScreen) DefineCharset(code, mode string) {
	d.record("define_charset", []interface{}{code, mode}, map[string]interface{}{})
}

// SetTitle updates the window title.
func (d *DebugScreen) SetTitle(param string) {
	d.record("set_title", []interface{}{param}, map[string]interface{}{})
}

// SetIconName updates the icon name.
func (d *DebugScreen) SetIconName(param string) {
	d.record("set_icon_name", []interface{}{param}, map[string]interface{}{})
}
