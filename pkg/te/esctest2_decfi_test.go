package te

import "testing"

// From esctest2/esctest/tests/decfi.py::test_DECFI_Basic
func TestEsctestDecfiTestDECFIBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECFI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 6})
}

// From esctest2/esctest/tests/decfi.py::test_DECFI_NoWrapOnRightEdge
func TestEsctestDecfiTestDECFINoWrapOnRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width, Y: 2})
	esctestDECFI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: size.Width, Y: 2})
}

// From esctest2/esctest/tests/decfi.py::test_DECFI_Scrolls
func TestEsctestDecfiTestDECFIScrolls(t *testing.T) {
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
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestDECFI(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 3, Right: 6, Bottom: 7}, []string{"abcde", "fhi" + esctestEmpty() + "j", "kmn" + esctestEmpty() + "o", "prs" + esctestEmpty() + "t", "uvwxy"})
}

// From esctest2/esctest/tests/decfi.py::test_DECFI_RightOfMargin
func TestEsctestDecfiTestDECFIRightOfMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 5)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 1})
	esctestDECFI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 7, Y: 1})
}

// From esctest2/esctest/tests/decfi.py::test_DECFI_WholeScreenScrolls
func TestEsctestDecfiTestDECFIWholeScreenScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width, Y: 1})
	esctestWrite(t, stream, "x")
	esctestDECFI(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width - 1, Top: 1, Right: size.Width, Bottom: 1}, []string{"x" + esctestEmpty()})
}
