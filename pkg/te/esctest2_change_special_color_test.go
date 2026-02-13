package te

import "testing"

func esctestDoChangeSpecialColorTest(t *testing.T, screen *Screen, stream *Stream, color, value, rgb string) {
	offset := esctestGetIndexedColors()
	esctestChangeSpecialColor(t, stream, color, value)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, color, "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != ";"+esctestItoa(offset+esctestAtoi(color))+";rgb:"+rgb {
		t.Fatalf("expected %q, got %q", ";"+esctestItoa(offset+esctestAtoi(color))+";rgb:"+rgb, got)
	}
}

func esctestDoChangeSpecialColorTest2(t *testing.T, screen *Screen, stream *Stream, color, value, rgb string) {
	esctestChangeSpecialColor2(t, stream, color, value)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor2(t, stream, color, "?")
	})
	if got := esctestReadOSC(t, response, "5"); got != ";"+color+";rgb:"+rgb {
		t.Fatalf("expected %q, got %q", ";"+color+";rgb:"+rgb, got)
	}
}

func esctestDoChangeSpecialColorTests(t *testing.T, screen *Screen, stream *Stream, color, value, rgb string) {
	esctestDoChangeSpecialColorTest(t, screen, stream, color, value, rgb)
	esctestDoChangeSpecialColorTest2(t, screen, stream, color, value, rgb)
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Multiple
func TestEsctestChangeSpecialColorTestChangeSpecialColorMultiple(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	offset := esctestGetIndexedColors()
	esctestChangeSpecialColor(t, stream, "0", "rgb:f0f0/f0f0/f0f0", "1", "rgb:f0f0/0000/0000")
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "4"); got != ";"+esctestItoa(offset)+";rgb:f0f0/f0f0/f0f0" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "4"); got != ";"+esctestItoa(offset+1)+";rgb:f0f0/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}

	esctestChangeSpecialColor(t, stream, "0", "rgb:8080/8080/8080", "1", "rgb:8080/0000/0000")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "4"); got != ";"+esctestItoa(offset)+";rgb:8080/8080/8080" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "4"); got != ";"+esctestItoa(offset+1)+";rgb:8080/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Multiple2
func TestEsctestChangeSpecialColorTestChangeSpecialColorMultiple2(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestChangeSpecialColor2(t, stream, "0", "rgb:f0f0/f0f0/f0f0", "1", "rgb:f0f0/0000/0000")
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor2(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "5"); got != ";0;rgb:f0f0/f0f0/f0f0" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "5"); got != ";1;rgb:f0f0/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}

	esctestChangeSpecialColor2(t, stream, "0", "rgb:8080/8080/8080", "1", "rgb:8080/0000/0000")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor2(t, stream, "0", "?", "1", "?")
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if got := esctestReadOSC(t, responses[0], "5"); got != ";0;rgb:8080/8080/8080" {
		t.Fatalf("unexpected response %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "5"); got != ";1;rgb:8080/0000/0000" {
		t.Fatalf("unexpected response %q", got)
	}
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_RGB
func TestEsctestChangeSpecialColorTestChangeSpecialColorRGB(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "rgb:f0f0/f0f0/f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "rgb:8080/8080/8080", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Hash3
func TestEsctestChangeSpecialColorTestChangeSpecialColorHash3(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#fff", "f0f0/f0f0/f0f0")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#888", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Hash6
func TestEsctestChangeSpecialColorTestChangeSpecialColorHash6(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#f0f0f0", "f0f0/f0f0/f0f0")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#808080", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Hash9
func TestEsctestChangeSpecialColorTestChangeSpecialColorHash9(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#f00f00f00", "f0f0/f0f0/f0f0")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#800800800", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_Hash12
func TestEsctestChangeSpecialColorTestChangeSpecialColorHash12(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#f000f000f000", "f0f0/f0f0/f0f0")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "#800080008000", "8080/8080/8080")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_RGBI
func TestEsctestChangeSpecialColorTestChangeSpecialColorRGBI(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "rgbi:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "rgbi:0.5/0.5/0.5", "c1c1/bbbb/bbbb")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_CIEXYZ
func TestEsctestChangeSpecialColorTestChangeSpecialColorCIEXYZ(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIEXYZ:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIEXYZ:0.5/0.5/0.5", "dddd/b5b5/a0a0")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_CIEuvY
func TestEsctestChangeSpecialColorTestChangeSpecialColorCIEuvY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIEuvY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIEuvY:0.5/0.5/0.5", "ffff/a3a3/aeae")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_CIExyY
func TestEsctestChangeSpecialColorTestChangeSpecialColorCIExyY(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIExyY:1/1/1", "ffff/ffff/ffff")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIExyY:0.5/0.5/0.5", "f7f7/b3b3/0e0e")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_CIELab
func TestEsctestChangeSpecialColorTestChangeSpecialColorCIELab(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIELab:1/1/1", "6c6c/6767/6767")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIELab:0.5/0.5/0.5", "5252/4f4f/4f4f")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_CIELuv
func TestEsctestChangeSpecialColorTestChangeSpecialColorCIELuv(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIELuv:1/1/1", "1616/1414/0e0e")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "CIELuv:0.5/0.5/0.5", "0e0e/1313/0e0e")
}

// From esctest2/esctest/tests/change_special_color.py::test_ChangeSpecialColor_TekHVC
func TestEsctestChangeSpecialColorTestChangeSpecialColorTekHVC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "TekHVC:1/1/1", "1a1a/1313/0f0f")
	esctestDoChangeSpecialColorTest(t, screen, stream, "0", "TekHVC:0.5/0.5/0.5", "1111/1313/0e0e")
}
