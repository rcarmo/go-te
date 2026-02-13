package te

import "testing"

// From esctest2/esctest/tests/el.py::test_EL_Default
func TestEsctestElTestELDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_0
func TestEsctestElTestEL0(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 0)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_1
func TestEsctestElTestEL1(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + "fghij"})
}

// From esctest2/esctest/tests/el.py::test_EL_2
func TestEsctestElTestEL2(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_IgnoresScrollRegion
func TestEsctestElTestELIgnoresScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestEL(t, stream, 2)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_doesNotRespectDECProtection
func TestEsctestElTestELDoesNotRespectDECProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlCSI+"1\"q")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlCSI+"0\"q")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_respectsISOProtection
func TestEsctestElTestELRespectsISOProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlESC+"W")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}
