package te

import (
	"encoding/base64"
	"testing"
)

// From esctest2/esctest/tests/manipulate_selection_data.py::test_ManipulateSelectionData_default
func TestEsctestManipulateSelectionDataDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	encoded := base64.StdEncoding.EncodeToString([]byte("testing 123"))
	esctestWrite(t, stream, ControlOSC+"52;s0;"+encoded+ControlST)
	response := esctestCaptureResponse(screen, func() {
		esctestWrite(t, stream, ControlOSC+"52;?"+ControlST)
	})
	value := esctestReadOSC(t, response, "52")
	esctestAssertEQ(t, value, ";s0;"+encoded)
}
