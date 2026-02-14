package te

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	DefaultSVGForeground = "#d3d7cf"
	DefaultSVGBackground = "#000000"
	DefaultSVGFontFamily = "ui-monospace, \"SFMono-Regular\", \"FiraCode Nerd Font\", \"FiraMono Nerd Font\", \"Fira Code\", \"Roboto Mono\", Menlo, Monaco, Consolas, \"Liberation Mono\", \"DejaVu Sans Mono\", \"Courier New\", monospace"
	DefaultSVGFontSize   = 14
	DefaultSVGLineHeight = 1.2
	DefaultSVGCharWidth  = 8
	DefaultSVGPadding    = 10
)

var ansiSVGColors = map[string]string{
	"black":         "#000000",
	"red":           "#cc0000",
	"green":         "#4e9a06",
	"yellow":        "#c4a000",
	"blue":          "#3465a4",
	"magenta":       "#75507b",
	"cyan":          "#06989a",
	"white":         "#d3d7cf",
	"brightblack":   "#555753",
	"brightred":     "#ef2929",
	"brightgreen":   "#8ae234",
	"brightyellow":  "#fce94f",
	"brightblue":    "#729fcf",
	"brightmagenta": "#ad7fa8",
	"brightcyan":    "#34e2e2",
	"brightwhite":   "#eeeeec",
	"gray":          "#555753",
	"grey":          "#555753",
	"lightgray":     "#d3d7cf",
	"lightgrey":     "#d3d7cf",
	"brown":         "#c4a000",
	"brightbrown":   "#fce94f",
}

var boxDrawingChars = "─━│┃┄┅┆┇┈┉┊┋┌┍┎┏┐┑┒┓└┕┖┗┘┙┚┛├┝┞┟┠┡┢┣┤┥┦┧┨┩┪┫┬┭┮┯┰┱┲┳┴┵┶┷┸┹┺┻┼┽┾┿╀╁╂╃╄╅╆╇╈╉╊╋═║╒╓╔╕╖╗╘╙╚╛╜╝╞╟╠╡╢╣╤╥╦╧╨╩╪╫╬╭╮╯╰╱╲╳╴╵╶╷╸╹╺╻╼╽╾╿"

type SVGOptions struct {
	Title      string
	FontFamily string
	FontSize   int
	LineHeight float64
	CharWidth  float64
	Padding    float64
	Background string
	Foreground string
	Palette    map[string]string
}

func DefaultSVGOptions() SVGOptions {
	return SVGOptions{
		Title:      "Terminal",
		FontFamily: DefaultSVGFontFamily,
		FontSize:   DefaultSVGFontSize,
		LineHeight: DefaultSVGLineHeight,
		CharWidth:  DefaultSVGCharWidth,
		Padding:    DefaultSVGPadding,
		Background: DefaultSVGBackground,
		Foreground: DefaultSVGForeground,
	}
}

func RenderScreenSVG(screen *Screen, opts SVGOptions) string {
	if screen == nil {
		return ""
	}
	return RenderTerminalSVG(screen.Buffer, screen.Columns, screen.Lines, opts)
}

func RenderTerminalSVG(screenBuffer [][]Cell, width, height int, opts SVGOptions) string {
	opts = normalizeSVGOptions(opts)
	palette := normalizePalette(opts.Palette)

	actualLineHeight := float64(opts.FontSize) * opts.LineHeight
	svgWidth := float64(width)*opts.CharWidth + 2*opts.Padding
	svgHeight := float64(height)*actualLineHeight + 2*opts.Padding

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %s %s\" class=\"terminal-svg\">", formatFloat(svgWidth), formatFloat(svgHeight)))
	b.WriteString("<title>")
	b.WriteString(escapeXML(opts.Title))
	b.WriteString("</title>")
	b.WriteString("<defs><style>")
	b.WriteString(fmt.Sprintf(".terminal-bg { fill: %s; }", opts.Background))
	b.WriteString(fmt.Sprintf(".terminal-text { font-family: %s; font-size: %dpx; fill: %s; white-space: pre; text-rendering: optimizeLegibility; }", opts.FontFamily, opts.FontSize, opts.Foreground))
	b.WriteString(".bold { font-weight: bold; }.italic { font-style: italic; }.underline { text-decoration: underline; }.strikethrough { text-decoration: line-through; }")
	b.WriteString("</style></defs>")
	b.WriteString(fmt.Sprintf("<rect class=\"terminal-bg\" x=\"0\" y=\"0\" width=\"%s\" height=\"%s\"/>", formatFloat(svgWidth), formatFloat(svgHeight)))
	b.WriteString("<g class=\"terminal-text\">")

	for rowIdx := 0; rowIdx < height && rowIdx < len(screenBuffer); rowIdx++ {
		rowData := screenBuffer[rowIdx]
		if len(rowData) == 0 {
			continue
		}
		rowWidth := width
		if len(rowData) < rowWidth {
			rowWidth = len(rowData)
		}

		rectY := opts.Padding + float64(rowIdx)*actualLineHeight
		textY := rectY + float64(opts.FontSize)

		var rowBackgrounds []string
		var rowTspans []string

		for col := 0; col < rowWidth; {
			cell := rowData[col]
			charData := cell.Data
			if charData == "" {
				col++
				continue
			}

			x := opts.Padding + float64(col)*opts.CharWidth
			fg := colorToHex(cell.Attr.Fg, true, palette, opts.Foreground, opts.Background)
			bg := colorToHex(cell.Attr.Bg, false, palette, opts.Foreground, opts.Background)
			if cell.Attr.Reverse {
				fg, bg = bg, fg
			}

			charCols := 1
			if col+1 < rowWidth && rowData[col+1].Data == "" {
				charCols = 2
			}

			if bg != opts.Background {
				bgWidth := float64(charCols)*opts.CharWidth + 0.5
				rowBackgrounds = append(rowBackgrounds, fmt.Sprintf("<rect x=\"%s\" y=\"%s\" width=\"%s\" height=\"%s\" fill=\"%s\"/>", formatFloat(x), formatFloat(rectY), formatFloat(bgWidth), formatFloat(actualLineHeight+0.5), bg))
			}

			classes := []string{}
			if cell.Attr.Bold {
				classes = append(classes, "bold")
			}
			if cell.Attr.Italics {
				classes = append(classes, "italic")
			}
			if cell.Attr.Underline {
				classes = append(classes, "underline")
			}
			if cell.Attr.Strikethrough {
				classes = append(classes, "strikethrough")
			}

			if isBoxDrawing(charData) {
				fillAttr := ""
				if fg != opts.Foreground {
					fillAttr = fmt.Sprintf(" fill=\"%s\"", fg)
				}
				classAttr := ""
				if len(classes) > 0 {
					classAttr = fmt.Sprintf(" class=\"%s\"", strings.Join(classes, " "))
				}
				rowBackgrounds = append(rowBackgrounds, fmt.Sprintf("<text x=\"%s\" y=\"%s\" transform=\"translate(0,%s) scale(1,%s) translate(0,%s)\"%s%s>%s</text>", formatFloat(x), formatFloat(textY), formatFloat(rectY), formatFloat(opts.LineHeight), formatFloat(-rectY), fillAttr, classAttr, escapeXML(charData)))
			} else {
				attrs := []string{fmt.Sprintf("x=\"%s\"", formatFloat(x))}
				if fg != opts.Foreground {
					attrs = append(attrs, fmt.Sprintf("fill=\"%s\"", fg))
				}
				if len(classes) > 0 {
					attrs = append(attrs, fmt.Sprintf("class=\"%s\"", strings.Join(classes, " ")))
				}
				rowTspans = append(rowTspans, fmt.Sprintf("<tspan %s>%s</tspan>", strings.Join(attrs, " "), escapeXML(charData)))
			}

			col += charCols
		}

		if len(rowBackgrounds) > 0 || len(rowTspans) > 0 {
			for _, rect := range rowBackgrounds {
				b.WriteString(rect)
			}
			if len(rowTspans) > 0 {
				b.WriteString(fmt.Sprintf("<text y=\"%s\">", formatFloat(textY)))
				for _, span := range rowTspans {
					b.WriteString(span)
				}
				b.WriteString("</text>")
			}
		}
	}

	b.WriteString("</g></svg>")
	return b.String()
}

