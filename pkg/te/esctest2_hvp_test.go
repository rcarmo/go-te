package te

import "testing"

// From esctest2/esctest/tests/hvp.py::test_HVP_DefaultParams
func TestEsctestHvpTestHVPDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	esctestHVPParams(t, stream, nil, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_RowOnly
func TestEsctestHvpTestHVPRowOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	row := 2
	esctestHVPParams(t, stream, &row, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_ColumnOnly
func TestEsctestHvpTestHVPColumnOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	col := 2
	esctestHVPParams(t, stream, nil, &col)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_ZeroIsTreatedAsOne
func TestEsctestHvpTestHVPZeroIsTreatedAsOne(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	row := 0
	col := 0
	esctestHVPParams(t, stream, &row, &col)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_OutOfBoundsParams
func TestEsctestHvpTestHVPOutOfBoundsParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestHVPParams(t, stream, intPtr(size.Height+10), intPtr(size.Width+10))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_RespectsOriginMode
func TestEsctestHvpTestHVPRespectsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 6, 11)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 10)
	esctestHVPParams(t, stream, intPtr(9), intPtr(7))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
	esctestAssertEQ(t, pos.Y, 9)

	esctestDECSET(t, stream, esctestModeDECOM)
	esctestHVPParams(t, stream, intPtr(1), intPtr(1))
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
	esctestWrite(t, stream, "X")

	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 5, Top: 6, Right: 5, Bottom: 6}, []string{"X"})
}

func intPtr(value int) *int {
	return &value
}
