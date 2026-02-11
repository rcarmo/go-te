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

type History struct {
	Top      historyDeque
	Bottom   historyDeque
	Ratio    float64
	Size     int
	Position int
}

type HistoryScreen struct {
	*Screen
	history History
}

func NewHistoryScreen(cols, lines, history int) *HistoryScreen {
	return NewHistoryScreenWithRatio(cols, lines, history, 0.5)
}

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

func (h *HistoryScreen) Reset() {
	h.Screen.Reset()
	h.resetHistory()
}

func (h *HistoryScreen) EraseInDisplay(how int, private bool, rest ...int) {
	h.beforeEvent("erase_in_display")
	h.Screen.EraseInDisplay(how, private, rest...)
	if how == 3 {
		h.resetHistory()
	}
	h.afterEvent("erase_in_display")
}

func (h *HistoryScreen) Index() {
	h.beforeEvent("index")
	h.indexInternal()
	h.afterEvent("index")
}

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

func (h *HistoryScreen) Draw(data string) {
	h.beforeEvent("draw")
	h.Screen.Draw(data)
	h.afterEvent("draw")
}

func (h *HistoryScreen) LineFeed() {
	h.beforeEvent("linefeed")
	h.indexInternal()
	if h.isModeSet(ModeLNM) {
		h.CarriageReturn()
	}
	h.afterEvent("linefeed")
}

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

func (h *HistoryScreen) History() [][]Cell {
	return append([][]Cell(nil), h.history.Top.items...)
}

func (h *HistoryScreen) Scrollback() int {
	return len(h.history.Top.items)
}

func (h *HistoryScreen) resetHistory() {
	h.history.Top.items = nil
	h.history.Bottom.items = nil
	h.history.Position = h.history.Size
}
