package te

import "testing"

func TestEsctest2DECDSRDECXCPR(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	setCursor(screen, 5, 6)
	if err := stream.Feed(ControlCSI + "?6n"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6;5;1R" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2DECRQM(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "4$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"4;2$y" {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlCSI + "4h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "4$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"4;1$y" {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlCSI + "?6$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6;2$y" {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlCSI + "?6h"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "?6$p"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6;1$y" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2DECRQSS(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	screen.SelectGraphicRendition([]int{1}, false)
	if err := stream.Feed(ControlDCS + "$qm" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlDCS+"1$r0;1m"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	screen.SetMargins(2, 4)
	if err := stream.Feed(ControlDCS + "$qr" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlDCS+"1$r2;4r"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	screen.SetMode([]int{69}, true)
	screen.SetLeftRightMargins(3, 5)
	if err := stream.Feed(ControlDCS + "$qs" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlDCS+"1$r3;5s"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlDCS + "$q q" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlDCS+"1$r0 q"+ControlST {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
