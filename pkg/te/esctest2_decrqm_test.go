package te

import "testing"

type esctestDecrqmFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecrqmFixture() esctestDecrqmFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecrqmFixture{screen: screen, stream: stream}
}

func (f esctestDecrqmFixture) requestAnsiMode(t *testing.T, mode int) []int {
	response := esctestCaptureResponse(f.screen, func() {
		esctestDECRQM(t, f.stream, mode, false)
	})
	return esctestReadCSI(t, response, "$y", 0)
}

func (f esctestDecrqmFixture) requestDecMode(t *testing.T, mode int) []int {
	response := esctestCaptureResponse(f.screen, func() {
		esctestDECRQM(t, f.stream, mode, true)
	})
	return esctestReadCSI(t, response, "$y", '?')
}

func (f esctestDecrqmFixture) doModifiableAnsiTest(t *testing.T, mode int) {
	before := f.requestAnsiMode(t, mode)
	if len(before) < 2 {
		t.Fatalf("expected 2 params, got %v", before)
	}
	if before[1] == 2 {
		esctestSM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 1})
		esctestRM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 2})
	} else {
		esctestRM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 2})
		esctestSM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 1})
	}
}

func (f esctestDecrqmFixture) doPermanentlyResetAnsiTest(t *testing.T, mode int) {
	esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 4})
}

func (f esctestDecrqmFixture) doModifiableDecTest(t *testing.T, mode int) {
	before := f.requestDecMode(t, mode)
	if len(before) < 2 {
		t.Fatalf("expected 2 params, got %v", before)
	}
	if before[1] == 2 {
		esctestDECSET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 1})
		esctestDECRESET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 2})
	} else {
		esctestDECRESET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 2})
		esctestDECSET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 1})
	}
}

func (f esctestDecrqmFixture) doPermanentlyResetDecTest(t *testing.T, mode int) {
	esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 4})
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM
func TestEsctestDecrqmTestDecrqm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	esctestAssertEQ(t, len(fixture.requestAnsiMode(t, esctestModeIRM)), 2)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_KAM
func TestEsctestDecrqmTestDecrqmAnsiKam(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, 2)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_IRM
func TestEsctestDecrqmTestDecrqmAnsiIrm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, esctestModeIRM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SRM
func TestEsctestDecrqmTestDecrqmAnsiSrm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, 12)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_LNM
func TestEsctestDecrqmTestDecrqmAnsiLnm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, esctestModeLNM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_GATM
func TestEsctestDecrqmTestDecrqmAnsiGatm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 1)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SRTM
func TestEsctestDecrqmTestDecrqmAnsiSrtm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_VEM
func TestEsctestDecrqmTestDecrqmAnsiVem(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 7)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_HEM
func TestEsctestDecrqmTestDecrqmAnsiHem(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 10)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_PUM
func TestEsctestDecrqmTestDecrqmAnsiPum(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 11)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_FEAM
func TestEsctestDecrqmTestDecrqmAnsiFeam(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 13)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_FETM
func TestEsctestDecrqmTestDecrqmAnsiFetm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 14)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_MATM
func TestEsctestDecrqmTestDecrqmAnsiMatm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 15)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_TTM
func TestEsctestDecrqmTestDecrqmAnsiTtm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 16)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SATM
func TestEsctestDecrqmTestDecrqmAnsiSatm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 17)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_TSM
func TestEsctestDecrqmTestDecrqmAnsiTsm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 18)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_EBM
func TestEsctestDecrqmTestDecrqmAnsiEbm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 19)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCKM
func TestEsctestDecrqmTestDecrqmDecDecckm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECCKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCOLM
func TestEsctestDecrqmTestDecrqmDecDeccolm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	fixture.doModifiableDecTest(t, ModeDECCOLM>>5)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECSCLM
func TestEsctestDecrqmTestDecrqmDecDecsclm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECSCLM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECSCNM
func TestEsctestDecrqmTestDecrqmDecDecscnm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, ModeDECSCNM>>5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECOM
func TestEsctestDecrqmTestDecrqmDecDecom(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECOM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECAWM
func TestEsctestDecrqmTestDecrqmDecDecawm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECAWM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECARM
func TestEsctestDecrqmTestDecrqmDecDecarm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECARM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPFF
func TestEsctestDecrqmTestDecrqmDecDecpff(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECPFF)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPEX
func TestEsctestDecrqmTestDecrqmDecDecpex(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECPEX)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECTCEM
func TestEsctestDecrqmTestDecrqmDecDectcem(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, ModeDECTCEM>>5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECRLM
func TestEsctestDecrqmTestDecrqmDecDecrlm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECRLM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHEBM
func TestEsctestDecrqmTestDecrqmDecDechebm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 35)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHEM
func TestEsctestDecrqmTestDecrqmDecDechem(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 36)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNRCM
func TestEsctestDecrqmTestDecrqmDecDecnrcm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNRCM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNAKB
func TestEsctestDecrqmTestDecrqmDecDecnakb(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 57)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECVCCM
func TestEsctestDecrqmTestDecrqmDecDecvccm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 61)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPCCM
func TestEsctestDecrqmTestDecrqmDecDecpccm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 64)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNKM
func TestEsctestDecrqmTestDecrqmDecDecnkm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECBKM
func TestEsctestDecrqmTestDecrqmDecDecbkm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECBKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECKBUM
func TestEsctestDecrqmTestDecrqmDecDeckbum(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECKBUM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECLRMM
func TestEsctestDecrqmTestDecrqmDecDeclrmm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECLRMM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECXRLM
func TestEsctestDecrqmTestDecrqmDecDecxrlm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 73)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECKPM
func TestEsctestDecrqmTestDecrqmDecDeckpm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 81)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNCSM
func TestEsctestDecrqmTestDecrqmDecDecncsm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	fixture.doModifiableDecTest(t, esctestModeDECNCSM)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECRLCM
func TestEsctestDecrqmTestDecrqmDecDecrlcm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 96)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCRTSM
func TestEsctestDecrqmTestDecrqmDecDeccrtsm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 97)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECARSM
func TestEsctestDecrqmTestDecrqmDecDecarsm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 98)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECMCM
func TestEsctestDecrqmTestDecrqmDecDecmcm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 99)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECAAM
func TestEsctestDecrqmTestDecrqmDecDecaam(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECAAM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCANSM
func TestEsctestDecrqmTestDecrqmDecDeccansm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECCANSM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNULM
func TestEsctestDecrqmTestDecrqmDecDecnulm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNULM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHDPXM
func TestEsctestDecrqmTestDecrqmDecDechdpxm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECHDPXM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECESKM
func TestEsctestDecrqmTestDecrqmDecDeceskym(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECESKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECOSCNM
func TestEsctestDecrqmTestDecrqmDecDecoscnm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECOSCNM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHCCM
func TestEsctestDecrqmTestDecrqmDecDechccm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetDecTest(t, esctestModeDECHCCM)
}
