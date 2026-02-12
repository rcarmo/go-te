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

func TestEsctest2XtermWinopsReportWindowState(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "11t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"1t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportWindowPosition(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "13t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"3;0;0t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportWindowSizePixels(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "14t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"4;2;10t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportScreenSizePixels(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "15t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"5;2;10t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportCharSizePixels(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "16t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"6;0;0t" {
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

func TestEsctest2XtermWinopsReportScreenSizeChars(t *testing.T) {
	screen := NewScreen(12, 4)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "19t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"9;4;12t" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2XtermWinopsReportIconTitle(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	screen.IconName = "icon"
	if err := stream.Feed(ControlCSI + "20t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlOSC+"L"+screen.IconName+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
