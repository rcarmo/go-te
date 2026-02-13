package te

import "testing"

func esctestExpectedDA2First(screen *Screen) int {
	switch screen.conformanceLevel {
	case 5:
		return 64
	case 4:
		return 41
	case 3:
		return 24
	case 2:
		return 1
	default:
		return 0
	}
}

func esctestAssertDA2Response(t *testing.T, params []int, screen *Screen) {
	if len(params) != 3 {
		t.Fatalf("expected 3 params, got %d", len(params))
	}
	if params[0] != esctestExpectedDA2First(screen) {
		t.Fatalf("expected first param %d, got %d", esctestExpectedDA2First(screen), params[0])
	}
	if params[1] < 314 || params[1] > 999 {
		t.Fatalf("expected second param in range, got %d", params[1])
	}
}

// From esctest2/esctest/tests/da2.py::test_DA2_NoParameter
func TestEsctestDa2TestDA2NoParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetConformance(64, 1)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlCSI+">"+EscDA)
	})
	params := esctestParseCSI(t, response, '>')
	esctestAssertDA2Response(t, params, screen)
}

// From esctest2/esctest/tests/da2.py::test_DA2_0
func TestEsctestDa2TestDA20(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetConformance(64, 1)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlCSI+">0"+EscDA)
	})
	params := esctestParseCSI(t, response, '>')
	esctestAssertDA2Response(t, params, screen)
}
