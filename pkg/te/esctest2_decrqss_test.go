package te

import (
	"strings"
	"testing"
)

func esctestReadDCSResponse(response string) string {
	payload := strings.TrimPrefix(response, ControlDCS)
	payload = strings.TrimSuffix(payload, ControlST)
	return payload
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECELF
func TestEsctestDecrqssTestDecrqssDecelf(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECELF(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "+q")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0+q")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECLFKC
func TestEsctestDecrqssTestDecrqssDeclfkc(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECLFKC(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "*}")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0*}")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSASD
func TestEsctestDecrqssTestDecrqssDecsasd(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSASD(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "$}")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0$}")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSACE
func TestEsctestDecrqssTestDecrqssDecsace(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSACE(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "*x")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0*x")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSCA
func TestEsctestDecrqssTestDecrqssDecsca(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "\"q")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r1\"q")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSCL
func TestEsctestDecrqssTestDecrqssDecscl(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	sevenBit := 1
	esctestDECSCL(t, stream, 65, &sevenBit)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "\"p")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r65;1\"p")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSTBM
func TestEsctestDecrqssTestDecrqssDecstbm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 5, 6)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "r")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r5;6r")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_SGR
func TestEsctestDecrqssTestDecrqssSgr(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSGR(t, stream, 1)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "m")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0;1m")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSCUSR
func TestEsctestDecrqssTestDecrqssDecscusr(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCUSR(t, stream, 4)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, " q")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r4 q")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSLRM
func TestEsctestDecrqssTestDecrqssDecslrm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECLRMM)
	esctestDECSLRM(t, stream, 3, 4)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "s")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r3;4s")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSLPP
func TestEsctestDecrqssTestDecrqssDecslpp(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestXtermWinops(t, stream, 27)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "t")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r27t")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSMKR
func TestEsctestDecrqssTestDecrqssDecsmkr(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSMKR(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "+r")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0+r")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSNLS
func TestEsctestDecrqssTestDecrqssDecsnls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSNLS(t, stream, 24)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "*|")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r24*|")
}

// From esctest2/esctest/tests/decrqss.py::test_DECRQSS_DECSSDT
func TestEsctestDecrqssTestDecrqssDecssdt(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSSDT(t, stream, 0)
	response := esctestCaptureResponse(screen, func() {
		esctestDECRQSS(t, stream, "$~")
	})
	esctestAssertEQ(t, esctestReadDCSResponse(response), "1$r0$~")
}
