package te

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrNoListener = errors.New("listener is not set")

type EventHandler interface {
	Bell()
	Backspace()
	Tab()
	LineFeed()
	NextLine()
	CarriageReturn()
	ShiftOut()
	ShiftIn()
	Reset()
	Index()
	ReverseIndex()
	SetTabStop()
	ClearTabStop(how ...int)
	SaveCursor()
	RestoreCursor()
	AlignmentDisplay()
	InsertCharacters(count ...int)
	CursorUp(count ...int)
	CursorDown(count ...int)
	CursorForward(count ...int)
	CursorBack(count ...int)
	CursorDown1(count ...int)
	CursorUp1(count ...int)
	CursorToColumn(column ...int)
	CursorToColumnAbsolute(column ...int)
	CursorPosition(params ...int)
	CursorBackTab(count ...int)
	CursorForwardTab(count ...int)
	ScrollUp(count ...int)
	ScrollDown(count ...int)
	RepeatLast(count ...int)
	EraseInDisplay(how int, private bool, rest ...int)
	EraseInLine(how int, private bool, rest ...int)
	InsertLines(count ...int)
	DeleteLines(count ...int)
	DeleteCharacters(count ...int)
	EraseCharacters(count ...int)
	ReportDeviceAttributes(mode int, private bool, prefix rune, rest ...int)
	CursorToLine(line ...int)
	ReportDeviceStatus(mode int, private bool, prefix rune, rest ...int)
	ReportMode(mode int, private bool)
	RequestStatusString(query string)
	SoftReset()
	SaveModes(modes []int)
	RestoreModes(modes []int)
	ForwardIndex()
	BackIndex()
	InsertColumns(count int)
	DeleteColumns(count int)
	EraseRectangle(top, left, bottom, right int)
	FillRectangle(ch rune, top, left, bottom, right int)
	CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft int)
	SetMargins(params ...int)
	SetLeftRightMargins(left, right int)
	SelectGraphicRendition(attrs []int, private bool)
	Draw(data string)
	Debug(params ...interface{})
	SetMode(modes []int, private bool)
	ResetMode(modes []int, private bool)
	DefineCharset(code, mode string)
	SetTitle(param string)
	SetIconName(param string)
	SetSelectionData(selection, data string)
	QuerySelectionData(selection string)
	SetColor(index int, value string)
	QueryColor(index int)
	ResetColor(index int, all bool)
	SetDynamicColor(index int, value string)
	QueryDynamicColor(index int)
	SetSpecialColor(index int, value string)
	QuerySpecialColor(index int)
	ResetSpecialColor(index int, all bool)
	SetTitleMode(params []int, reset bool)
	SetConformance(level int, sevenBit int)
	WindowOp(params []int)
}

type Stream struct {
	listener   EventHandler
	strict          bool
	useUTF8         bool
	use8BitControls bool
	state           parserState
	params          []int
	current         string
	private         bool
	csiPrefix       rune
	csiIntermediate rune
	takingText      bool
	oscEsc          bool
	skipNext        bool
	dcsData         string
	paramStrings    []string
}

type parserState int

const (
	stateGround parserState = iota
	stateEscape
	stateCSI
	stateSharp
	stateOSC
	stateDCS
	stateAPC
	statePM
	stateSOS
	stateCharset
	stateEscapePercent
	stateEscapeSpace
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
	st.csiPrefix = 0
	st.csiIntermediate = 0
	st.dcsData = ""
	st.takingText = false
	st.oscEsc = false
	st.skipNext = false
	st.use8BitControls = false
	st.paramStrings = nil
}

func (st *Stream) Detach(screen EventHandler) {
	if st.listener == screen {
		st.listener = nil
	}
}

func (st *Stream) Feed(data string) (err error) {
	return st.FeedBytes([]byte(data))
}

