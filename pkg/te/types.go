package te

type ColorMode uint8

const (
	ColorDefault ColorMode = iota
	ColorANSI16
	ColorANSI256
	ColorTrueColor
)

type Color struct {
	Mode  ColorMode
	Index uint8
	Name  string
}

type Attr struct {
	Fg            Color
	Bg            Color
	Bold          bool
	Italics       bool
	Underline     bool
	Strikethrough bool
	Reverse       bool
	Blink         bool
	Conceal       bool
	Protected     bool
	ISOProtected  bool
}

type Cell struct {
	Data string
	Attr Attr
}

type Cursor struct {
	Row    int
	Col    int
	Attr   Attr
	Hidden bool
}
