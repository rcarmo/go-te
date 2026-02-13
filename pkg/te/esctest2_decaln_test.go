package te

import "testing"

// From esctest2/esctest/tests/decaln.py::test_DECALN_FillsScreen
func TestEsctestDecalnTestDECALNFillsScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECALN(t, stream)
	size := esctestGetScreenSize(screen)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width, Top: 1, Right: size.Width, Bottom: 1}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: size.Height, Right: 1, Bottom: size.Height}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width, Top: size.Height, Right: size.Width, Bottom: size.Height}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width / 2, Top: size.Height / 2, Right: size.Width / 2, Bottom: size.Height / 2}, []string{"E"})
}

// From esctest2/esctest/tests/decaln.py::test_DECALN_MovesCursorHome
func TestEsctestDecalnTestDECALNMovesCursorHome(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestDECALN(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decaln.py::test_DECALN_ClearsMargins
func TestEsctestDecalnTestDECALNClearsMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 3)
	esctestDECSTBM(t, stream, 4, 5)
	esctestDECALN(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestCUU(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestCUD(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 6})
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestCUB(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 4})
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 4})
	esctestCUF(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 4, Y: 4})
}
