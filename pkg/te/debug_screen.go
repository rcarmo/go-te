package te

import (
	"encoding/json"
	"io"
)

type DebugScreen struct {
	To   io.Writer
	Only map[string]struct{}
}

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

func (d *DebugScreen) Bell()           { d.record("bell", nil, map[string]interface{}{}) }
func (d *DebugScreen) Backspace()      { d.record("backspace", nil, map[string]interface{}{}) }
func (d *DebugScreen) Tab()            { d.record("tab", nil, map[string]interface{}{}) }
func (d *DebugScreen) LineFeed()       { d.record("linefeed", nil, map[string]interface{}{}) }
func (d *DebugScreen) NextLine()       { d.record("next_line", nil, map[string]interface{}{}) }
func (d *DebugScreen) CarriageReturn() { d.record("carriage_return", nil, map[string]interface{}{}) }
func (d *DebugScreen) ShiftOut()       { d.record("shift_out", nil, map[string]interface{}{}) }
func (d *DebugScreen) ShiftIn()        { d.record("shift_in", nil, map[string]interface{}{}) }
func (d *DebugScreen) Reset()          { d.record("reset", nil, map[string]interface{}{}) }
func (d *DebugScreen) Index()          { d.record("index", nil, map[string]interface{}{}) }
func (d *DebugScreen) ReverseIndex()   { d.record("reverse_index", nil, map[string]interface{}{}) }
func (d *DebugScreen) SetTabStop()     { d.record("set_tab_stop", nil, map[string]interface{}{}) }
func (d *DebugScreen) StartProtectedArea() {
	d.record("start_protected_area", nil, map[string]interface{}{})
}
func (d *DebugScreen) EndProtectedArea() {
	d.record("end_protected_area", nil, map[string]interface{}{})
}
func (d *DebugScreen) SetCharacterProtection(mode int) {
	d.record("set_character_protection", []interface{}{mode}, map[string]interface{}{})
}
func (d *DebugScreen) ClearTabStop(how ...int) {
	d.record("clear_tab_stop", intsToInterfaces(how), map[string]interface{}{})
}
func (d *DebugScreen) SaveCursor()    { d.record("save_cursor", nil, map[string]interface{}{}) }
func (d *DebugScreen) RestoreCursor() { d.record("restore_cursor", nil, map[string]interface{}{}) }
func (d *DebugScreen) SaveCursorDEC() { d.record("save_cursor_dec", nil, map[string]interface{}{}) }
func (d *DebugScreen) RestoreCursorDEC() {
	d.record("restore_cursor_dec", nil, map[string]interface{}{})
}
func (d *DebugScreen) AlignmentDisplay() {
	d.record("alignment_display", nil, map[string]interface{}{})
}
func (d *DebugScreen) InsertCharacters(count ...int) {
	d.record("insert_characters", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorUp(count ...int) {
	d.record("cursor_up", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorDown(count ...int) {
	d.record("cursor_down", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorForward(count ...int) {
	d.record("cursor_forward", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorBack(count ...int) {
	d.record("cursor_back", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorDown1(count ...int) {
	d.record("cursor_down1", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorUp1(count ...int) {
	d.record("cursor_up1", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorToColumn(column ...int) {
	d.record("cursor_to_column", intsToInterfaces(column), map[string]interface{}{})
}
func (d *DebugScreen) CursorToColumnAbsolute(column ...int) {
	d.record("cursor_to_column_absolute", intsToInterfaces(column), map[string]interface{}{})
}
func (d *DebugScreen) CursorPosition(params ...int) {
	d.record("cursor_position", intsToInterfaces(params), map[string]interface{}{})
}
func (d *DebugScreen) CursorBackTab(count ...int) {
	d.record("cursor_back_tab", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) CursorForwardTab(count ...int) {
	d.record("cursor_forward_tab", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) ScrollUp(count ...int) {
	d.record("scroll_up", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) ScrollDown(count ...int) {
	d.record("scroll_down", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) RepeatLast(count ...int) {
	d.record("repeat_last", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) EraseInDisplay(how int, private bool, rest ...int) {
	args := append([]interface{}{how}, intsToInterfaces(rest)...)
	d.record("erase_in_display", args, map[string]interface{}{"private": private})
}
func (d *DebugScreen) EraseInLine(how int, private bool, rest ...int) {
	args := append([]interface{}{how}, intsToInterfaces(rest)...)
	d.record("erase_in_line", args, map[string]interface{}{"private": private})
}
func (d *DebugScreen) InsertLines(count ...int) {
	d.record("insert_lines", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) DeleteLines(count ...int) {
	d.record("delete_lines", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) DeleteCharacters(count ...int) {
	d.record("delete_characters", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) EraseCharacters(count ...int) {
	d.record("erase_characters", intsToInterfaces(count), map[string]interface{}{})
}
func (d *DebugScreen) ReportDeviceAttributes(mode int, private bool, prefix rune, rest ...int) {
	args := append([]interface{}{mode}, intsToInterfaces(rest)...)
	d.record("report_device_attributes", args, map[string]interface{}{"private": private, "prefix": string(prefix)})
}
func (d *DebugScreen) CursorToLine(line ...int) {
	d.record("cursor_to_line", intsToInterfaces(line), map[string]interface{}{})
}
func (d *DebugScreen) ReportDeviceStatus(mode int, private bool, prefix rune, rest ...int) {
	args := append([]interface{}{mode}, intsToInterfaces(rest)...)
	d.record("report_device_status", args, map[string]interface{}{"private": private, "prefix": string(prefix)})
}

func (d *DebugScreen) ReportMode(mode int, private bool) {
	d.record("report_mode", []interface{}{mode}, map[string]interface{}{"private": private})
}

func (d *DebugScreen) RequestStatusString(query string) {
	d.record("request_status_string", []interface{}{query}, map[string]interface{}{})
}

func (d *DebugScreen) SoftReset() {
	d.record("soft_reset", nil, map[string]interface{}{})
}

func (d *DebugScreen) SaveModes(modes []int) {
	d.record("save_modes", []interface{}{modes}, map[string]interface{}{})
}

func (d *DebugScreen) RestoreModes(modes []int) {
	d.record("restore_modes", []interface{}{modes}, map[string]interface{}{})
}

func (d *DebugScreen) ForwardIndex() {
	d.record("forward_index", nil, map[string]interface{}{})
}

func (d *DebugScreen) BackIndex() {
	d.record("back_index", nil, map[string]interface{}{})
}

func (d *DebugScreen) InsertColumns(count int) {
	d.record("insert_columns", []interface{}{count}, map[string]interface{}{})
}

func (d *DebugScreen) DeleteColumns(count int) {
	d.record("delete_columns", []interface{}{count}, map[string]interface{}{})
}

func (d *DebugScreen) EraseRectangle(top, left, bottom, right int) {
	d.record("erase_rectangle", []interface{}{top, left, bottom, right}, map[string]interface{}{})
}

func (d *DebugScreen) SelectiveEraseRectangle(top, left, bottom, right int) {
	d.record("selective_erase_rectangle", []interface{}{top, left, bottom, right}, map[string]interface{}{})
}

func (d *DebugScreen) FillRectangle(ch rune, top, left, bottom, right int) {
	d.record("fill_rectangle", []interface{}{string(ch), top, left, bottom, right}, map[string]interface{}{})
}

func (d *DebugScreen) CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft int) {
	d.record("copy_rectangle", []interface{}{srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft}, map[string]interface{}{})
}

func (d *DebugScreen) SetSelectionData(selection, data string) {
	d.record("set_selection_data", []interface{}{selection, data}, map[string]interface{}{})
}

func (d *DebugScreen) QuerySelectionData(selection string) {
	d.record("query_selection_data", []interface{}{selection}, map[string]interface{}{})
}

func (d *DebugScreen) SetColor(index int, value string) {
	d.record("set_color", []interface{}{index, value}, map[string]interface{}{})
}

func (d *DebugScreen) QueryColor(index int) {
	d.record("query_color", []interface{}{index}, map[string]interface{}{})
}

func (d *DebugScreen) ResetColor(index int, all bool) {
	d.record("reset_color", []interface{}{index}, map[string]interface{}{"all": all})
}

func (d *DebugScreen) SetDynamicColor(index int, value string) {
	d.record("set_dynamic_color", []interface{}{index, value}, map[string]interface{}{})
}

func (d *DebugScreen) QueryDynamicColor(index int) {
	d.record("query_dynamic_color", []interface{}{index}, map[string]interface{}{})
}

func (d *DebugScreen) SetSpecialColor(index int, value string) {
	d.record("set_special_color", []interface{}{index, value}, map[string]interface{}{})
}

func (d *DebugScreen) QuerySpecialColor(index int) {
	d.record("query_special_color", []interface{}{index}, map[string]interface{}{})
}

func (d *DebugScreen) ResetSpecialColor(index int, all bool) {
	d.record("reset_special_color", []interface{}{index}, map[string]interface{}{"all": all})
}

func (d *DebugScreen) ResetDynamicColor(index int, all bool) {
	d.record("reset_dynamic_color", []interface{}{index}, map[string]interface{}{"all": all})
}

func (d *DebugScreen) SetTitleMode(params []int, reset bool) {
	d.record("set_title_mode", []interface{}{params}, map[string]interface{}{"reset": reset})
}

func (d *DebugScreen) SetConformance(level int, sevenBit int) {
	d.record("set_conformance", []interface{}{level, sevenBit}, map[string]interface{}{})
}

func (d *DebugScreen) WindowOp(params []int) {
	d.record("window_op", []interface{}{params}, map[string]interface{}{})
}
func (d *DebugScreen) SetMargins(params ...int) {
	d.record("set_margins", intsToInterfaces(params), map[string]interface{}{})
}
func (d *DebugScreen) SetLeftRightMargins(left, right int) {
	d.record("set_left_right_margins", []interface{}{left, right}, map[string]interface{}{})
}
func (d *DebugScreen) SelectGraphicRendition(attrs []int, private bool) {
	args := make([]interface{}, len(attrs))
	for i, v := range attrs {
		args[i] = v
	}
	d.record("select_graphic_rendition", args, map[string]interface{}{"private": private})
}
func (d *DebugScreen) Draw(data string) {
	d.record("draw", []interface{}{data}, map[string]interface{}{})
}
func (d *DebugScreen) Debug(params ...interface{}) {
	d.record("debug", params, map[string]interface{}{})
}
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
func (d *DebugScreen) DefineCharset(code, mode string) {
	d.record("define_charset", []interface{}{code, mode}, map[string]interface{}{})
}
func (d *DebugScreen) SetTitle(param string) {
	d.record("set_title", []interface{}{param}, map[string]interface{}{})
}
func (d *DebugScreen) SetIconName(param string) {
	d.record("set_icon_name", []interface{}{param}, map[string]interface{}{})
}
