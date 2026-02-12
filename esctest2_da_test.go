package te

import "testing"

func TestEsctest2DA(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + "c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlCSI + "0c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+"?6c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}

func TestEsctest2DA2(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	responses := []string{}
	screen.WriteProcessInput = func(data string) { responses = append(responses, data) }
	if err := stream.Feed(ControlCSI + ">c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+">0;0;0c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
	responses = nil
	if err := stream.Feed(ControlCSI + ">0c"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(responses) == 0 || responses[len(responses)-1] != ControlCSI+">0;0;0c" {
		t.Fatalf("unexpected response: %#v", responses)
	}
}
