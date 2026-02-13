package te

import "testing"

// From esctest2/esctest/tests/decbi.py::test_DECBI_Basic
func TestEsctestDecbiTestDECBIBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECBI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 4, Y: 6})
}

// From esctest2/esctest/tests/decbi.py::test_DECBI_NoWrapOnLeftEdge
func TestEsctestDecbiTestDECBINoWrapOnLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestDECBI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 2})
}

// From esctest2/esctest/tests/decbi.py::test_DECBI_Scrolls
func TestEsctestDecbiTestDECBIScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	y := 3
	for _, line := range lines {
		esctestCUP(t, stream, esctestPoint{X: 2, Y: y})
		esctestWrite(t, stream, line)
		y++
	}
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 5)
	esctestDECSTBM(t, stream, 4, 6)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 5})
	esctestDECBI(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 3, Right: 6, Bottom: 7}, []string{"abcde", "f" + esctestBlank() + "ghj", "k" + esctestBlank() + "lmo", "p" + esctestBlank() + "qrt", "uvwxy"})
}

// From esctest2/esctest/tests/decbi.py::test_DECBI_LeftOfMargin
func TestEsctestDecbiTestDECBILeftOfMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestDECBI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decbi.py::test_DECBI_WholeScreenScrolls
func TestEsctestDecbiTestDECBIWholeScreenScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECBI(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{esctestBlank() + "x"})
}
