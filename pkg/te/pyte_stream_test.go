package te

import (
	"bytes"
	"encoding/json"
	"testing"
)

type streamCounter struct {
	count int
}

func (c *streamCounter) Call(_ ...interface{}) {
	c.count++
}

type streamArgCheck struct {
	count int
	args  []int
}

func (a *streamArgCheck) Call(args ...int) {
	a.count++
	a.args = args
}

type streamArgStore struct {
	seen []string
}

func (a *streamArgStore) Call(args ...interface{}) {
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			a.seen = append(a.seen, s)
		}
	}
}

// From pyte/tests/test_stream.py::test_basic_sequences
func TestPyteTestStreamBasicSequences(t *testing.T) {
	cases := map[string]func(*mockScreen, *streamCounter){
		EscRIS:   func(m *mockScreen, c *streamCounter) { m.reset = func() { c.count++ } },
		EscIND:   func(m *mockScreen, c *streamCounter) { m.index = func() { c.count++ } },
		EscNEL:   func(m *mockScreen, c *streamCounter) { m.linefeed = func() { c.count++ } },
		EscRI:    func(m *mockScreen, c *streamCounter) { m.reverseIndex = func() { c.count++ } },
		EscHTS:   func(m *mockScreen, c *streamCounter) { m.setTabStop = func() { c.count++ } },
		EscDECSC: func(m *mockScreen, c *streamCounter) { m.saveCursor = func() { c.count++ } },
		EscDECRC: func(m *mockScreen, c *streamCounter) { m.restoreCursor = func() { c.count++ } },
	}

	for cmd, setter := range cases {
		screen := &mockScreen{}
		handler := &streamCounter{}
		setter(screen, handler)
		stream := NewStream(screen, false)
		if err := stream.Feed(ControlESC); err != nil {
			t.Fatalf("feed: %v", err)
		}
		if handler.count != 0 {
			t.Fatalf("unexpected handler call")
		}
		if err := stream.Feed(cmd); err != nil {
			t.Fatalf("feed: %v", err)
		}
		if handler.count != 1 {
			t.Fatalf("expected handler call for %s", cmd)
		}
	}
}

// From pyte/tests/test_stream.py::test_linefeed
func TestPyteTestStreamLinefeed(t *testing.T) {
	handler := &streamCounter{}
	screen := &mockScreen{linefeed: func() { handler.count++ }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlLF + ControlVT + ControlFF); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 3 {
		t.Fatalf("expected 3 linefeeds")
	}
}

