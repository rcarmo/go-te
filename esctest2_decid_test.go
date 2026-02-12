package te

import "testing"

func TestEsctest2DECIDBasic(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlESC + "Z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2DECID8Bit(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed("\x1b G\u009a\x1b F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
