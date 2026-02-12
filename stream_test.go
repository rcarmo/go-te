package te

import (
	"bytes"
	"io"
	"testing"
)

type counter struct {
	count int
}

func (c *counter) Call(_ ...interface{}) {
	c.count++
}

type argCheck struct {
	count int
	args  []interface{}
}

func (a *argCheck) Call(args ...interface{}) {
	a.count++
	a.args = args
}

type mockScreen struct {
	bell                func()
	backspace           func()
	tab                 func()
	linefeed            func()
	nextLine            func()
	carriageReturn      func()
	shiftOut            func()
	shiftIn             func()
	reset               func()
	index               func()
	reverseIndex        func()
	setTabStop          func()
	clearTabStop        func(int)
	saveCursor          func()
	restoreCursor       func()
	alignmentDisplay    func()
	insertCharacters    func(int)
	cursorUp            func(int)
	cursorDown          func(int)
	cursorForward       func(int)
	cursorBack          func(int)
	cursorDown1         func(int)
	cursorUp1           func(int)
	cursorToColumn      func(int)
	cursorToColumnAbs   func(int)
	cursorPosition      func(int, int)
	cursorBackTab       func(int)
	cursorForwardTab    func(int)
	scrollUp            func(int)
	scrollDown          func(int)
	repeatLast          func(int)
	eraseInDisplay      func(int, bool, ...int)
	eraseInLine         func(int, bool)
	insertLines         func(int)
	deleteLines         func(int)
	deleteCharacters    func(int)
	eraseCharacters     func(int)
	reportDeviceAttrs   func(int, bool, rune)
	cursorToLine        func(int)
	reportDeviceStatus  func(int, bool, rune)
	reportMode          func(int, bool)
	requestStatusString func(string)
	softReset           func()
	setMargins          func(int, int)
	setLeftRightMargins func(int, int)
	selectGraphic       func([]int, bool)
	draw                func(string)
	debug               func(...interface{})
	setMode             func([]int, bool)
	resetMode           func([]int, bool)
	defineCharset       func(string, string)
	setTitle            func(string)
	setIconName         func(string)
}

