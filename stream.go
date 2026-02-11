package te

import (
	"errors"
	"unicode/utf8"
)

var ErrNoListener = errors.New("listener is not set")

type EventHandler interface {
	Bell()
	Backspace()
	Tab()
	LineFeed()
	CarriageReturn()
	ShiftOut()
	ShiftIn()
	Reset()
	Index()
	ReverseIndex()
	SetTabStop()
	ClearTabStop(how int)
	SaveCursor()
	RestoreCursor()
	AlignmentDisplay()
	InsertCharacters(count int)
	CursorUp(count int)
	CursorDown(count int)
	CursorForward(count int)
	CursorBack(count int)
	CursorDown1(count int)
	CursorUp1(count int)
	CursorToColumn(column int)
	CursorPosition(line, col int)
	EraseInDisplay(how int, private bool, rest ...int)
	EraseInLine(how int, private bool)
	InsertLines(count int)
	DeleteLines(count int)
	DeleteCharacters(count int)
	EraseCharacters(count int)
	ReportDeviceAttributes(mode int, private bool)
	CursorToLine(line int)
	ReportDeviceStatus(mode int)
	SetMargins(top, bottom int)
	SelectGraphicRendition(attrs []int, private bool)
	Draw(data string)
	Debug(params ...interface{})
	SetMode(modes []int, private bool)
	ResetMode(modes []int, private bool)
	DefineCharset(code, mode string)
	SetTitle(param string)
	SetIconName(param string)
}

type Stream struct {
	listener   EventHandler
	strict     bool
	useUTF8    bool
	state      parserState
	params     []int
	current    string
	private    bool
	takingText bool
	oscEsc     bool
	skipNext   bool
}

type parserState int

const (
	stateGround parserState = iota
	stateEscape
	stateCSI
	stateSharp
	stateOSC
	stateCharset
	stateEscapePercent
)

func NewStream(screen EventHandler, strict bool) *Stream {
	st := &Stream{strict: strict, useUTF8: true}
	if screen != nil {
		st.Attach(screen)
	}
	return st
}

func (st *Stream) Attach(screen EventHandler) {
	st.listener = screen
	st.state = stateGround
	st.params = nil
	st.current = ""
	st.private = false
	st.takingText = false
	st.oscEsc = false
	st.skipNext = false
}

func (st *Stream) Detach(screen EventHandler) {
	if st.listener == screen {
		st.listener = nil
	}
}

func (st *Stream) Feed(data string) (err error) {
	if st.listener == nil {
		return ErrNoListener
	}
	defer func() {
		if r := recover(); r != nil {
			st.state = stateGround
			st.resetCSI()
			st.oscEsc = false
			st.skipNext = false
			err = errors.New("handler panic")
		}
	}()

	var textBuf []rune
	flush := func() {
		if len(textBuf) > 0 {
			st.listener.Draw(string(textBuf))
			textBuf = textBuf[:0]
		}
	}

	for _, ch := range data {
		if st.skipNext {
			flush()
			if err := st.feedRune(ch); err != nil {
				return err
			}
			continue
		}
		if st.isPlainText(ch) && st.state == stateGround {
			textBuf = append(textBuf, ch)
			continue
		}
		flush()
		if err := st.feedRune(ch); err != nil {
			return err
		}
	}
	flush()
	return nil
}

func (st *Stream) feedRune(ch rune) error {
	if st.skipNext {
		st.skipNext = false
		return nil
	}
	switch st.state {
	case stateGround:
		return st.handleGround(ch)
	case stateEscape:
		return st.handleEscape(ch)
	case stateSharp:
		st.state = stateGround
		if ch == '8' {
			st.listener.AlignmentDisplay()
		}
		return nil
	case stateCharset:
		st.state = stateGround
		if st.useUTF8 {
			return nil
		}
		st.listener.DefineCharset(string(ch), st.current)
		return nil
	case stateEscapePercent:
		st.state = stateGround
		st.SelectOtherCharset(string(ch))
		return nil
	case stateCSI:
		return st.handleCSI(ch)
	case stateOSC:
		return st.handleOSC(ch)
	default:
		st.state = stateGround
	}
	return nil
}

