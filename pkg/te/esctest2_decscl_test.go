package te

import "testing"

// From esctest2/esctest/tests/decscl.py::test_DECSCL_Level2DoesntSupportDECRQM
func TestEsctestDecsclTestDecsclLevel2DoesntSupportDecrqm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "Hello world.")
	esctestGetScreenSize(screen)
	esctestDECSCL(t, stream, 62, esctestIntPtr(1))
	esctestGetScreenSize(screen)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQM(t, stream, esctestModeIRM, false)
	})
	if response != "" {
		t.Fatalf("expected DECRQM to be unsupported at level 2")
	}
}

// From esctest2/esctest/tests/decscl.py::test_DSCSCL_Level2Supports7BitControls
func TestEsctestDecsclTestDecsclLevel2Supports7BitControls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCL(t, stream, 62, esctestIntPtr(1))
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 2})
}

// From esctest2/esctest/tests/decscl.py::test_DSCSCL_Level3_SupportsDECRQMDoesntSupportDECSLRM
func TestEsctestDecsclTestDecsclLevel3SupportsDecrqmDoesntSupportDecslrm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCL(t, stream, 63, esctestIntPtr(1))
	_ = esctestCaptureResponse(screen, func() {
		esctestDECRQM(t, stream, esctestModeIRM, false)
	})
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 6)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, "abc")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 8)
}

// From esctest2/esctest/tests/decscl.py::test_DECSCL_Level4_SupportsDECSLRMDoesntSupportDECNCSM
func TestEsctestDecsclTestDecsclLevel4SupportsDecslrmDoesntSupportDecncsm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCL(t, stream, 64, esctestIntPtr(1))
	esctestDECSET(t, stream, esctestModeAllow80To132)
	esctestDECRESET(t, stream, ModeDECCOLM>>5)
	esctestDECSET(t, stream, esctestModeDECNCSM)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "1")
	esctestDECSET(t, stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 5, 6)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, "abc")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 6)
}

// From esctest2/esctest/tests/decscl.py::test_DECSCL_Level5_SupportsDECNCSM
func TestEsctestDecsclTestDecsclLevel5SupportsDecncsm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCL(t, stream, 65, esctestIntPtr(1))
	esctestDECRESET(t, stream, ModeDECCOLM>>5)
	esctestDECSET(t, stream, esctestModeDECNCSM)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "1")
	esctestDECSET(t, stream, ModeDECCOLM>>5)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"1"})
}

// From esctest2/esctest/tests/decscl.py::test_DECSCL_RISOnChange
func TestEsctestDecsclTestDecsclRisOnChange(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSCL(t, stream, 61, nil)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"b"})
}
