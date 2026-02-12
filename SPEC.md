# go-te

## Overview

go-te is a Go port of https://github.com/selectel/pyte, a small in-memory terminal emulator.
It emulates a VT100/VT220/VT520-style screen (TERM=linux-like behavior) for screen scraping,
TUI rendering, and terminal playback without a real PTY.

## Goals

- Match pyte behavior for screen state, cursor movement, and attribute handling.
- Provide a stream parser that accepts terminal byte streams and updates a screen model.
- Offer screen variants for history (scrollback) and diffs (dirty line tracking).
- Keep the API minimal, deterministic, and easy to embed.

## Non-Goals

- No PTY, shell, or process management.
- No input handling (keyboard/mouse events are out of scope).
- No full xterm feature set; unsupported sequences are ignored.

## Compatibility Scope (pyte-aligned)

- VT100/VT220/VT520 control and CSI sequences for cursor, erase, insert/delete, and SGR.
- TERM=linux style semantics (line wrapping, scrolling, and tab stops).
- 16-color + 256-color SGR support; default/bright colors must be supported.
- Common SGR attributes: bold, underline, blink, reverse, and conceal.

## Architecture

### Stream Parser

- `Stream` is a state machine that consumes decoded runes and emits screen actions.
- `ByteStream` wraps `Stream` and accepts raw bytes, decoding UTF-8 with buffering for
  partial sequences.
- Supported control classes:
  - C0 controls: BEL, BS, HT, LF, VT, FF, CR, ESC.
  - ESC sequences for save/restore cursor, alternate charset toggles, and mode changes.
  - CSI sequences for cursor movement, scrolling region, erase, insert/delete, and SGR.
- Unsupported or malformed sequences are ignored; in strict mode they return errors.

### Screen Model

- The screen is a 2D grid of cells sized by `columns` x `lines`.
- Each cell contains a rune and `Attr` (style + colors). Empty cells contain a space rune
  with default attributes.
- The `Cursor` tracks row/column (0-based), visibility, and style state for future writes.
- The screen tracks modes:
  - insert/overwrite
  - origin mode (relative to scroll region)
  - autowrap
  - newline mode (LF implies CR)
- Tab stops default every 8 columns and can be set/cleared.

### Buffers

- Primary and alternate screen buffers are supported, including cursor save/restore.
- A scroll region (top/bottom margins) constrains scrolling operations.

### Screen Variants

- `Screen`: base screen implementation.
- `DiffScreen`: tracks which line indexes changed since the last clear.
- `HistoryScreen`: retains scrollback in a bounded ring buffer.
- `DebugScreen`: logs received actions and state changes for diagnostics.

## Public API (Current)

### Core Types

```
package te

type Color struct {
  Mode  ColorMode // Default, ANSI16, ANSI256, TrueColor
  Index uint8
  Name  string
}

type Attr struct {
  Fg, Bg       Color
  Bold         bool
  Italics      bool
  Underline    bool
  Strikethrough bool
  Reverse      bool
  Blink        bool
  Conceal      bool
}

type Cell struct {
  Data string
  Attr Attr
}

type Cursor struct {
  Row, Col int
  Attr     Attr
  Hidden   bool
}
```

### Screen

```
func NewScreen(cols, lines int) *Screen
func (s *Screen) Resize(lines, cols int)
func (s *Screen) Reset()
func (s *Screen) Display() []string
func (s *Screen) LinesCells() [][]Cell
```

- `Display()` returns the visible screen as `[]string`, one entry per line.
- `LinesCells()` exposes the full cell matrix for renderers needing attributes.
- `Screen` exposes fields for buffer, modes, cursor, margins, tab stops, title, and icon name.

### HistoryScreen

```
func NewHistoryScreen(cols, lines, history int) *HistoryScreen
func NewHistoryScreenWithRatio(cols, lines, history int, ratio float64) *HistoryScreen
func (s *HistoryScreen) History() [][]Cell
func (s *HistoryScreen) Scrollback() int
```

### DiffScreen

```
func NewDiffScreen(cols, lines int) *DiffScreen
```

### Stream

```
func NewStream(screen EventHandler, strict bool) *Stream
func (st *Stream) Attach(screen EventHandler)
func (st *Stream) Detach(screen EventHandler)
func (st *Stream) Feed(data string) error
```

### ByteStream

```
func NewByteStream(screen EventHandler, strict bool) *ByteStream
func (st *ByteStream) Feed(data []byte) error
func (st *ByteStream) SelectOtherCharset(code string)
```

`EventHandler` is implemented by Screen/HistoryScreen/DiffScreen/DebugScreen.

## Control/CSI Coverage (Minimum)

- Cursor: CUU, CUD, CUF, CUB, CUP/HVP, CNL, CPL, CHA, VPA, SCP/RCP.
- Erase: ED, EL, ECH, DCH.
- Insert/Delete: ICH, IL, DL.
- Scrolling: SU, SD, DECSTBM (set scroll region).
- Modes: DECAWM (autowrap), DECOM (origin), SM/RM for insert mode.
- SGR: reset, bold, underline, blink, reverse, conceal, 16-color, 256-color.

## Behavior Notes

- Line feed scrolls within the active scroll region when the cursor is on the last row.
- Autowrap inserts a wrap when writing past the last column.
- `CR` moves the cursor to column 0; `LF` moves down one row; `BS` moves left one column.
- `HT` moves to the next tab stop or last column if none remain.
- Insert mode shifts existing cells right before writing a new cell.
- Erase operations clear to space with default attributes.

## Concurrency

- Screen and Stream types are not thread-safe; callers must serialize access.
- Stream methods must not mutate the screen concurrently with reads of the buffer.

## Testing Strategy

- Port the full pyte test suite as Go tests (screen/stream/history/diff).
- Validate captured fixtures from `pyte/tests/captured` via ByteStream playback.
- Add esctest2-derived tests to exercise full escape/CSI coverage (VT100/VT220/xterm extensions).
- Keep parity tests for cursor movement, SGR, erase, insert/delete, scrolling, and history pagination.