func (st *Stream) FeedBytes(data []byte) (err error) {
	if st.listener == nil {
		return ErrNoListener
	}
	defer func() {
		if r := recover(); r != nil {
			st.state = stateGround
			st.resetCSI()
			st.oscEsc = false
			st.skipNext = false
			err = fmt.Errorf("handler panic: %v", r)
		}
	}()

	var textBuf []rune
	flush := func() {
		if len(textBuf) > 0 {
			st.listener.Draw(string(textBuf))
			textBuf = textBuf[:0]
		}
	}

	for i := 0; i < len(data); {
		b := data[i]
		var ch rune
		size := 1
		if b < 0x80 {
			ch = rune(b)
		} else if b >= 0xC0 {
			r, n := utf8.DecodeRune(data[i:])
			if r != utf8.RuneError || n > 1 {
				ch = r
				size = n
			} else {
				ch = rune(b)
			}
		} else {
			ch = rune(b)
		}

		if st.skipNext {
			flush()
			st.skipNext = false
			if err := st.feedRune(ch); err != nil {
				return err
			}
			i += size
			continue
		}
		if st.isPlainText(ch) && st.state == stateGround {
			textBuf = append(textBuf, ch)
			i += size
			continue
		}
		flush()
		if err := st.feedRune(ch); err != nil {
			return err
		}
		i += size
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
	case stateDCS, stateAPC, statePM, stateSOS:
		return st.handleString(ch)
	case stateEscapeSpace:
		st.state = stateGround
		switch ch {
		case 'F':
			st.use8BitControls = false
		case 'G':
			st.use8BitControls = true
		}
		return nil
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
	case '\x84':
		st.listener.Index()
	case '\x85':
		st.listener.LineFeed()
		st.listener.CarriageReturn()
	case '\x88':
		st.listener.SetTabStop()
	case '\x8d':
		st.listener.ReverseIndex()
	case '\x90':
		st.state = stateDCS
		st.dcsData = ""
	case '\x9a':
		st.listener.ReportDeviceAttributes(0, true, '?')
	case '\x98':
		st.state = stateSOS
	case '\x9b':
		st.state = stateCSI
		st.resetCSI()
	case '\x9d':
		st.state = stateOSC
		st.current = ""
		st.oscEsc = false
	case '\x9e':
		st.state = statePM
	case '\x9f':
		st.state = stateAPC
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
	case 'P':
		st.state = stateDCS
		st.oscEsc = false
		st.dcsData = ""
	case '_':
		st.state = stateAPC
		st.oscEsc = false
	case '^':
		st.state = statePM
		st.oscEsc = false
	case 'X':
		st.state = stateSOS
		st.oscEsc = false
	case '#':
		st.state = stateSharp
	case ' ':
		st.state = stateEscapeSpace
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
		st.listener.CarriageReturn()
	case 'M':
		st.listener.ReverseIndex()
	case 'H':
		st.listener.SetTabStop()
	case '7':
		st.listener.SaveCursor()
	case '8':
		st.listener.RestoreCursor()
	case '6':
		st.listener.BackIndex()
	case '9':
		st.listener.ForwardIndex()
	case 'Z':
		st.listener.ReportDeviceAttributes(0, true, '?')
	default:
		return nil
	}
	return nil
}