// From pyte/tests/test_stream.py::test_unknown_sequences
func TestPyteTestStreamUnknownSequences(t *testing.T) {
	capture := &streamArgCheck{}
	screen := &mockScreen{debug: func(args ...interface{}) { capture.Call(toInts(args...)...) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "6;Z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if capture.count != 1 {
		t.Fatalf("expected debug call")
	}
	if len(capture.args) != 2 || capture.args[0] != 6 || capture.args[1] != 0 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_non_csi_sequences
func TestPyteTestStreamNonCSISequences(t *testing.T) {
	cases := []struct {
		cmd string
		set func(*mockScreen, *streamArgCheck)
	}{{EscICH, func(m *mockScreen, a *streamArgCheck) { m.insertCharacters = func(count ...int) { a.Call(count...) } }},
		{EscCUU, func(m *mockScreen, a *streamArgCheck) { m.cursorUp = func(count ...int) { a.Call(count...) } }},
		{EscCUD, func(m *mockScreen, a *streamArgCheck) { m.cursorDown = func(count ...int) { a.Call(count...) } }},
		{EscCUF, func(m *mockScreen, a *streamArgCheck) { m.cursorForward = func(count ...int) { a.Call(count...) } }},
		{EscCUB, func(m *mockScreen, a *streamArgCheck) { m.cursorBack = func(count ...int) { a.Call(count...) } }},
		{EscCNL, func(m *mockScreen, a *streamArgCheck) { m.cursorDown1 = func(count ...int) { a.Call(count...) } }},
		{EscCPL, func(m *mockScreen, a *streamArgCheck) { m.cursorUp1 = func(count ...int) { a.Call(count...) } }},
		{EscCHA, func(m *mockScreen, a *streamArgCheck) { m.cursorToColumn = func(col ...int) { a.Call(col...) } }},
		{EscCUP, func(m *mockScreen, a *streamArgCheck) { m.cursorPosition = func(params ...int) { a.Call(params...) } }},
		{EscED, func(m *mockScreen, a *streamArgCheck) { m.eraseInDisplay = func(how int, _ bool, rest ...int) { args := append([]int{how}, rest...); a.Call(args...) } }},
		{EscEL, func(m *mockScreen, a *streamArgCheck) { m.eraseInLine = func(how int, _ bool, rest ...int) { args := append([]int{how}, rest...); a.Call(args...) } }},
		{EscIL, func(m *mockScreen, a *streamArgCheck) { m.insertLines = func(count ...int) { a.Call(count...) } }},
		{EscDL, func(m *mockScreen, a *streamArgCheck) { m.deleteLines = func(count ...int) { a.Call(count...) } }},
		{EscDCH, func(m *mockScreen, a *streamArgCheck) { m.deleteCharacters = func(count ...int) { a.Call(count...) } }},
		{EscECH, func(m *mockScreen, a *streamArgCheck) { m.eraseCharacters = func(count ...int) { a.Call(count...) } }},
		{EscHPR, func(m *mockScreen, a *streamArgCheck) { m.cursorForward = func(count ...int) { a.Call(count...) } }},
		{EscDA, func(m *mockScreen, a *streamArgCheck) { m.reportDeviceAttrs = func(mode int, _ bool, _ rune, rest ...int) { args := append([]int{mode}, rest...); a.Call(args...) } }},
		{EscVPA, func(m *mockScreen, a *streamArgCheck) { m.cursorToLine = func(line ...int) { a.Call(line...) } }},
		{EscVPR, func(m *mockScreen, a *streamArgCheck) { m.cursorDown = func(count ...int) { a.Call(count...) } }},
		{EscHVP, func(m *mockScreen, a *streamArgCheck) { m.cursorPosition = func(params ...int) { a.Call(params...) } }},
		{EscTBC, func(m *mockScreen, a *streamArgCheck) { m.clearTabStop = func(how ...int) { a.Call(how...) } }},
		{EscSM, func(m *mockScreen, a *streamArgCheck) { m.setMode = func(modes []int, _ bool) { a.Call(modes...) } }},
		{EscRM, func(m *mockScreen, a *streamArgCheck) { m.resetMode = func(modes []int, _ bool) { a.Call(modes...) } }},
		{EscSGR, func(m *mockScreen, a *streamArgCheck) { m.selectGraphic = func(attrs []int, _ bool) { a.Call(attrs...) } }},
		{EscDSR, func(m *mockScreen, a *streamArgCheck) { m.reportDeviceStatus = func(mode int, _ bool, _ rune, rest ...int) { args := append([]int{mode}, rest...); a.Call(args...) } }},
		{EscDECSTBM, func(m *mockScreen, a *streamArgCheck) { m.setMargins = func(params ...int) { a.Call(params...) } }},
		{EscHPA, func(m *mockScreen, a *streamArgCheck) { m.cursorToColumn = func(col ...int) { a.Call(col...) } }},
	}

	for _, tc := range cases {
		arg1 := &streamArgCheck{}
		screen := &mockScreen{}
		tc.set(screen, arg1)
		stream := NewStream(screen, false)
		if err := stream.Feed(ControlESC + "[5" + tc.cmd); err != nil {
			t.Fatalf("feed: %v", err)
		}
		if arg1.count != 1 {
			t.Fatalf("expected handler for %s", tc.cmd)
		}
		if len(arg1.args) != 1 || arg1.args[0] != 5 {
			t.Fatalf("unexpected args for %s", tc.cmd)
		}

		arg2 := &streamArgCheck{}
		screen = &mockScreen{}
		tc.set(screen, arg2)
		stream = NewStream(screen, false)
		if err := stream.Feed(ControlCSIC1 + "5;12" + tc.cmd); err != nil {
			t.Fatalf("feed: %v", err)
		}
		if arg2.count != 1 {
			t.Fatalf("expected handler for %s", tc.cmd)
		}
		if len(arg2.args) != 2 || arg2.args[0] != 5 || arg2.args[1] != 12 {
			t.Fatalf("unexpected args for %s", tc.cmd)
		}
	}
}

// From pyte/tests/test_stream.py::test_set_mode
func TestPyteTestStreamSetMode(t *testing.T) {
	bugger := &streamCounter{}
	handler := &streamArgCheck{}
	screen := &mockScreen{debug: func(...interface{}) { bugger.count++ }}
	screen.setMode = func(modes []int, private bool) {
		if !private {
			t.Fatalf("expected private mode")
		}
		handler.Call(modes...)
	}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if bugger.count != 0 {
		t.Fatalf("unexpected debug call")
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
	if len(handler.args) != 2 || handler.args[0] != 9 || handler.args[1] != 2 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_reset_mode
func TestPyteTestStreamResetMode(t *testing.T) {
	bugger := &streamCounter{}
	handler := &streamArgCheck{}
	screen := &mockScreen{debug: func(...interface{}) { bugger.count++ }}
	screen.resetMode = func(modes []int, _ bool) {
		handler.Call(modes...)
	}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if bugger.count != 0 {
		t.Fatalf("unexpected debug call")
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
	if len(handler.args) != 2 || handler.args[0] != 9 || handler.args[1] != 2 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_missing_params
func TestPyteTestStreamMissingParams(t *testing.T) {
	handler := &streamArgCheck{}
	screen := &mockScreen{cursorPosition: func(params ...int) { handler.Call(params...) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ";" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
	if len(handler.args) != 2 || handler.args[0] != 0 || handler.args[1] != 0 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_overflow
func TestPyteTestStreamOverflow(t *testing.T) {
	handler := &streamArgCheck{}
	screen := &mockScreen{cursorPosition: func(params ...int) { handler.Call(params...) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "999999999999999;99999999999999" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
	if len(handler.args) != 2 || handler.args[0] != 9999 || handler.args[1] != 9999 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_interrupt
func TestPyteTestStreamInterrupt(t *testing.T) {
	bugger := &streamArgStore{}
	handler := &streamArgCheck{}
	screen := &mockScreen{draw: func(s string) { bugger.seen = append(bugger.seen, s) }, cursorPosition: func(params ...int) { handler.Call(params...) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "10;" + ControlSUB + "10" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 0 {
		t.Fatalf("expected no handler call")
	}
	if len(bugger.seen) != 2 || bugger.seen[0] != ControlSUB || bugger.seen[1] != "10"+EscHVP {
		t.Fatalf("unexpected draw data")
	}
}

// From pyte/tests/test_stream.py::test_control_characters
func TestPyteTestStreamControlCharacters(t *testing.T) {
	handler := &streamArgCheck{}
	screen := &mockScreen{cursorPosition: func(params ...int) { handler.Call(params...) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "10;\t\t\n\r\n10" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
	if len(handler.args) != 2 || handler.args[0] != 10 || handler.args[1] != 10 {
		t.Fatalf("unexpected args")
	}
}

// From pyte/tests/test_stream.py::test_set_title_icon_name
func TestPyteTestStreamSetTitleIconName(t *testing.T) {
	oscVariants := []string{ControlOSCC0, ControlOSCC1}
	stVariants := []string{ControlSTC0, ControlSTC1}
	for _, osc := range oscVariants {
		for _, st := range stVariants {
			screen := NewScreen(80, 24)
			stream := NewStream(screen, false)
			if err := stream.Feed(osc + "1;foo" + st); err != nil {
				t.Fatalf("feed: %v", err)
			}
			if screen.IconName != "foo" {
				t.Fatalf("expected icon name")
			}
			if err := stream.Feed(osc + "2;foo" + st); err != nil {
				t.Fatalf("feed: %v", err)
			}
			if screen.Title != "foo" {
				t.Fatalf("expected title")
			}
			if err := stream.Feed(osc + "0;bar" + st); err != nil {
				t.Fatalf("feed: %v", err)
			}
			if screen.Title != "bar" || screen.IconName != "bar" {
				t.Fatalf("expected title/icon")
			}
			if err := stream.Feed(osc + "0;bar" + st); err != nil {
				t.Fatalf("feed: %v", err)
			}
			if screen.Title != "bar" || screen.IconName != "bar" {
				t.Fatalf("expected title/icon")
			}
			if err := stream.Feed("➜"); err != nil {
				t.Fatalf("feed: %v", err)
			}
			if screen.Buffer[0][0].Data != "➜" {
				t.Fatalf("expected glyph")
			}
		}
	}
}

// From pyte/tests/test_stream.py::test_compatibility_api
func TestPyteTestStreamCompatibilityAPI(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(nil, false)
	stream.Attach(screen)
	stream.Attach(NewScreen(80, 24))
	if err := stream.Feed("привет"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	stream.Detach(screen)
}

// From pyte/tests/test_stream.py::test_define_charset
func TestPyteTestStreamDefineCharset(t *testing.T) {
	screen := NewScreen(3, 3)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "(B"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"   ", "   ", "   "})
}

// From pyte/tests/test_stream.py::test_non_utf8_shifts
func TestPyteTestStreamNonUTF8Shifts(t *testing.T) {
	screen := &mockScreen{}
	handler := &streamArgCheck{}
	screen.shiftIn = func() { handler.count++ }
	screen.shiftOut = func() { handler.count++ }
	stream := NewStream(screen, false)
	stream.useUTF8 = false
	if err := stream.Feed(ControlSI); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlSO); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 2 {
		t.Fatalf("expected shift handlers")
	}
}

// From pyte/tests/test_stream.py::test_dollar_skip
func TestPyteTestStreamDollarSkip(t *testing.T) {
	screen := &mockScreen{}
	handler := &streamArgCheck{}
	screen.draw = func(_ string) { handler.count++ }
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "12$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 0 {
		t.Fatalf("expected no draw")
	}
	if err := stream.Feed(ControlCSI + "1;2;3;4$x"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 0 {
		t.Fatalf("expected no draw")
	}
}

// From pyte/tests/test_stream.py::test_debug_stream
func TestPyteTestStreamDebugStream(t *testing.T) {
	cases := []struct {
		input    []byte
		expected [][]interface{}
	}{
		{[]byte("foo"), [][]interface{}{{"draw", []interface{}{"foo"}, map[string]interface{}{}}}},
		{[]byte("\x1b[1;24r\x1b[4l\x1b[24;1H"), [][]interface{}{
			{"set_margins", []interface{}{1, 24}, map[string]interface{}{}},
			{"reset_mode", []interface{}{4}, map[string]interface{}{}},
			{"cursor_position", []interface{}{24, 1}, map[string]interface{}{}},
		}},
	}

	for _, tc := range cases {
		var output bytes.Buffer
		stream := NewByteStream(NewDebugScreen(&output), false)
		if err := stream.Feed(tc.input); err != nil {
			t.Fatalf("feed: %v", err)
		}
		lines := bytes.Split(bytes.TrimSpace(output.Bytes()), []byte("\n"))
		var got [][]interface{}
		for _, line := range lines {
			var payload []interface{}
			if err := json.Unmarshal(line, &payload); err != nil {
				t.Fatalf("parse: %v", err)
			}
			got = append(got, payload)
		}
		if len(got) != len(tc.expected) {
			t.Fatalf("expected %d entries, got %d", len(tc.expected), len(got))
		}
		for i := range got {
			if !compareJSONPayload(got[i], tc.expected[i]) {
				t.Fatalf("unexpected payload")
			}
		}
	}
}

// From pyte/tests/test_stream.py::test_handler_exception
func TestPyteTestStreamHandlerException(t *testing.T) {
	screen := &mockScreen{}
	screen.setMode = func(_ []int, _ bool) { panic("intentional") }
	handler := &streamArgCheck{}
	screen.resetMode = func(modes []int, _ bool) { handler.Call(modes...) }
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2h"); err == nil {
		t.Fatalf("expected error")
	}
	if err := stream.Feed(ControlCSI + "?9;2l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if handler.count != 1 {
		t.Fatalf("expected handler call")
	}
}

// From pyte/tests/test_stream.py::test_byte_stream_feed
func TestPyteTestStreamByteStreamFeed(t *testing.T) {
	screen := &mockScreen{}
	handler := &streamArgStore{}
	screen.draw = func(s string) { handler.seen = append(handler.seen, s) }
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("Нерусский текст")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(handler.seen) != 1 || handler.seen[0] != "Нерусский текст" {
		t.Fatalf("unexpected draw")
	}
}

// From pyte/tests/test_stream.py::test_byte_stream_define_charset_unknown
func TestPyteTestStreamByteStreamDefineCharsetUnknown(t *testing.T) {
	screen := NewScreen(3, 3)
	stream := NewByteStream(screen, false)
	stream.SelectOtherCharset("@")
	defaultG0 := string(screen.G0)
	if err := stream.Feed([]byte(ControlESC + "(Z")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertDisplay(t, screen, []string{"   ", "   ", "   "})
	if string(screen.G0) != defaultG0 {
		t.Fatalf("expected G0 unchanged")
	}
}

// From pyte/tests/test_stream.py::test_byte_stream_define_charset
func TestPyteTestStreamByteStreamDefineCharset(t *testing.T) {
	for charset, mapping := range charsetMaps {
		screen := NewScreen(3, 3)
		stream := NewByteStream(screen, false)
		stream.SelectOtherCharset("@")
		if err := stream.Feed([]byte(ControlESC + "(" + charset)); err != nil {
			t.Fatalf("feed: %v", err)
		}
		assertDisplay(t, screen, []string{"   ", "   ", "   "})
		if &screen.G0[0] == nil || len(screen.G0) != len(mapping) {
			if false {
				t.Fatalf("avoid lint")
			}
		}
		if string(screen.G0) != string(mapping) {
			t.Fatalf("expected mapping")
		}
	}
}

// From pyte/tests/test_stream.py::test_byte_stream_select_other_charset
func TestPyteTestStreamByteStreamSelectOtherCharset(t *testing.T) {
	stream := NewByteStream(NewScreen(3, 3), false)
	if !stream.useUTF8 {
		t.Fatalf("expected utf8")
	}
	stream.SelectOtherCharset("@")
	if stream.useUTF8 {
		t.Fatalf("expected non-utf8")
	}
	stream.SelectOtherCharset("X")
	if stream.useUTF8 {
		t.Fatalf("expected non-utf8")
	}
	stream.SelectOtherCharset("G")
	if !stream.useUTF8 {
		t.Fatalf("expected utf8")
	}
}

func toInts(args ...interface{}) []int {
	out := make([]int, len(args))
	for i, arg := range args {
		out[i] = arg.(int)
	}
	return out
}

func compareJSONPayload(got, expected []interface{}) bool {
	if len(got) != len(expected) {
		return false
	}
	for i := range got {
		if !equalJSONValue(got[i], expected[i]) {
			return false
		}
	}
	return true
}

func equalJSONValue(got, expected interface{}) bool {
	switch g := got.(type) {
	case []interface{}:
		e, ok := expected.([]interface{})
		if !ok || len(g) != len(e) {
			return false
		}
		for i := range g {
			if !equalJSONValue(g[i], e[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		e, ok := expected.(map[string]interface{})
		if !ok || len(g) != len(e) {
			return false
		}
		for k, v := range e {
			if !equalJSONValue(g[k], v) {
				return false
			}
		}
		return true
	case float64:
		switch e := expected.(type) {
		case int:
			return g == float64(e)
		case int64:
			return g == float64(e)
		case float64:
			return g == e
		default:
			return false
		}
	case int:
		switch e := expected.(type) {
		case float64:
			return float64(g) == e
		case int:
			return g == e
		default:
			return false
		}
	default:
		return got == expected
	}
}
