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

func (d *DebugScreen) Bell()           { d.record("bell", nil, map[string]interface{}{}) }
func (d *DebugScreen) Backspace()      { d.record("backspace", nil, map[string]interface{}{}) }
func (d *DebugScreen) Tab()            { d.record("tab", nil, map[string]interface{}{}) }
func (d *DebugScreen) LineFeed()       { d.record("linefeed", nil, map[string]interface{}{}) }
func (d *DebugScreen) CarriageReturn() { d.record("carriage_return", nil, map[string]interface{}{}) }
func (d *DebugScreen) ShiftOut()       { d.record("shift_out", nil, map[string]interface{}{}) }
func (d *DebugScreen) ShiftIn()        { d.record("shift_in", nil, map[string]interface{}{}) }
func (d *DebugScreen) Reset()          { d.record("reset", nil, map[string]interface{}{}) }
func (d *DebugScreen) Index()          { d.record("index", nil, map[string]interface{}{}) }
func (d *DebugScreen) ReverseIndex()   { d.record("reverse_index", nil, map[string]interface{}{}) }
func (d *DebugScreen) SetTabStop()     { d.record("set_tab_stop", nil, map[string]interface{}{}) }
func (d *DebugScreen) ClearTabStop(how int) {
	d.record("clear_tab_stop", []interface{}{how}, map[string]interface{}{})
}
func (d *DebugScreen) SaveCursor()    { d.record("save_cursor", nil, map[string]interface{}{}) }
func (d *DebugScreen) RestoreCursor() { d.record("restore_cursor", nil, map[string]interface{}{}) }
func (d *DebugScreen) AlignmentDisplay() {
	d.record("alignment_display", nil, map[string]interface{}{})
}
func (d *DebugScreen) InsertCharacters(count int) {
	d.record("insert_characters", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorUp(count int) {
	d.record("cursor_up", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorDown(count int) {
	d.record("cursor_down", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorForward(count int) {
	d.record("cursor_forward", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorBack(count int) {
	d.record("cursor_back", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorDown1(count int) {
	d.record("cursor_down1", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorUp1(count int) {
	d.record("cursor_up1", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) CursorToColumn(column int) {
	d.record("cursor_to_column", []interface{}{column}, map[string]interface{}{})
}
func (d *DebugScreen) CursorPosition(line, col int) {
	d.record("cursor_position", []interface{}{line, col}, map[string]interface{}{})
}
func (d *DebugScreen) CursorBackTab(count int) {
	d.record("cursor_back_tab", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) ScrollUp(count int) {
	d.record("scroll_up", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) ScrollDown(count int) {
	d.record("scroll_down", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) RepeatLast(count int) {
	d.record("repeat_last", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) EraseInDisplay(how int, private bool, _ ...int) {
	d.record("erase_in_display", []interface{}{how}, map[string]interface{}{"private": private})
}
func (d *DebugScreen) EraseInLine(how int, private bool) {
	d.record("erase_in_line", []interface{}{how}, map[string]interface{}{"private": private})
}
func (d *DebugScreen) InsertLines(count int) {
	d.record("insert_lines", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) DeleteLines(count int) {
	d.record("delete_lines", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) DeleteCharacters(count int) {
	d.record("delete_characters", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) EraseCharacters(count int) {
	d.record("erase_characters", []interface{}{count}, map[string]interface{}{})
}
func (d *DebugScreen) ReportDeviceAttributes(mode int, private bool) {
	d.record("report_device_attributes", []interface{}{mode}, map[string]interface{}{"private": private})
}
func (d *DebugScreen) CursorToLine(line int) {
	d.record("cursor_to_line", []interface{}{line}, map[string]interface{}{})
}
func (d *DebugScreen) ReportDeviceStatus(mode int) {
	d.record("report_device_status", []interface{}{mode}, map[string]interface{}{})
}
func (d *DebugScreen) SetMargins(top, bottom int) {
	d.record("set_margins", []interface{}{top, bottom}, map[string]interface{}{})
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
	d.record("set_mode", args, map[string]interface{}{"private": private})
}
func (d *DebugScreen) ResetMode(modes []int, private bool) {
	args := make([]interface{}, len(modes))
	for i, v := range modes {
		args[i] = v
	}
	d.record("reset_mode", args, map[string]interface{}{"private": private})
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