func (m *mockScreen) Bell() {
	if m.bell != nil {
		m.bell()
	}
}
func (m *mockScreen) Backspace() {
	if m.backspace != nil {
		m.backspace()
	}
}
func (m *mockScreen) Tab() {
	if m.tab != nil {
		m.tab()
	}
}
func (m *mockScreen) LineFeed() {
	if m.linefeed != nil {
		m.linefeed()
	}
}
func (m *mockScreen) CarriageReturn() {
	if m.carriageReturn != nil {
		m.carriageReturn()
	}
}
func (m *mockScreen) NextLine() {
	if m.nextLine != nil {
		m.nextLine()
	}
}
func (m *mockScreen) ShiftOut() {
	if m.shiftOut != nil {
		m.shiftOut()
	}
}
func (m *mockScreen) ShiftIn() {
	if m.shiftIn != nil {
		m.shiftIn()
	}
}
func (m *mockScreen) Reset() {
	if m.reset != nil {
		m.reset()
	}
}
func (m *mockScreen) Index() {
	if m.index != nil {
		m.index()
	}
}
func (m *mockScreen) ReverseIndex() {
	if m.reverseIndex != nil {
		m.reverseIndex()
	}
}
func (m *mockScreen) SetTabStop() {
	if m.setTabStop != nil {
		m.setTabStop()
	}
}
func (m *mockScreen) ClearTabStop(how int) {
	if m.clearTabStop != nil {
		m.clearTabStop(how)
	}
}
func (m *mockScreen) SaveCursor() {
	if m.saveCursor != nil {
		m.saveCursor()
	}
}
func (m *mockScreen) RestoreCursor() {
	if m.restoreCursor != nil {
		m.restoreCursor()
	}
}
func (m *mockScreen) AlignmentDisplay() {
	if m.alignmentDisplay != nil {
		m.alignmentDisplay()
	}
}
func (m *mockScreen) InsertCharacters(count int) {
	if m.insertCharacters != nil {
		m.insertCharacters(count)
	}
}
func (m *mockScreen) CursorUp(count int) {
	if m.cursorUp != nil {
		m.cursorUp(count)
	}
}
func (m *mockScreen) CursorDown(count int) {
	if m.cursorDown != nil {
		m.cursorDown(count)
	}
}
func (m *mockScreen) CursorForward(count int) {
	if m.cursorForward != nil {
		m.cursorForward(count)
	}
}
func (m *mockScreen) CursorBack(count int) {
	if m.cursorBack != nil {
		m.cursorBack(count)
	}
}
func (m *mockScreen) CursorDown1(count int) {
	if m.cursorDown1 != nil {
		m.cursorDown1(count)
	}
}
func (m *mockScreen) CursorUp1(count int) {
	if m.cursorUp1 != nil {
		m.cursorUp1(count)
	}
}
func (m *mockScreen) CursorToColumn(column int) {
	if m.cursorToColumn != nil {
		m.cursorToColumn(column)
	}
}
func (m *mockScreen) CursorToColumnAbsolute(column int) {
	if m.cursorToColumnAbs != nil {
		m.cursorToColumnAbs(column)
	}
}
func (m *mockScreen) CursorPosition(line, col int) {
	if m.cursorPosition != nil {
		m.cursorPosition(line, col)
	}
}
func (m *mockScreen) CursorBackTab(count int) {
	if m.cursorBackTab != nil {
		m.cursorBackTab(count)
	}
}
func (m *mockScreen) CursorForwardTab(count int) {
	if m.cursorForwardTab != nil {
		m.cursorForwardTab(count)
	}
}
func (m *mockScreen) ScrollUp(count int) {
	if m.scrollUp != nil {
		m.scrollUp(count)
	}
}
func (m *mockScreen) ScrollDown(count int) {
	if m.scrollDown != nil {
		m.scrollDown(count)
	}
}
func (m *mockScreen) RepeatLast(count int) {
	if m.repeatLast != nil {
		m.repeatLast(count)
	}
}
func (m *mockScreen) EraseInDisplay(how int, private bool, rest ...int) {
	if m.eraseInDisplay != nil {
		m.eraseInDisplay(how, private, rest...)
	}
}
func (m *mockScreen) EraseInLine(how int, private bool) {
	if m.eraseInLine != nil {
		m.eraseInLine(how, private)
	}
}
func (m *mockScreen) InsertLines(count int) {
	if m.insertLines != nil {
		m.insertLines(count)
	}
}
func (m *mockScreen) DeleteLines(count int) {
	if m.deleteLines != nil {
		m.deleteLines(count)
	}
}
func (m *mockScreen) DeleteCharacters(count int) {
	if m.deleteCharacters != nil {
		m.deleteCharacters(count)
	}
}
func (m *mockScreen) EraseCharacters(count int) {
	if m.eraseCharacters != nil {
		m.eraseCharacters(count)
	}
}
func (m *mockScreen) ReportDeviceAttributes(mode int, private bool, prefix rune) {
	if m.reportDeviceAttrs != nil {
		m.reportDeviceAttrs(mode, private, prefix)
	}
}
func (m *mockScreen) CursorToLine(line int) {
	if m.cursorToLine != nil {
		m.cursorToLine(line)
	}
}
func (m *mockScreen) ReportDeviceStatus(mode int, private bool, prefix rune) {
	if m.reportDeviceStatus != nil {
		m.reportDeviceStatus(mode, private, prefix)
	}
}

func (m *mockScreen) ReportMode(mode int, private bool) {
	if m.reportMode != nil {
		m.reportMode(mode, private)
	}
}

func (m *mockScreen) RequestStatusString(query string) {
	if m.requestStatusString != nil {
		m.requestStatusString(query)
	}
}

func (m *mockScreen) SoftReset() {
	if m.softReset != nil {
		m.softReset()
	}
}
func (m *mockScreen) SetMargins(top, bottom int) {
	if m.setMargins != nil {
		m.setMargins(top, bottom)
	}
}
func (m *mockScreen) SetLeftRightMargins(left, right int) {
	if m.setLeftRightMargins != nil {
		m.setLeftRightMargins(left, right)
	}
}
func (m *mockScreen) SelectGraphicRendition(attrs []int, private bool) {
	if m.selectGraphic != nil {
		m.selectGraphic(attrs, private)
	}
}
func (m *mockScreen) Draw(data string) {
	if m.draw != nil {
		m.draw(data)
	}
}
func (m *mockScreen) Debug(params ...interface{}) {
	if m.debug != nil {
		m.debug(params...)
	}
}
func (m *mockScreen) SetMode(modes []int, private bool) {
	if m.setMode != nil {
		m.setMode(modes, private)
	}
}
func (m *mockScreen) ResetMode(modes []int, private bool) {
	if m.resetMode != nil {
		m.resetMode(modes, private)
	}
}
func (m *mockScreen) DefineCharset(code, mode string) {
	if m.defineCharset != nil {
		m.defineCharset(code, mode)
	}
}
func (m *mockScreen) SetTitle(param string) {
	if m.setTitle != nil {
		m.setTitle(param)
	}
}
func (m *mockScreen) SetIconName(param string) {
	if m.setIconName != nil {
		m.setIconName(param)
	}
}

