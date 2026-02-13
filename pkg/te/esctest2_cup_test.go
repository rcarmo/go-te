package te

import "testing"

// From esctest2/esctest/tests/cup.py::test_CUP_DefaultParams
func TestEsctestCupTestCUPDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	esctestCUPParams(t, stream, nil, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_RowOnly
func TestEsctestCupTestCUPRowOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	row := 2
	esctestCUPParams(t, stream, &row, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cup.py::test_CUP_ColumnOnly
func TestEsctestCupTestCUPColumnOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	col := 2
	esctestCUPParams(t, stream, nil, &col)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_ZeroIsTreatedAsOne
func TestEsctestCupTestCUPZeroIsTreatedAsOne(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	row := 0
	col := 0
	esctestCUPParams(t, stream, &row, &col)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_OutOfBoundsParams
func TestEsctestCupTestCUPOutOfBoundsParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width + 10, Y: size.Height + 10})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/cup.py::test_CUP_RespectsOriginMode
func TestEsctestCupTestCUPRespectsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestCUP(t, stream, esctestPoint{X: 7, Y: 9})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
	esctestAssertEQ(t, pos.Y, 9)

	esctestDECSET(t, stream, esctestModeDECOM)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
	esctestWrite(t, stream, "X")

	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 6, Right: 5, Bottom: 6}, []string{"X"})
}