func (st *Stream) handleGround(ch rune) error {
	switch ch {
	case '\x07':
		st.listener.Bell()
	case '\x08':
		st.listener.Backspace()
	case '\t':
		st.listener.Tab()
	case '\n', '\x0b', '\x0c':
		st.listener.LineFeed()
	case '\r':
		st.listener.CarriageReturn()
	case '\x0e':
		if !st.useUTF8 {
			st.listener.ShiftOut()
		}
	case '\x0f':
		if !st.useUTF8 {
			st.listener.ShiftIn()
		}
	case '\x1b':
		st.state = stateEscape
	case '\x9b':
		st.state = stateCSI
		st.resetCSI()
	case '\x9d':
		st.state = stateOSC
		st.current = ""
		st.oscEsc = false
	case '\x00', '\x7f':
		return nil
	default:
		st.listener.Draw(string(ch))
	}
	return nil
}

func (st *Stream) handleEscape(ch rune) error {
	st.state = stateGround
	switch ch {
	case '[':
		st.state = stateCSI
		st.resetCSI()
	case ']':
		st.state = stateOSC
		st.current = ""
		st.oscEsc = false
	case '#':
		st.state = stateSharp
	case '%':
		st.current = ""
		st.state = stateEscapePercent
	case '(', ')':
		st.current = string(ch)
		st.state = stateCharset
	case 'c':
		st.listener.Reset()
	case 'D':
		st.listener.Index()
	case 'E':
		st.listener.LineFeed()
	case 'M':
		st.listener.ReverseIndex()
	case 'H':
		st.listener.SetTabStop()
	case '7':
		st.listener.SaveCursor()
	case '8':
		st.listener.RestoreCursor()
	default:
		return nil
	}
	return nil
}

func (st *Stream) handleCSI(ch rune) error {
	if ch == '?' {
		st.private = true
		return nil
	}
	if ch >= '0' && ch <= '9' {
		st.current += string(ch)
		return nil
	}
	if ch == ';' {
		st.appendParam()
		return nil
	}
	if isControlInCSI(ch) {
		st.dispatchControl(ch)
		return nil
	}
	if ch == '$' {
		st.state = stateGround
		st.skipNext = true
		return nil
	}
	if ch == '\x18' || ch == '\x1a' {
		st.state = stateGround
		st.listener.Draw(string(ch))
		return nil
	}
	if ch == ' ' || ch == '>' {
		return nil
	}

	st.appendParam()
	params := st.params
	st.state = stateGround
	st.dispatchCSI(ch, params)
	return nil
}

func (st *Stream) handleOSC(ch rune) error {
	if st.oscEsc {
		st.oscEsc = false
		if ch == '\\' {
			st.state = stateGround
			st.finishOSC()
			return nil
		}
		st.current += string('\x1b')
		st.current += string(ch)
		return nil
	}
	if ch == '\x07' || ch == '\x9c' {
		st.state = stateGround
		st.finishOSC()
		return nil
	}
	if ch == '\x1b' {
		st.oscEsc = true
		return nil
	}
	st.current += string(ch)
	return nil
}

func (st *Stream) finishOSC() {
	if st.current == "" {
		return
	}
	code, rest := st.current[0], ""
	if len(st.current) > 1 {
		rest = st.current[1:]
	}
	if len(rest) > 0 && rest[0] == ';' {
		rest = rest[1:]
	}
	switch code {
	case '0':
		st.listener.SetIconName(rest)
		st.listener.SetTitle(rest)
	case '1':
		st.listener.SetIconName(rest)
	case '2':
		st.listener.SetTitle(rest)
	}
	st.current = ""
}

func (st *Stream) resetCSI() {
	st.params = st.params[:0]
	st.current = ""
	st.private = false
}

func (st *Stream) appendParam() {
	value := 0
	if st.current != "" {
		for _, r := range st.current {
			value = value*10 + int(r-'0')
			if value > 9999 {
				value = 9999
				break
			}
		}
	}
	st.params = append(st.params, value)
	st.current = ""
}