func (st *Stream) handleCSI(ch rune) error {
	if ch == '?' {
		st.private = true
		st.csiPrefix = ch
		return nil
	}
	if ch == '>' {
		st.csiPrefix = ch
		return nil
	}
	if ch == '\'' {
		st.appendParam()
		params := st.params
		st.state = stateGround
		st.dispatchCSI(ch, params)
		return nil
	}
	if ch >= ' ' && ch <= '/' {
		st.csiIntermediate = ch
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
		st.csiIntermediate = ch
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

func (st *Stream) handleString(ch rune) error {
	if st.state == stateDCS {
		if st.oscEsc {
			st.oscEsc = false
			if ch == '\\' {
				st.state = stateGround
				st.finishDCS()
				return nil
			}
			return nil
		}
		if ch == '\x07' || ch == '\x9c' {
			st.state = stateGround
			st.finishDCS()
			return nil
		}
		if ch == '\x1b' {
			st.oscEsc = true
			return nil
		}
		st.dcsData += string(ch)
		return nil
	}
	if st.oscEsc {
		st.oscEsc = false
		if ch == '\\' {
			st.state = stateGround
			return nil
		}
		return nil
	}
	if ch == '\x07' || ch == '\x9c' {
		st.state = stateGround
		return nil
	}
	if ch == '\x1b' {
		st.oscEsc = true
	}
	return nil
}

func (st *Stream) finishDCS() {
	data := st.dcsData
	st.dcsData = ""
	if len(data) >= 2 && data[0] == '$' && data[1] == 'q' {
		st.listener.RequestStatusString(data[2:])
	}
}

func (st *Stream) finishOSC() {
	if st.current == "" {
		return
	}
	chunks := strings.Split(st.current, ";")
	if len(chunks) == 0 {
		return
	}
	option := chunks[0]
	rest := strings.Join(chunks[1:], ";")
	switch option {
	case "0":
		st.listener.SetIconName(rest)
		st.listener.SetTitle(rest)
	case "1":
		st.listener.SetIconName(rest)
	case "2":
		st.listener.SetTitle(rest)
	case "52":
		selection := "s0"
		if len(chunks) > 1 && chunks[1] != "" {
			selection = chunks[1]
		}
		if len(chunks) > 2 {
			if chunks[2] == "?" {
				st.listener.QuerySelectionData(selection)
			} else {
				st.listener.SetSelectionData(selection, chunks[2])
			}
		}
	case "4":
		for i := 1; i+1 < len(chunks); i += 2 {
			idx, err := strconv.Atoi(chunks[i])
			if err != nil {
				continue
			}
			spec := chunks[i+1]
			if spec == "?" {
				st.listener.QueryColor(idx)
				continue
			}
			st.listener.SetColor(idx, spec)
		}
	case "104":
		if len(chunks) == 1 || chunks[1] == "" {
			st.listener.ResetColor(0, true)
			return
		}
		if idx, err := strconv.Atoi(chunks[1]); err == nil {
			st.listener.ResetColor(idx, false)
		}
	case "10", "11":
		idx, err := strconv.Atoi(option)
		if err != nil {
			break
		}
		if len(chunks) > 1 {
			spec := chunks[1]
			if spec == "?" {
				st.listener.QueryDynamicColor(idx)
				break
			}
			st.listener.SetDynamicColor(idx, spec)
		}
	case "5":
		if len(chunks) > 2 {
			idx, err := strconv.Atoi(chunks[1])
			if err != nil {
				break
			}
			spec := chunks[2]
			if spec == "?" {
				st.listener.QuerySpecialColor(idx)
				break
			}
			st.listener.SetSpecialColor(idx, spec)
		}
	case "105":
		if len(chunks) == 1 || chunks[1] == "" {
			st.listener.ResetSpecialColor(0, true)
			break
		}
		if idx, err := strconv.Atoi(chunks[1]); err == nil {
			st.listener.ResetSpecialColor(idx, false)
		}
	}
	st.current = ""
}

func (st *Stream) resetCSI() {
	st.params = st.params[:0]
	st.paramStrings = st.paramStrings[:0]
	st.current = ""
	st.private = false
	st.csiPrefix = 0
	st.csiIntermediate = 0
}

func (st *Stream) appendParam() {
	value := 0
	current := st.current
	if current != "" {
		for _, r := range current {
			value = value*10 + int(r-'0')
			if value > 9999 {
				value = 9999
				break
			}
		}
	}
	st.params = append(st.params, value)
	st.paramStrings = append(st.paramStrings, current)
	st.current = ""
}

func (st *Stream) dispatchCSI(final rune, params []int) {
	switch final {
	case '@':
		st.listener.InsertCharacters(params...)
	case 'A':
		st.listener.CursorUp(params...)
	case 'B':
		st.listener.CursorDown(params...)
	case 'C':
		st.listener.CursorForward(params...)
	case 'D':
		st.listener.CursorBack(params...)
	case 'E':
		st.listener.CursorDown1(params...)
	case 'F':
		st.listener.CursorUp1(params...)
	case 'G':
		st.listener.CursorToColumn(params...)
	case 'H', 'f':
		st.listener.CursorPosition(params...)
	case 'J':
		how := defaultParam(params, 0, 0)
		rest := []int{}
		if len(params) > 1 {
			rest = params[1:]
		}
		st.listener.EraseInDisplay(how, st.private, rest...)
	case 'K':
		how := defaultParam(params, 0, 0)
		rest := []int{}
		if len(params) > 1 {
			rest = params[1:]
		}
		st.listener.EraseInLine(how, st.private, rest...)
	case 'L':
		st.listener.InsertLines(params...)
	case 'M':
		st.listener.DeleteLines(params...)
	case 'P':
		st.listener.DeleteCharacters(params...)
	case 'X':
		st.listener.EraseCharacters(params...)
	case 'I':
		st.listener.CursorForwardTab(params...)
	case 'S':
		st.listener.ScrollUp(params...)
	case 'T':
		st.listener.ScrollDown(params...)
	case 'b':
		st.listener.RepeatLast(params...)
	case 'a':
		st.listener.CursorForward(params...)
	case '`':
		st.listener.CursorToColumnAbsolute(params...)
	case 'c':
		mode := defaultParam(params, 0, 0)
		rest := []int{}
		if len(params) > 1 {
			rest = params[1:]
		}
		st.listener.ReportDeviceAttributes(mode, st.private, st.csiPrefix, rest...)
	case 'd':
		st.listener.CursorToLine(params...)
	case 'e':
		st.listener.CursorDown(params...)
	case 'g':
		st.listener.ClearTabStop(params...)
	case 'h':
		st.listener.SetMode(params, st.private)
	case 'l':
		st.listener.ResetMode(params, st.private)
	case 'm':
		st.listener.SelectGraphicRendition(params, st.private)
	case 'n':
		mode := defaultParam(params, 0, 0)
		rest := []int{}
		if len(params) > 1 {
			rest = params[1:]
		}
		st.listener.ReportDeviceStatus(mode, st.private, st.csiPrefix, rest...)
	case 'r':
		if st.private {
			st.listener.RestoreModes(params)
			return
		}
		st.listener.SetMargins(params...)
	case 's':
		if st.private {
			st.listener.SaveModes(params)
			return
		}
		if len(params) >= 2 {
			st.listener.SetLeftRightMargins(defaultParam(params, 0, 0), defaultParam(params, 1, 0))
			return
		}
		if mc, ok := st.listener.(interface{ isModeSet(int) bool }); ok {
			if mc.isModeSet(ModeDECLRMM) {
				return
			}
		}
		st.listener.SaveCursor()
	case 'u':
		st.listener.RestoreCursor()
	case '\'':
		st.listener.CursorToColumn(params...)
	default:
		if st.csiPrefix == '>' && final == 't' {
			st.listener.SetTitleMode(params, false)
			return
		}
		if st.csiPrefix == '>' && final == 'T' {
			st.listener.SetTitleMode(params, true)
			return
		}
		if st.csiIntermediate == '"' && final == 'p' {
			level := defaultParam(params, 0, 0)
			seven := defaultParam(params, 1, 0)
			st.listener.SetConformance(level, seven)
			return
		}
		if final == 't' {
			paramsWithOmitted := make([]int, len(params))
			for i, p := range params {
				if i < len(st.paramStrings) && st.paramStrings[i] == "" {
					paramsWithOmitted[i] = -1
				} else {
					paramsWithOmitted[i] = p
				}
			}
			st.listener.WindowOp(paramsWithOmitted)
			return
		}
		if st.csiIntermediate == '$' && final == 'p' {
			st.listener.ReportMode(defaultParam(params, 0, 0), st.private)
			return
		}
		if st.csiIntermediate == '!' && final == 'p' {
			st.listener.SoftReset()
			return
		}
		if st.csiIntermediate == '\'' && final == '}' {
			st.listener.InsertColumns(defaultParam(params, 0, 1))
			return
		}
		if st.csiIntermediate == '\'' && final == '~' {
			st.listener.DeleteColumns(defaultParam(params, 0, 1))
			return
		}
		if st.csiIntermediate == '$' && final == 'z' {
			st.listener.EraseRectangle(defaultParam(params, 0, 1), defaultParam(params, 1, 1), defaultParam(params, 2, 1), defaultParam(params, 3, 1))
			return
		}
		if st.csiIntermediate == '$' && final == '{' {
			st.listener.EraseRectangle(defaultParam(params, 0, 1), defaultParam(params, 1, 1), defaultParam(params, 2, 1), defaultParam(params, 3, 1))
			return
		}
		if st.csiIntermediate == '$' && final == 'x' {
			ch := rune(defaultParam(params, 0, 0))
			st.listener.FillRectangle(ch, defaultParam(params, 1, 1), defaultParam(params, 2, 1), defaultParam(params, 3, 1), defaultParam(params, 4, 1))
			return
		}
		if st.csiIntermediate == '$' && final == 'v' {
			srcTop := defaultParam(params, 0, 1)
			srcLeft := defaultParam(params, 1, 1)
			srcBottom := defaultParam(params, 2, 1)
			srcRight := defaultParam(params, 3, 1)
			dstTop := defaultParam(params, 5, 1)
			dstLeft := defaultParam(params, 6, 1)
			st.listener.CopyRectangle(srcTop, srcLeft, srcBottom, srcRight, dstTop, dstLeft)
			return
		}
		if st.csiIntermediate == '$' && final == 'y' {
			return
		}
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
	case '\x07', '\x08', '\t', '\n', '\x0b', '\x0c', '\r', '\x0e', '\x0f', '\x84', '\x85', '\x88', '\x8d', '\x90', '\x98', '\x9a', '\x9e', '\x9f':
		return false
	}
	return true
}