func TestBasicSequences(t *testing.T) {
	cases := []struct {
		seq  string
		name string
	}{
		{EscRIS, "reset"},
		{EscIND, "index"},
		{EscNEL, "next_line"},
		{EscRI, "reverse_index"},
		{EscHTS, "set_tab_stop"},
		{EscDECSC, "save_cursor"},
		{EscDECRC, "restore_cursor"},
	}
	for _, tc := range cases {
		hits := counter{}
		screen := &mockScreen{}
		switch tc.name {
		case "reset":
			screen.reset = func() { hits.Call() }
		case "index":
			screen.index = func() { hits.Call() }
		case "linefeed":
			screen.linefeed = func() { hits.Call() }
		case "next_line":
			screen.nextLine = func() { hits.Call() }
		case "reverse_index":
			screen.reverseIndex = func() { hits.Call() }
		case "set_tab_stop":
			screen.setTabStop = func() { hits.Call() }
		case "save_cursor":
			screen.saveCursor = func() { hits.Call() }
		case "restore_cursor":
			screen.restoreCursor = func() { hits.Call() }
		}
		stream := NewStream(screen, false)
		if err := stream.Feed(ControlESC); err != nil {
			t.Fatalf("feed ESC: %v", err)
		}
		if hits.count != 0 {
			t.Fatalf("expected 0 hits before sequence")
		}
		if err := stream.Feed(tc.seq); err != nil {
			t.Fatalf("feed seq: %v", err)
		}
		if hits.count != 1 {
			t.Fatalf("expected 1 hit for %s, got %d", tc.name, hits.count)
		}
	}
}

func TestLinefeedControls(t *testing.T) {
	count := 0
	stream := NewStream(&mockScreen{linefeed: func() { count++ }}, false)
	if err := stream.Feed(ControlLF + ControlVT + ControlFF); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 linefeeds, got %d", count)
	}
}

func TestUnknownCSISequences(t *testing.T) {
	args := argCheck{}
	screen := &mockScreen{debug: args.Call}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "6;z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if args.count != 1 {
		t.Fatalf("expected 1 debug call, got %d", args.count)
	}
	if len(args.args) != 2 || args.args[0].(int) != 6 || args.args[1].(int) != 0 {
		t.Fatalf("unexpected args: %#v", args.args)
	}
}

