package te

import "testing"

// From esctest2/esctest/tests/pm.py::test_PM_Basic
func TestEsctestPmTestPMBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWritePM(t, stream, "xyz")
	esctestWrite(t, stream, "A")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{"A" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/pm.py::test_PM_8bit
func TestEsctestPmTestPM8bit(t *testing.T) {
	t.Skip("requires DISABLE_WIDE_CHARS / ALLOW_C2_CONTROLS options")
}
