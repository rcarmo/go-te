package te

import (
	"strings"
	"testing"
)

func esctestGetVTLevel(t *testing.T, screen *Screen) int {
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceAttributes(0, false, '>', 0)
	})
	params := esctestReadCSI(t, response, 'c', '>')
	if len(params) == 0 {
		return 0
	}
	level := params[0]
	if level <= 24 {
		if level < 18 {
			return 2
		}
		return 3
	}
	return 4
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DECXCPR
func TestEsctestDecdsrTestDecdsrDecxcpr(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	level := esctestGetVTLevel(t, screen)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(6, true, '?')
	})
	params := esctestReadCSI(t, response, 'R', '?')
	if level >= 4 {
		esctestAssertEQ(t, params, []int{6, 5, 1})
	} else {
		esctestAssertEQ(t, params, []int{6, 5})
	}
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRPrinterPort
func TestEsctestDecdsrTestDecdsrDsrPrinterPort(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(15, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	valid := map[int]bool{10: true, 11: true, 13: true, 18: true, 19: true}
	if !valid[params[0]] {
		t.Fatalf("unexpected printer status %d", params[0])
	}
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRUDKLocked
func TestEsctestDecdsrTestDecdsrDsrUDKLocked(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(25, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	valid := map[int]bool{20: true, 21: true}
	if !valid[params[0]] {
		t.Fatalf("unexpected UDK status %d", params[0])
	}
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRKeyboard
func TestEsctestDecdsrTestDecdsrDsrKeyboard(t *testing.T) {
	screen := NewScreen(80, 24)
	level := esctestGetVTLevel(t, screen)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(26, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if level <= 2 {
		if len(params) != 2 {
			t.Fatalf("expected 2 params, got %d", len(params))
		}
	} else if level == 3 {
		if len(params) != 3 {
			t.Fatalf("expected 3 params, got %d", len(params))
		}
	} else {
		if len(params) != 4 {
			t.Fatalf("expected 4 params, got %d", len(params))
		}
	}
	if params[0] != 27 {
		t.Fatalf("expected 27, got %d", params[0])
	}
	if len(params) > 1 {
		valid := map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true, 22: true, 28: true, 29: true, 30: true, 31: true, 33: true, 35: true, 36: true, 38: true, 39: true, 40: true}
		if !valid[params[1]] {
			t.Fatalf("unexpected keyboard language %d", params[1])
		}
	}
	if len(params) > 2 {
		valid := map[int]bool{0: true, 3: true, 8: true}
		if !valid[params[2]] {
			t.Fatalf("unexpected keyboard status %d", params[2])
		}
	}
	if len(params) > 3 {
		valid := map[int]bool{0: true, 1: true, 4: true, 5: true}
		if !valid[params[3]] {
			t.Fatalf("unexpected keyboard type %d", params[3])
		}
	}
}

func esctestDoLocatorStatusTest(t *testing.T, screen *Screen, code int) {
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(code, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	valid := map[int]bool{50: true, 53: true, 55: true}
	if !valid[params[0]] {
		t.Fatalf("unexpected locator status %d", params[0])
	}
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRDECLocatorStatus
func TestEsctestDecdsrTestDecdsrDsrDecLocatorStatus(t *testing.T) {
	screen := NewScreen(80, 24)
	esctestDoLocatorStatusTest(t, screen, 55)
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRXtermLocatorStatus
func TestEsctestDecdsrTestDecdsrDsrXtermLocatorStatus(t *testing.T) {
	screen := NewScreen(80, 24)
	esctestDoLocatorStatusTest(t, screen, 55)
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_LocatorType
func TestEsctestDecdsrTestDecdsrLocatorType(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(56, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if params[0] != 57 {
		t.Fatalf("expected 57, got %d", params[0])
	}
	valid := map[int]bool{0: true, 1: true, 2: true}
	if !valid[params[1]] {
		t.Fatalf("unexpected locator type %d", params[1])
	}
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DECMSR
func TestEsctestDecdsrTestDecdsrDecmsr(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(62, true, '?')
	})
	params := esctestReadCSI(t, response, '{', 0)
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	esctestAssertEQ(t, params[0], 0)
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DECCKSR
func TestEsctestDecdsrTestDecdsrDeccksr(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(63, true, '?', 123)
	})
	payload := strings.TrimPrefix(response, ControlESC+"P")
	payload = strings.TrimSuffix(payload, ControlST)
	esctestAssertEQ(t, payload, "123!~0000")
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRDataIntegrity
func TestEsctestDecdsrTestDecdsrDsrDataIntegrity(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(75, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	esctestAssertEQ(t, params[0], 70)
}

// From esctest2/esctest/tests/decdsr.py::test_DECDSR_DSRMultipleSessionStatus
func TestEsctestDecdsrTestDecdsrDsrMultipleSessionStatus(t *testing.T) {
	screen := NewScreen(80, 24)
	response := esctestCaptureResponse(screen, func() {
		screen.ReportDeviceStatus(85, true, '?')
	})
	params := esctestReadCSI(t, response, 'n', '?')
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	valid := map[int]bool{83: true, 87: true}
	if !valid[params[0]] {
		t.Fatalf("unexpected session status %d", params[0])
	}
}
