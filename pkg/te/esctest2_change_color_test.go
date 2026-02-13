package te

import "testing"

func esctestDoChangeColorTest(t *testing.T, screen *Screen, stream *Stream, color, value, rgb string) {
	esctestChangeColor(t, stream, color, value)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, color, "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != ";"+color+";rgb:"+rgb {
		t.Fatalf("expected %q, got %q", ";"+color+";rgb:"+rgb, got)
	}
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_Multiple
func TestEsctestChangeColorTestChangeColorMultiple(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestChangeColor(t, stream, "0", "rgb:f0f0/f0f0/f0f0", "1", "rgb:f0f0/0000/0000")
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeColor(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "4"); got != ";0;rgb:f0f0/f0f0/f0f0" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "4"); got != ";1;rgb:f0f0/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
	esctestChangeColor(t, stream, "0", "rgb:8080/8080/8080", "1", "rgb:8080/0000/0000")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeColor(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "4"); got != ";0;rgb:8080/8080/8080" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "4"); got != ";1;rgb:8080/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_RGB
func TestEsctestChangeColorTestChangeColorRGB(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "rgb:f0f0/f0f0/f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeColorTest(t, screen, stream, "0", "rgb:8080/8080/8080", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_Hash3
func TestEsctestChangeColorTestChangeColorHash3(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "#fff", "f0f0/f0f0/f0f0")
	esctestDoChangeColorTest(t, screen, stream, "0", "#888", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_Hash6
func TestEsctestChangeColorTestChangeColorHash6(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "#f0f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeColorTest(t, screen, stream, "0", "#808080", "8080/8080/8080")
}
