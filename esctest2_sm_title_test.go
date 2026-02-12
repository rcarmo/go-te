package te

import "testing"

func TestEsctest2SMTitleSetHexQueryUTF8(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ">2;1T"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + ">0;3t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "2;6162" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "ab" {
		t.Fatalf("expected title ab, got %q", screen.Title)
	}
}

func TestEsctest2SMTitleSetUTF8QueryHex(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ">0;3T"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + ">2;1t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "2;ab" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "6162" {
		t.Fatalf("expected title 6162, got %q", screen.Title)
	}
}

func TestEsctest2SMTitleSetUTF8QueryUTF8(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ">0;1T"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + ">2;3t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "2;ab" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "ab" {
		t.Fatalf("expected title ab, got %q", screen.Title)
	}
}

func TestEsctest2SMTitleSetHexQueryHex(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + ">2;3T"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + ">0;1t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlOSC + "2;6162" + ControlST); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "6162" {
		t.Fatalf("expected title 6162, got %q", screen.Title)
	}
}
