package te

type ColorMode uint8

const (
	ColorDefault ColorMode = iota
	ColorANSI16
	ColorANSI256
)

type Color struct {
	Mode  ColorMode
	Index uint8
}

type Attr struct {
	Fg        Color
	Bg        Color
	Bold      bool
	Underline bool
	Blink     bool
	Reverse   bool
	Conceal   bool
}

type Cell struct {
	Ch   rune
	Attr Attr
}

type Cursor struct {
	Row    int
	Col    int
	Hidden bool
}
