package te

import "testing"

func TestEsctest2XtermWinopsReportTitle(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	screen.Title = "hello"
	if err := stream.Feed(ControlCSI + "21t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlOSC+"l"+screen.Title+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportTextAreaChars(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "18t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"8;2;10t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
