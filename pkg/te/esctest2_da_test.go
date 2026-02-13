package te

import "testing"

func esctestExpectedDAParams(screen *Screen) []int {
	switch screen.conformanceLevel {
	case 5:
		return []int{65, 1, 2, 6, 9, 15, 16, 17, 18, 21, 22, 28, 29}
	case 4:
		return []int{64, 1, 2, 6, 9, 15, 16, 17, 18, 21, 22, 28, 29}
	case 3:
		return []int{63, 1, 2, 6, 9, 15, 22, 29}
	case 2:
		return []int{62, 1, 2, 6, 9, 15, 22, 29}
	case 1:
		return []int{1, 2}
	default:
		return []int{0}
	}
}

func esctestAssertDAResponse(t *testing.T, params []int, expected []int) {
	if len(params) < len(expected) {
		t.Fatalf("expected at least %d params, got %d", len(expected), len(params))
	}
	if params[0] != expected[0] {
		t.Fatalf("expected first param %d, got %d", expected[0], params[0])
	}
	for _, value := range expected {
		found := false
		for _, param := range params {
			if param == value {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected param %d in response", value)
		}
	}
}

// From esctest2/esctest/tests/da.py::test_DA_NoParameter
func TestEsctestDaTestDANoParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetConformance(64, 1)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlCSI+EscDA)
	})
	params := esctestParseCSI(t, response, '?')
	esctestAssertDAResponse(t, params, esctestExpectedDAParams(screen))
}

// From esctest2/esctest/tests/da.py::test_DA_0
func TestEsctestDaTestDA0(t *testing.T) {
	screen := NewScreen(80, 24)
	screen.SetConformance(64, 1)
	stream := NewStream(screen, false)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlCSI+"0"+EscDA)
	})
	params := esctestParseCSI(t, response, '?')
	esctestAssertDAResponse(t, params, esctestExpectedDAParams(screen))
}