func TestCSIParamParsing(t *testing.T) {
	args := argCheck{}
	screen := &mockScreen{cursorPosition: func(line, col int) { args.Call(line, col) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ";" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if args.count != 1 || args.args[0].(int) != 0 || args.args[1].(int) != 0 {
		t.Fatalf("unexpected args: %#v", args.args)
	}

	args.count = 0
	if err := stream.Feed(ControlCSI + "999999999999999;99999999999999" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if args.args[0].(int) != 9999 || args.args[1].(int) != 9999 {
		t.Fatalf("expected overflow capped")
	}
}

func TestInterruptCSI(t *testing.T) {
	seen := []string{}
	screen := &mockScreen{draw: func(data string) { seen = append(seen, data) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "10;" + ControlSUB + "10" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(seen) != 2 || seen[0] != ControlSUB || seen[1] != "10"+EscHVP {
		t.Fatalf("unexpected seen: %#v", seen)
	}
}

func TestControlCharactersInCSI(t *testing.T) {
	args := argCheck{}
	screen := &mockScreen{cursorPosition: func(line, col int) { args.Call(line, col) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "10;\t\t\n\r\n10" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if args.args[0].(int) != 10 || args.args[1].(int) != 10 {
		t.Fatalf("unexpected args: %#v", args.args)
	}
}

func TestOSCSequences(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlOSCC0 + "1;foo" + ControlSTC0); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.IconName != "foo" {
		t.Fatalf("expected icon name")
	}
	if err := stream.Feed(ControlOSCC0 + "2;bar" + ControlSTC0); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "bar" {
		t.Fatalf("expected title")
	}
	if err := stream.Feed(ControlOSCC0 + "0;baz" + ControlSTC0); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.IconName != "baz" || screen.Title != "baz" {
		t.Fatalf("expected title and icon name")
	}
}

func TestCompatibilityAPI(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(nil, false)
	stream.Attach(screen)
	stream.Attach(NewScreen(80, 24))
	if err := stream.Feed("привет"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	stream.Detach(screen)
}

func TestDefineCharsetNoop(t *testing.T) {
	screen := NewScreen(3, 3)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "(B"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "   " {
		t.Fatalf("unexpected display")
	}
}

func TestNonUTF8Shifts(t *testing.T) {
	calls := 0
	screen := &mockScreen{shiftIn: func() { calls++ }, shiftOut: func() { calls++ }}
	stream := NewStream(screen, false)
	stream.useUTF8 = false
	if err := stream.Feed(ControlSI + ControlSO); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDollarSkip(t *testing.T) {
	calls := 0
	screen := &mockScreen{draw: func(string) { calls++ }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "12$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "1;2;3;4$x"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected no draw calls")
	}
}

func TestDebugStream(t *testing.T) {
	buf := &bytes.Buffer{}
	debug := NewDebugScreen(buf)
	stream := NewByteStream(debug, false)
	if err := stream.Feed([]byte("\x1b[1;24r\x1b[4l\x1b[24;1H")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestByteStreamFeed(t *testing.T) {
	count := 0
	var seen string
	screen := &mockScreen{draw: func(data string) { count++; seen = data }}
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("Нерусский текст")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if count != 1 || seen != "Нерусский текст" {
		t.Fatalf("unexpected draw %d %q", count, seen)
	}
}

func TestByteStreamDefineCharsetUnknown(t *testing.T) {
	screen := NewScreen(3, 3)
	stream := NewByteStream(screen, false)
	stream.SelectOtherCharset("@")
	defaultCharset := screen.G0
	if err := stream.Feed([]byte(ControlESC + "(Z")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "   " {
		t.Fatalf("unexpected display")
	}
	if &screen.G0[0] != &defaultCharset[0] {
		t.Fatalf("charset should remain unchanged")
	}
}

func TestByteStreamDefineCharset(t *testing.T) {
	for charset := range charsetMaps {
		screen := NewScreen(3, 3)
		stream := NewByteStream(screen, false)
		stream.SelectOtherCharset("@")
		if err := stream.Feed([]byte(ControlESC + "(" + charset)); err != nil {
			t.Fatalf("feed: %v", err)
		}
		if screen.Display()[0] != "   " {
			t.Fatalf("unexpected display")
		}
		if &screen.G0[0] != &charsetMaps[charset][0] {
			t.Fatalf("expected charset %s", charset)
		}
	}
}

func TestByteStreamSelectOtherCharset(t *testing.T) {
	stream := NewByteStream(NewScreen(3, 3), false)
	if !stream.useUTF8 {
		t.Fatalf("expected utf8 by default")
	}
	stream.SelectOtherCharset("@")
	if stream.useUTF8 {
		t.Fatalf("expected utf8 disabled")
	}
	stream.SelectOtherCharset("X")
	if stream.useUTF8 {
		t.Fatalf("expected utf8 to remain disabled")
	}
	stream.SelectOtherCharset("G")
	if !stream.useUTF8 {
		t.Fatalf("expected utf8 enabled")
	}
}

func TestHandlerException(t *testing.T) {
	called := 0
	screen := &mockScreen{
		setMode:   func([]int, bool) { panic("boom") },
		resetMode: func([]int, bool) { called++ },
	}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2h"); err == nil {
		t.Fatalf("expected error")
	}
	if err := stream.Feed(ControlCSI + "?9;2l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected reset_mode call")
	}
}

func TestSetResetModePrivate(t *testing.T) {
	args := []interface{}{}
	screen := &mockScreen{setMode: func(modes []int, private bool) { args = append(args, modes, private) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(args) == 0 {
		t.Fatalf("expected set mode")
	}

	resetArgs := []interface{}{}
	screen = &mockScreen{resetMode: func(modes []int, private bool) { resetArgs = append(resetArgs, modes, private) }}
	stream = NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "?9;2l"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(resetArgs) == 0 {
		t.Fatalf("expected reset mode")
	}
}

func TestNonCSISequences(t *testing.T) {
	args := argCheck{}
	screen := &mockScreen{cursorUp: func(count int) { args.Call(count) }}
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "[5" + EscCUU); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if args.count != 1 || args.args[0].(int) != 5 {
		t.Fatalf("unexpected args: %#v", args.args)
	}
}

func TestSetTitleIconFromOSC(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlOSC + "1;foo" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.IconName != "foo" {
		t.Fatalf("expected icon name")
	}
	if err := stream.Feed(ControlOSC + "2;bar" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "bar" {
		t.Fatalf("expected title")
	}
}

func TestDrawAfterOSC(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlOSC + "0;bar" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("➜"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Buffer[0][0].Data != "➜" {
		t.Fatalf("expected symbol drawn")
	}
}

func TestStreamDetachNoListener(t *testing.T) {
	stream := NewStream(nil, false)
	if err := stream.Feed("foo"); err == nil {
		t.Fatalf("expected error")
	}
	stream.Detach(nil)
}

func TestDebugScreenOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	screen := NewDebugScreen(buf)
	stream := NewByteStream(screen, false)
	if err := stream.Feed([]byte("foo")); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("draw")) {
		t.Fatalf("expected draw event")
	}
}

func TestDebugScreenWritesToWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	screen := NewDebugScreen(io.Discard)
	screen.To = buf
	screen.Draw("foo")
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
}
