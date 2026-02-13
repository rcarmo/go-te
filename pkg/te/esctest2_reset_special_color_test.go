package te

import "testing"

// From esctest2/esctest/tests/reset_special_color.py::test_ResetSpecialColor_Single
func TestEsctestResetSpecialColorTestResetSpecialColorSingle(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "0", "?")
	})
	original := esctestReadOSC(t, response, "4")
	esctestChangeSpecialColor(t, stream, "0", "#aaaabbbbcccc")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != ";16;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetSpecialColor(t, stream, "0")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}

// From esctest2/esctest/tests/reset_special_color.py::test_ResetSpecialColor_Single2
func TestEsctestResetSpecialColorTestResetSpecialColorSingle2(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor2(t, stream, "0", "?")
	})
	original := esctestReadOSC(t, response, "5")
	esctestChangeSpecialColor2(t, stream, "0", "#aaaabbbbcccc")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor2(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "5"); got != ";0;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetSpecialColor(t, stream, "0")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor2(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "5"); got != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}

// From esctest2/esctest/tests/reset_special_color.py::test_ResetSpecialColor_Multiple
func TestEsctestResetSpecialColorTestResetSpecialColorMultiple(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor(t, stream, "0", "?", "1", "?")
	if len(responses) < 2 {
		t.Fatalf("expected responses")
	}
	original1 := esctestReadOSC(t, responses[0], "4")
	original2 := esctestReadOSC(t, responses[1], "4")
	esctestChangeSpecialColor(t, stream, "0", "#aaaabbbbcccc")
	esctestChangeSpecialColor(t, stream, "1", "#ddddeeeeffff")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor(t, stream, "0", "?")
	esctestChangeSpecialColor(t, stream, "1", "?")
	if got := esctestReadOSC(t, responses[0], "4"); got != ";16;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "4"); got != ";17;rgb:dddd/eeee/ffff" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetSpecialColor(t, stream, "0", "1")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor(t, stream, "0", "?", "1", "?")
	if esctestReadOSC(t, responses[0], "4") != original1 {
		t.Fatalf("expected %q", original1)
	}
	if esctestReadOSC(t, responses[1], "4") != original2 {
		t.Fatalf("expected %q", original2)
	}
}

// From esctest2/esctest/tests/reset_special_color.py::test_ResetSpecialColor_Multiple2
func TestEsctestResetSpecialColorTestResetSpecialColorMultiple2(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor2(t, stream, "0", "?", "1", "?")
	if len(responses) < 2 {
		t.Fatalf("expected responses")
	}
	original1 := esctestReadOSC(t, responses[0], "5")
	original2 := esctestReadOSC(t, responses[1], "5")
	esctestChangeSpecialColor2(t, stream, "0", "#aaaabbbbcccc")
	esctestChangeSpecialColor2(t, stream, "1", "#ddddeeeeffff")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor2(t, stream, "0", "?")
	esctestChangeSpecialColor2(t, stream, "1", "?")
	if got := esctestReadOSC(t, responses[0], "5"); got != ";0;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	if got := esctestReadOSC(t, responses[1], "5"); got != ";1;rgb:dddd/eeee/ffff" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetSpecialColor(t, stream, "0", "1")
	responses = []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	esctestChangeSpecialColor2(t, stream, "0", "?", "1", "?")
	if esctestReadOSC(t, responses[0], "5") != original1 {
		t.Fatalf("expected %q", original1)
	}
	if esctestReadOSC(t, responses[1], "5") != original2 {
		t.Fatalf("expected %q", original2)
	}
}

// From esctest2/esctest/tests/reset_special_color.py::test_ResetSpecialColor_Dynamic
func TestEsctestResetSpecialColorTestResetSpecialColorDynamic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "10", "?")
	})
	original := esctestReadOSC(t, response, "10")
	esctestChangeSpecialColor(t, stream, "10", "#aaaabbbbcccc")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "10", "?")
	})
	if got := esctestReadOSC(t, response, "10"); got != ";rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetDynamicColor(t, stream, "110")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeSpecialColor(t, stream, "10", "?")
	})
	if got := esctestReadOSC(t, response, "10"); got != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}
