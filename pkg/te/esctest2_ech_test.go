package te

import "testing"

// From esctest2/esctest/tests/ech.py::test_ECH_DefaultParam
func TestEsctestEchTestECHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + "bc"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_ExplicitParam
func TestEsctestEchTestECHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_IgnoresScrollRegion
func TestEsctestEchTestECHIgnoresScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcdefg")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 1})
	esctestECH(t, stream, 4)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 1}, []string{"ab" + esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + "g"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_OutsideScrollRegion
func TestEsctestEchTestECHOutsideScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcdefg")
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 4)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 7, Bottom: 1}, []string{esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + "efg"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_doesNotRespectDECPRotection
func TestEsctestEchTestECHDoesNotRespectDECPRotection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlCSI+"1\"q")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlCSI+"0\"q")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + esctestBlank()})
}

// From esctest2/esctest/tests/ech.py::test_ECH_respectsISOProtection
func TestEsctestEchTestECHRespectsISOProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlESC+"W")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}
