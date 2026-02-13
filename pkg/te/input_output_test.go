package te

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// From pyte/tests/test_input_output.py::test_input_output
func TestCapturedInputOutput(t *testing.T) {
	cases := []string{"cat-gpl3", "find-etc", "htop", "ls", "mc", "top", "vi"}
	for _, name := range cases {
		name := name
		t.Run(name, func(t *testing.T) {
			inputPath := filepath.Join("..", "..", "vendor", "pyte", "tests", "captured", name+".input")
			outputPath := filepath.Join("..", "..", "vendor", "pyte", "tests", "captured", name+".output")

			input, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("read input: %v", err)
			}
			outputRaw, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("read output: %v", err)
			}
			var output []string
			if err := json.Unmarshal(outputRaw, &output); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}

			screen := NewScreen(80, 24)
			stream := NewByteStream(screen, false)
			if err := stream.Feed(input); err != nil {
				t.Fatalf("feed input: %v", err)
			}

			if got := screen.Display(); len(got) != len(output) {
				t.Fatalf("expected %d lines, got %d", len(output), len(got))
			}
			for i := range output {
				if output[i] != screen.Display()[i] {
					t.Fatalf("line %d mismatch: expected %q got %q", i, output[i], screen.Display()[i])
				}
			}
		})
	}
}
