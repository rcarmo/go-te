package te

import "testing"

func esctestDeccraPrepare(t *testing.T, stream *Stream) {
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefgh"+ControlCR+ControlLF)
	esctestWrite(t, stream, "ijklmnop"+ControlCR+ControlLF)
	esctestWrite(t, stream, "qrstuvwx"+ControlCR+ControlLF)
	esctestWrite(t, stream, "yz012345"+ControlCR+ControlLF)
	esctestWrite(t, stream, "ABCDEFGH"+ControlCR+ControlLF)
	esctestWrite(t, stream, "IJKLMNOP"+ControlCR+ControlLF)
	esctestWrite(t, stream, "QRSTUVWX"+ControlCR+ControlLF)
	esctestWrite(t, stream, "YZ6789!@"+ControlCR+ControlLF)
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_nonOverlappingSourceAndDest
func TestEsctestDeccraTestDeccraNonOverlappingSourceAndDest(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), esctestIntPtr(5), esctestIntPtr(5), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDjklH",
		"IJKLrstP",
		"QRSTz01X",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_overlappingSourceAndDest
func TestEsctestDeccraTestDeccraOverlappingSourceAndDest(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), esctestIntPtr(3), esctestIntPtr(3), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrjklvwx",
		"yzrst345",
		"ABz01FGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_destinationPartiallyOffscreen
func TestEsctestDeccraTestDeccraDestinationPartiallyOffscreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	size := esctestGetScreenSize(screen)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), esctestIntPtr(size.Height-1), esctestIntPtr(size.Width-1), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width - 1, Top: size.Height - 1, Right: size.Width, Bottom: size.Height}, []string{
		"jk",
		"rs",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_defaultValuesInSource
func TestEsctestDeccraTestDeccraDefaultValuesInSource(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECCRA(t, stream, nil, nil, esctestIntPtr(2), esctestIntPtr(2), nil, esctestIntPtr(5), esctestIntPtr(5), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDabGH",
		"IJKLijOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_defaultValuesInDest
func TestEsctestDeccraTestDeccraDefaultValuesInDest(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), nil, nil, nil)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"jkldefgh",
		"rstlmnop",
		"z01tuvwx",
		"yz012345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_invalidSourceRectDoesNothing
func TestEsctestDeccraTestDeccraInvalidSourceRectDoesNothing(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(1), esctestIntPtr(1), esctestIntPtr(1), esctestIntPtr(5), esctestIntPtr(5), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_respectsOriginMode
func TestEsctestDeccraTestDeccraRespectsOriginMode(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 2, 9)
	esctestDECSTBM(t, stream, 2, 9)
	esctestDECSET(t, stream, esctestModeDECOM)
	esctestDECCRA(t, stream, esctestIntPtr(1), esctestIntPtr(1), esctestIntPtr(3), esctestIntPtr(3), esctestIntPtr(1), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1))
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestDECRESET(t, stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDjklH",
		"IJKLrstP",
		"QRSTz01X",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_ignoresMargins
func TestEsctestDeccraTestDeccraIgnoresMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 6)
	esctestDECSTBM(t, stream, 3, 6)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), esctestIntPtr(5), esctestIntPtr(5), esctestIntPtr(1))
	esctestDECRESET(t, stream, esctestModeDECLRMM)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDjklH",
		"IJKLrstP",
		"QRSTz01X",
		"YZ6789!@",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_overlyLargeSourceClippedToScreenSize
func TestEsctestDeccraTestDeccraOverlyLargeSourceClippedToScreenSize(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: size.Height - 1})
	esctestWrite(t, stream, "ab")
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: size.Height})
	esctestWrite(t, stream, "cX")
	esctestDECCRA(t, stream, esctestIntPtr(size.Height), esctestIntPtr(size.Width), esctestIntPtr(size.Height+1), esctestIntPtr(size.Width+1), esctestIntPtr(1), esctestIntPtr(size.Height-1), esctestIntPtr(size.Width-1), esctestIntPtr(1))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width - 1, Top: size.Height - 1, Right: size.Width, Bottom: size.Height}, []string{
		"Xb",
		"cX",
	})
}

// From esctest2/esctest/tests/deccra.py::test_DECCRA_cursorDoesNotMove
func TestEsctestDeccraTestDeccraCursorDoesNotMove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDeccraPrepare(t, stream)
	position := esctestPoint{X: 3, Y: 4}
	esctestCUP(t, stream, position)
	esctestDECCRA(t, stream, esctestIntPtr(2), esctestIntPtr(2), esctestIntPtr(4), esctestIntPtr(4), esctestIntPtr(1), esctestIntPtr(5), esctestIntPtr(5), esctestIntPtr(1))
	esctestAssertEQ(t, esctestGetCursorPosition(screen), position)
}
