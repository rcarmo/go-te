package te

import "testing"

// From esctest2/esctest/tests/reset_color.py::test_ResetColor_Standard
func TestEsctestResetColorTestResetColorStandard(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "0", "?")
	})
	original := esctestReadOSC(t, response, "4")
	esctestChangeColor(t, stream, "0", "#aaaabbbbcccc")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != ";0;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetColor(t, stream, "0")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "0", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}

// From esctest2/esctest/tests/reset_color.py::test_ResetColor_All
func TestEsctestResetColorTestResetColorAll(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "3", "?")
	})
	original := esctestReadOSC(t, response, "4")
	esctestChangeColor(t, stream, "3", "#aabbcc")
	response = esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "3", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != ";3;rgb:aaaa/bbbb/cccc" {
		t.Fatalf("unexpected color %q", got)
	}
	esctestResetColor(t, stream)
	response = esctestCaptureResponse(screen, func() {
		esctestChangeColor(t, stream, "3", "?")
	})
	if got := esctestReadOSC(t, response, "4"); got != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}
