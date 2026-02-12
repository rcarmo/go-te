package te

import "testing"

func TestEsctest2ChangeSpecialColorRGB(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlOSC + "5;0;rgb:f0f0/f0f0/f0f0" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "5;0;?" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlOSC+"5;0;rgb:f0f0/f0f0/f0f0"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2ResetSpecialColor(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlOSC + "5;0;#fff" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "105;0" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "5;0;?" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlOSC+"5;0;rgb:0000/0000/0000"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
