package te

import "testing"

func TestEsctest2APCBasic(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "_xyz" + ControlST + "A"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}

func TestEsctest2APC8Bit(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("\x1b G\u009fxyz\u009cA\x1b F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}

func TestEsctest2PMBasic(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "^xyz" + ControlST + "A"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}

func TestEsctest2PM8Bit(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("\x1b G\u009exyz\u009cA\x1b F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}

func TestEsctest2SOSBasic(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "Xxyz" + ControlST + "A"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}

func TestEsctest2SOS8Bit(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("\x1b G\u0098xyz\u009cA\x1b F"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	assertCell(t, screen, 1, 1, "A")
	assertCell(t, screen, 2, 1, " ")
}
