package te

import "testing"

// From esctest2/esctest/tests/sgr.py::test_SGR_Bold
func TestEsctestSgrTestSgrBold(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "x")
	esctestSGR(t, stream, 1)
	esctestWrite(t, stream, "y")
	esctestAssertCharHasSGR(t, screen, esctestPoint{X: 1, Y: 1}, []int{1}, []int{39, 49})
	esctestAssertCharHasSGR(t, screen, esctestPoint{X: 2, Y: 1}, []int{}, []int{39, 49, 1})
}