func (st *Stream) dispatchCSI(final rune, params []int) {
	switch final {
	case '@':
		st.listener.InsertCharacters(defaultParam(params, 0, 1))
	case 'A':
		st.listener.CursorUp(defaultParam(params, 0, 1))
	case 'B':
		st.listener.CursorDown(defaultParam(params, 0, 1))
	case 'C':
		st.listener.CursorForward(defaultParam(params, 0, 1))
	case 'D':
		st.listener.CursorBack(defaultParam(params, 0, 1))
	case 'E':
		st.listener.CursorDown1(defaultParam(params, 0, 1))
	case 'F':
		st.listener.CursorUp1(defaultParam(params, 0, 1))
	case 'G':
		st.listener.CursorToColumn(defaultParam(params, 0, 1))
	case 'H', 'f':
		st.listener.CursorPosition(defaultParam(params, 0, 0), defaultParam(params, 1, 0))
	case 'J':
		st.listener.EraseInDisplay(defaultParam(params, 0, 0), st.private)
	case 'K':
		st.listener.EraseInLine(defaultParam(params, 0, 0), st.private)
	case 'L':
		st.listener.InsertLines(defaultParam(params, 0, 1))
	case 'M':
		st.listener.DeleteLines(defaultParam(params, 0, 1))
	case 'P':
		st.listener.DeleteCharacters(defaultParam(params, 0, 1))
	case 'X':
		st.listener.EraseCharacters(defaultParam(params, 0, 1))
	case 'a':
		st.listener.CursorForward(defaultParam(params, 0, 1))
	case 'c':
		st.listener.ReportDeviceAttributes(defaultParam(params, 0, 0), st.private)
	case 'd':
		st.listener.CursorToLine(defaultParam(params, 0, 1))
	case 'e':
		st.listener.CursorDown(defaultParam(params, 0, 1))
	case 'g':
		st.listener.ClearTabStop(defaultParam(params, 0, 0))
	case 'h':
		st.listener.SetMode(params, st.private)
	case 'l':
		st.listener.ResetMode(params, st.private)
	case 'm':
		st.listener.SelectGraphicRendition(params, st.private)
	case 'n':
		st.listener.ReportDeviceStatus(defaultParam(params, 0, 0))
	case 'r':
		if len(params) == 0 {
			st.listener.SetMargins(0, 0)
			return
		}
		top := defaultParam(params, 0, 0)
		bottom := 0
		if len(params) > 1 {
			bottom = params[1]
		}
		st.listener.SetMargins(top, bottom)
	case '\'':
		st.listener.CursorToColumn(defaultParam(params, 0, 1))
	default:
		args := make([]interface{}, len(params))
		for i, v := range params {
			args[i] = v
		}
		st.listener.Debug(args...)
	}
}

func (st *Stream) SelectOtherCharset(code string) {
	if code == "@" {
		st.useUTF8 = false
	} else if code == "G" || code == "8" {
		st.useUTF8 = true
	}
}

func isControlInCSI(ch rune) bool {
	switch ch {
	case '\x07', '\x08', '\t', '\n', '\x0b', '\x0c', '\r':
		return true
	default:
		return false
	}
}

func (st *Stream) dispatchControl(ch rune) {
	switch ch {
	case '\x07':
		st.listener.Bell()
	case '\x08':
		st.listener.Backspace()
	case '\t':
		st.listener.Tab()
	case '\n', '\x0b', '\x0c':
		st.listener.LineFeed()
	case '\r':
		st.listener.CarriageReturn()
	}
}

func defaultParam(params []int, index int, fallback int) int {
	if index < len(params) {
		if params[index] == 0 {
			return fallback
		}
		return params[index]
	}
	return fallback
}

func (st *Stream) isPlainText(ch rune) bool {
	if ch == utf8.RuneError {
		return false
	}
	if ch == '\x1b' || ch == '\x9b' || ch == '\x9d' || ch == '\x00' || ch == '\x7f' {
		return false
	}
	switch ch {
	case '\x07', '\x08', '\t', '\n', '\x0b', '\x0c', '\r', '\x0e', '\x0f':
		return false
	}
	return true
}
