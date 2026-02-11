package te

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidSequence = errors.New("invalid control sequence")

type ScreenLike interface {
	PutRune(ch rune)
	Backspace()
	Tab()
	LineFeed()
	CarriageReturn()
	MoveCursor(rowDelta, colDelta int)
	SetCursor(row, col int)
	Cursor() Cursor
	Size() (cols, lines int)
	SaveCursor()
	RestoreCursor()
	EraseInDisplay(mode int)
	EraseInLine(mode int)
	EraseChars(n int)
	DeleteChars(n int)
	InsertChars(n int)
	DeleteLines(n int)
	InsertLines(n int)
	ScrollUp(n int)
	ScrollDown(n int)
	SetScrollRegion(top, bottom int)
	SetTabStop(col int)
	ClearTabStop(col int)
	SetInsertMode(enabled bool)
	SetAutowrap(enabled bool)
	SetOriginMode(enabled bool)
	SetNewlineMode(enabled bool)
	ApplySGR(params []int)
	Reset()
}

type Stream struct {
	screen       ScreenLike
	strict       bool
	state        parserState
	params       []int
	paramBuf     strings.Builder
	private      bool
	seenParam    bool
	intermediate strings.Builder
}

type tabStopClearAll interface {
	ClearAllTabStops()
}

type alternateBuffer interface {
	EnableAlternateBuffer(clear bool)
	DisableAlternateBuffer()
}

type parserState int

const (
	stateGround parserState = iota
	stateEscape
	stateCSI
	stateEscapeCharset
)

func NewStream(screen ScreenLike, strict bool) *Stream {
	return &Stream{screen: screen, strict: strict}
}

func (st *Stream) Attach(screen ScreenLike) {
	st.screen = screen
}

func (st *Stream) FeedString(data string) error {
	if st.screen == nil {
		return nil
	}
	for _, ch := range data {
		if err := st.feedRune(ch); err != nil {
			return err
		}
	}
	return nil
}

func (st *Stream) feedRune(ch rune) error {
	switch st.state {
	case stateGround:
		return st.handleGround(ch)
	case stateEscape:
		return st.handleEscape(ch)
	case stateCSI:
		return st.handleCSI(ch)
	case stateEscapeCharset:
		return st.handleCharset(ch)
	default:
		st.state = stateGround
	}
	return nil
}

func (st *Stream) handleGround(ch rune) error {
	switch ch {
	case 0x07:
		return nil
	case 0x08:
		st.screen.Backspace()
		return nil
	case 0x09:
		st.screen.Tab()
		return nil
	case 0x0a, 0x0b, 0x0c:
		st.screen.LineFeed()
		return nil
	case 0x0d:
		st.screen.CarriageReturn()
		return nil
	case 0x1b:
		st.state = stateEscape
		return nil
	case 0x7f:
		return nil
	}
	if ch >= 0x20 {
		st.screen.PutRune(ch)
	}
	return nil
}

func (st *Stream) handleEscape(ch rune) error {
	st.state = stateGround
	switch ch {
	case '[':
		st.state = stateCSI
		st.params = st.params[:0]
		st.paramBuf.Reset()
		st.private = false
		st.seenParam = false
		st.intermediate.Reset()
		return nil
	case '7':
		st.screen.SaveCursor()
		return nil
	case '8':
		st.screen.RestoreCursor()
		return nil
	case 'D':
		st.screen.LineFeed()
		return nil
	case 'M':
		st.screen.ScrollDown(1)
		return nil
	case 'E':
		st.screen.CarriageReturn()
		st.screen.LineFeed()
		return nil
	case 'c':
		st.screen.Reset()
		return nil
	case 'H':
		cursor := st.screen.Cursor()
		st.screen.SetTabStop(cursor.Col)
		return nil
	case '(', ')':
		st.state = stateEscapeCharset
		return nil
	default:
		if st.strict {
			return ErrInvalidSequence
		}
		return nil
	}
}

func (st *Stream) handleCharset(ch rune) error {
	st.state = stateGround
	_ = ch
	return nil
}

func (st *Stream) handleCSI(ch rune) error {
	if ch >= 0x30 && ch <= 0x3f {
		st.seenParam = true
		if ch == '?' && st.paramBuf.Len() == 0 && len(st.params) == 0 {
			st.private = true
			return nil
		}
		st.paramBuf.WriteRune(ch)
		return nil
	}
	if ch >= 0x20 && ch <= 0x2f {
		st.intermediate.WriteRune(ch)
		return nil
	}
	if ch >= 0x40 && ch <= 0x7e {
		params := st.parseParams()
		st.state = stateGround
		err := st.dispatchCSI(ch, params)
		st.private = false
		return err
	}
	st.state = stateGround
	st.private = false
	if st.strict {
		return ErrInvalidSequence
	}
	return nil
}

