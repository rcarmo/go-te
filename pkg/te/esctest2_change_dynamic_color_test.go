package te

import "testing"

func esctestDoChangeDynamicColorTest(t *testing.T, screen *Screen, stream *Stream, code, value, rgb string) {
	esctestChangeDynamicColor(t, stream, code, value)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeDynamicColor(t, stream, code, "?")
	})
	if got := esctestReadOSC(t, response, code); got != ";rgb:"+rgb {
		t.Fatalf("expected %q, got %q", ";rgb:"+rgb, got)
	}
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_Multiple
func TestEsctestChangeDynamicColorTestChangeDynamicColorMultiple(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestChangeDynamicColor(t, stream, "10", "rgb:f0f0/f0f0/f0f0", "rgb:f0f0/0000/0000")
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeDynamicColor(t, stream, "10", "?", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "10"); got != ";rgb:f0f0/f0f0/f0f0" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "11"); got != ";rgb:f0f0/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
	esctestChangeDynamicColor(t, stream, "10", "rgb:8080/8080/8080", "rgb:8080/0000/0000")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeDynamicColor(t, stream, "10", "?", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "10"); got != ";rgb:8080/8080/8080" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "11"); got != ";rgb:8080/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_RGB
func TestEsctestChangeDynamicColorTestChangeDynamicColorRGB(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "rgb:f0f0/f0f0/f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "rgb:8080/8080/8080", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_Hash3
func TestEsctestChangeDynamicColorTestChangeDynamicColorHash3(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#fff", "f0f0/f0f0/f0f0")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#888", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_Hash6
func TestEsctestChangeDynamicColorTestChangeDynamicColorHash6(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#f0f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#808080", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_Hash9
func TestEsctestChangeDynamicColorTestChangeDynamicColorHash9(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#f00f00f00", "f0f0/f0f0/f0f0")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#800800800", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_Hash12
func TestEsctestChangeDynamicColorTestChangeDynamicColorHash12(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#f000f000f000", "f0f0/f0f0/f0f0")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "#800080008000", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_RGBI
func TestEsctestChangeDynamicColorTestChangeDynamicColorRGBI(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "rgbi:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "rgbi:0.5/0.5/0.5", "c1c1/bbbb/bbbb")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_CIEXYZ
func TestEsctestChangeDynamicColorTestChangeDynamicColorCIEXYZ(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIEXYZ:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIEXYZ:0.5/0.5/0.5", "dddd/b5b5/a0a0")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_CIEuvY
func TestEsctestChangeDynamicColorTestChangeDynamicColorCIEuvY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIEuvY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIEuvY:0.5/0.5/0.5", "ffff/a3a3/aeae")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_CIExyY
func TestEsctestChangeDynamicColorTestChangeDynamicColorCIExyY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIExyY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIExyY:0.5/0.5/0.5", "f7f7/b3b3/0e0e")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_CIELab
func TestEsctestChangeDynamicColorTestChangeDynamicColorCIELab(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIELab:1/1/1", "6c6c/6767/6767")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIELab:0.5/0.5/0.5", "5252/4f4f/4f4f")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_CIELuv
func TestEsctestChangeDynamicColorTestChangeDynamicColorCIELuv(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIELuv:1/1/1", "1616/1414/0e0e")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "CIELuv:0.5/0.5/0.5", "0e0e/1313/0e0e")
}

// From esctest2/esctest/tests/change_dynamic_color.py::test_ChangeDynamicColor_TekHVC
func TestEsctestChangeDynamicColorTestChangeDynamicColorTekHVC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "TekHVC:1/1/1", "1a1a/1313/0f0f")
	esctestDoChangeDynamicColorTest(t, screen, stream, "10", "TekHVC:0.5/0.5/0.5", "1111/1313/0e0e")
}
