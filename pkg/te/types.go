package te

// ColorMode identifies a color encoding mode.
type ColorMode uint8

const (
	// ColorDefault selects default terminal colors.
	ColorDefault ColorMode = iota
	// ColorANSI16 selects the 16-color palette.
	ColorANSI16
	// ColorANSI256 selects the 256-color palette.
	ColorANSI256
	// ColorTrueColor selects 24-bit true color.
	ColorTrueColor
)

// Color describes a terminal color.
type Color struct {
	Mode  ColorMode
	Index uint8
	Name  string
}

// Attr describes character attributes and colors.
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

// Cell represents a single screen cell.
type Cell struct {
	Data string
	Attr Attr
}

// Cursor tracks cursor position and attributes.
type Cursor struct {
	Row    int
	Col    int
	Attr   Attr
	Hidden bool
}