func normalizeSVGOptions(opts SVGOptions) SVGOptions {
	defaults := DefaultSVGOptions()
	if opts.Title == "" {
		opts.Title = defaults.Title
	}
	if opts.FontFamily == "" {
		opts.FontFamily = defaults.FontFamily
	}
	if opts.FontSize == 0 {
		opts.FontSize = defaults.FontSize
	}
	if opts.LineHeight == 0 {
		opts.LineHeight = defaults.LineHeight
	}
	if opts.CharWidth == 0 {
		opts.CharWidth = defaults.CharWidth
	}
	if opts.Padding == 0 {
		opts.Padding = defaults.Padding
	}
	if opts.Background == "" {
		opts.Background = defaults.Background
	}
	if opts.Foreground == "" {
		opts.Foreground = defaults.Foreground
	}
	return opts
}

func normalizePalette(palette map[string]string) map[string]string {
	if palette == nil {
		return ansiSVGColors
	}
	normalized := make(map[string]string, len(palette)+6)
	for key, value := range palette {
		normalized[strings.ToLower(key)] = value
	}
	ensurePaletteAlias(normalized, "gray", "brightblack", ansiSVGColors["gray"])
	ensurePaletteAlias(normalized, "grey", "brightblack", ansiSVGColors["grey"])
	ensurePaletteAlias(normalized, "lightgray", "white", ansiSVGColors["lightgray"])
	ensurePaletteAlias(normalized, "lightgrey", "white", ansiSVGColors["lightgrey"])
	ensurePaletteAlias(normalized, "brown", "yellow", ansiSVGColors["brown"])
	ensurePaletteAlias(normalized, "brightbrown", "brightyellow", ansiSVGColors["brightbrown"])
	return normalized
}

func ensurePaletteAlias(palette map[string]string, alias, source, fallback string) {
	if _, ok := palette[alias]; ok {
		return
	}
	if value, ok := palette[source]; ok {
		palette[alias] = value
		return
	}
	palette[alias] = fallback
}

func colorToHex(color Color, isForeground bool, palette map[string]string, defaultFg, defaultBg string) string {
	fallback := defaultFg
	if !isForeground {
		fallback = defaultBg
	}
	if color.Mode == ColorDefault || color.Name == "" || color.Name == "default" {
		return fallback
	}
	name := color.Name
	if strings.HasPrefix(name, "#") {
		return name
	}
	if isHexColor(name) {
		return "#" + name
	}
	lower := strings.ToLower(name)
	if lower == "default" {
		return fallback
	}
	if lowerPaletteValue, ok := palette[lower]; ok {
		return lowerPaletteValue
	}
	if strings.HasPrefix(lower, "rgb(") {
		return fallback
	}
	return fallback
}

func isHexColor(value string) bool {
	if len(value) != 6 {
		return false
	}
	for _, r := range value {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}

func escapeXML(value string) string {
	return html.EscapeString(value)
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 1, 64)
}

func isBoxDrawing(char string) bool {
	r, size := utf8.DecodeRuneInString(char)
	if r == utf8.RuneError || size == 0 || size != len(char) {
		return false
	}
	return strings.ContainsRune(boxDrawingChars, r)
}
