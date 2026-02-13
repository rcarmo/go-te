package te

import "testing"

// From esctest2/esctest/tests/sos.py::test_SOS_Basic
func TestEsctestSosTestSOSBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWriteSOS(t, stream, "xyz")
	esctestWrite(t, stream, "A")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{"A" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/sos.py::test_SOS_8bit
func TestEsctestSosTestSOS8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
