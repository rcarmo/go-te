package te

import "testing"

// From esctest2/esctest/tests/sm_title.py::test_SMTitle_SetHexQueryUTF8
func TestEsctestSmTitleTestSmTitleSetHexQueryUtf8(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryHex)
	esctestSMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryUTF8)
	esctestChangeWindowTitle(t, stream, "6162")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "ab")
	esctestChangeWindowTitle(t, stream, "61")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "a")
	esctestChangeIconTitle(t, stream, "6162")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "ab")
	esctestChangeIconTitle(t, stream, "61")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "a")
}

// From esctest2/esctest/tests/sm_title.py::test_SMTitle_SetUTF8QueryUTF8
func TestEsctestSmTitleTestSmTitleSetUtf8QueryUtf8(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryHex)
	esctestSMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryUTF8)
	esctestChangeWindowTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "ab")
	esctestChangeWindowTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "a")
	esctestChangeIconTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "ab")
	esctestChangeIconTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "a")
}

// From esctest2/esctest/tests/sm_title.py::test_SMTitle_SetUTF8QueryHex
func TestEsctestSmTitleTestSmTitleSetUtf8QueryHex(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryUTF8)
	esctestSMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryHex)
	esctestChangeWindowTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "6162")
	esctestChangeWindowTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "61")
	esctestChangeIconTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "6162")
	esctestChangeIconTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "61")
}

// From esctest2/esctest/tests/sm_title.py::test_SMTitle_SetHexQueryHex
func TestEsctestSmTitleTestSmTitleSetHexQueryHex(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryUTF8)
	esctestSMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryHex)
	esctestChangeWindowTitle(t, stream, "6162")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "6162")
	esctestChangeWindowTitle(t, stream, "61")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "61")
	esctestChangeIconTitle(t, stream, "6162")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "6162")
	esctestChangeIconTitle(t, stream, "61")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "61")
}