func (st *Stream) parseParams() []int {
	if !st.seenParam || st.paramBuf.Len() == 0 {
		return nil
	}
	parts := strings.Split(st.paramBuf.String(), ";")
	params := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			params = append(params, 0)
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			params = append(params, 0)
			continue
		}
		params = append(params, value)
	}
	return params
}

func (st *Stream) dispatchCSI(final rune, params []int) error {
	switch final {
	case 'A':
		st.screen.MoveCursor(-defaultParam(params, 0, 1), 0)
	case 'B':
		st.screen.MoveCursor(defaultParam(params, 0, 1), 0)
	case 'C':
		st.screen.MoveCursor(0, defaultParam(params, 0, 1))
	case 'D':
		st.screen.MoveCursor(0, -defaultParam(params, 0, 1))
	case 'E':
		st.screen.MoveCursor(defaultParam(params, 0, 1), 0)
		st.screen.CarriageReturn()
	case 'F':
		st.screen.MoveCursor(-defaultParam(params, 0, 1), 0)
		st.screen.CarriageReturn()
	case 'G':
		st.screen.SetCursor(st.screen.Cursor().Row, defaultParam(params, 0, 1)-1)
	case 'd':
		st.screen.SetCursor(defaultParam(params, 0, 1)-1, st.screen.Cursor().Col)
	case 'H', 'f':
		row := defaultParam(params, 0, 1) - 1
		col := defaultParam(params, 1, 1) - 1
		st.screen.SetCursor(row, col)
	case 'J':
		st.screen.EraseInDisplay(defaultParam(params, 0, 0))
	case 'K':
		st.screen.EraseInLine(defaultParam(params, 0, 0))
	case 'g':
		mode := defaultParam(params, 0, 0)
		switch mode {
		case 0:
			cursor := st.screen.Cursor()
			st.screen.ClearTabStop(cursor.Col)
		case 3:
			if target, ok := st.screen.(tabStopClearAll); ok {
				target.ClearAllTabStops()
			}
		}
	case 'X':
		st.screen.EraseChars(defaultParam(params, 0, 1))
	case 'P':
		st.screen.DeleteChars(defaultParam(params, 0, 1))
	case '@':
		st.screen.InsertChars(defaultParam(params, 0, 1))
	case 'L':
		st.screen.InsertLines(defaultParam(params, 0, 1))
	case 'M':
		st.screen.DeleteLines(defaultParam(params, 0, 1))
	case 'S':
		st.screen.ScrollUp(defaultParam(params, 0, 1))
	case 'T':
		st.screen.ScrollDown(defaultParam(params, 0, 1))
	case 'r':
		_, lines := st.screen.Size()
		if len(params) == 0 {
			st.screen.SetScrollRegion(0, lines-1)
			return nil
		}
		top := defaultParam(params, 0, 1) - 1
		bottom := defaultParam(params, 1, lines) - 1
		st.screen.SetScrollRegion(top, bottom)
	case 'm':
		st.screen.ApplySGR(params)
	case 's':
		st.screen.SaveCursor()
	case 'u':
		st.screen.RestoreCursor()
	case 'h':
		return st.setMode(params, true)
	case 'l':
		return st.setMode(params, false)
	default:
		if st.strict {
			return fmt.Errorf("%w: CSI %c", ErrInvalidSequence, final)
		}
	}
	return nil
}

func (st *Stream) setMode(params []int, enabled bool) error {
	if st.private {
		for _, p := range params {
			switch p {
			case 6:
				st.screen.SetOriginMode(enabled)
			case 7:
				st.screen.SetAutowrap(enabled)
			case 47, 1047:
				if target, ok := st.screen.(alternateBuffer); ok {
					if enabled {
						target.EnableAlternateBuffer(true)
					} else {
						target.DisableAlternateBuffer()
					}
				}
			case 1049:
				if target, ok := st.screen.(alternateBuffer); ok {
					if enabled {
						st.screen.SaveCursor()
						target.EnableAlternateBuffer(true)
					} else {
						target.DisableAlternateBuffer()
						st.screen.RestoreCursor()
					}
				}
			}
		}
		return nil
	}
	for _, p := range params {
		switch p {
		case 4:
			st.screen.SetInsertMode(enabled)
		case 20:
			st.screen.SetNewlineMode(enabled)
		}
	}
	return nil
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
