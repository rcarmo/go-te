package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rcarmo/go-te/pkg/te"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: render_svg <capture-file> <output-svg> [width] [height]")
		os.Exit(2)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	width := 120
	height := 40

	if len(os.Args) >= 4 {
		if v, err := strconv.Atoi(os.Args[3]); err == nil {
			width = v
		}
	}
	if len(os.Args) >= 5 {
		if v, err := strconv.Atoi(os.Args[4]); err == nil {
			height = v
		}
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read capture file: %v\n", err)
		os.Exit(1)
	}

	screen := te.NewScreen(width, height)
	stream := te.NewByteStream(screen, false)
	if err := stream.Feed(data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse capture data: %v\n", err)
		os.Exit(1)
	}

	opts := te.DefaultSVGOptions()
	opts.Title = "copilot-tmux"
	svg := te.RenderScreenSVG(screen, opts)

	if err := os.WriteFile(outputPath, []byte(svg), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write SVG: %v\n", err)
		os.Exit(1)
	}
}
