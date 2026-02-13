package te

import "testing"

// From esctest2/esctest/tests/decid.py::test_DECID_Basic
func TestEsctestDecidTestDECIDBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlESC+"Z")
	})
	params := esctestParseCSI(t, response, '?')
	if len(params) == 0 {
		t.Fatalf("expected response params")
	}
}

// From esctest2/esctest/tests/decid.py::test_DECID_8bit
func TestEsctestDecidTestDECID8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
