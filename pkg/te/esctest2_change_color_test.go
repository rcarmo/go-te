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

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_Hash9
func TestEsctestChangeColorTestChangeColorHash9(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "#f00f00f00", "f0f0/f0f0/f0f0")
	esctestDoChangeColorTest(t, screen, stream, "0", "#800800800", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_Hash12
func TestEsctestChangeColorTestChangeColorHash12(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "#f000f000f000", "f0f0/f0f0/f0f0")
	esctestDoChangeColorTest(t, screen, stream, "0", "#800080008000", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_RGBI
func TestEsctestChangeColorTestChangeColorRGBI(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "rgbi:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeColorTest(t, screen, stream, "0", "rgbi:0.5/0.5/0.5", "c1c1/bbbb/bbbb")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_CIEXYZ
func TestEsctestChangeColorTestChangeColorCIEXYZ(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "CIEXYZ:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeColorTest(t, screen, stream, "0", "CIEXYZ:0.5/0.5/0.5", "dddd/b5b5/a0a0")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_CIEuvY
func TestEsctestChangeColorTestChangeColorCIEuvY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "CIEuvY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeColorTest(t, screen, stream, "0", "CIEuvY:0.5/0.5/0.5", "ffff/a3a3/aeae")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_CIExyY
func TestEsctestChangeColorTestChangeColorCIExyY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "CIExyY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeColorTest(t, screen, stream, "0", "CIExyY:0.5/0.5/0.5", "f7f7/b3b3/0e0e")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_CIELab
func TestEsctestChangeColorTestChangeColorCIELab(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "CIELab:1/1/1", "6c6c/6767/6767")
	esctestDoChangeColorTest(t, screen, stream, "0", "CIELab:0.5/0.5/0.5", "5252/4f4f/4f4f")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_CIELuv
func TestEsctestChangeColorTestChangeColorCIELuv(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "CIELuv:1/1/1", "1616/1414/0e0e")
	esctestDoChangeColorTest(t, screen, stream, "0", "CIELuv:0.5/0.5/0.5", "0e0e/1313/0e0e")
}

// From esctest2/esctest/tests/change_color.py::test_ChangeColor_TekHVC
func TestEsctestChangeColorTestChangeColorTekHVC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeColorTest(t, screen, stream, "0", "TekHVC:1/1/1", "1a1a/1313/0f0f")
	esctestDoChangeColorTest(t, screen, stream, "0", "TekHVC:0.5/0.5/0.5", "1111/1313/0e0e")
}
