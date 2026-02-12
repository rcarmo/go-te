package te

import "testing"

func TestEsctest2CUPDefaults(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "3;6" + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 2 || screen.Cursor.Col != 5 {
		t.Fatalf("expected cursor 3,6")
	}
	if err := stream.Feed(ControlCSI + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor 1,1")
	}
	if err := stream.Feed(ControlCSI + "2" + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected row 2 col 1")
	}
	if err := stream.Feed(ControlCSI + ";2" + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 1 {
		t.Fatalf("expected row 1 col 2")
	}
	if err := stream.Feed(ControlCSI + "0;0" + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 0 || screen.Cursor.Col != 0 {
		t.Fatalf("expected row 1 col 1")
	}
	if err := stream.Feed(ControlCSI + "999;999" + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 4 || screen.Cursor.Col != 9 {
		t.Fatalf("expected clamp to bounds")
	}
}

func TestEsctest2CursorMoves(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "2" + EscCUD); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 2 {
		t.Fatalf("expected row 2")
	}
	if err := stream.Feed(ControlCSI + "3" + EscCUF); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Col != 3 {
		t.Fatalf("expected col 3")
	}
	if err := stream.Feed(ControlCSI + "2" + EscCUB); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Col != 1 {
		t.Fatalf("expected col 1")
	}
	if err := stream.Feed(ControlCSI + "1" + EscCUU); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 1 {
		t.Fatalf("expected row 1")
	}
}

func TestEsctest2LineCursorModes(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "3" + EscCNL); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 3 || screen.Cursor.Col != 0 {
		t.Fatalf("expected CNL to move down and CR")
	}
	if err := stream.Feed(ControlCSI + "2" + EscCPL); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 0 {
		t.Fatalf("expected CPL to move up and CR")
	}
	if err := stream.Feed(ControlCSI + "5" + EscCHA); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Col != 4 {
		t.Fatalf("expected CHA column 5")
	}
}

func TestEsctest2VPAHVP(t *testing.T) {
	screen := NewScreen(10, 5)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlCSI + "3" + EscVPA); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 2 {
		t.Fatalf("expected row 3")
	}
	if err := stream.Feed(ControlCSI + "2" + EscVPR); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 4 {
		t.Fatalf("expected row 5")
	}
	if err := stream.Feed(ControlCSI + "2;3" + EscHVP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Row != 1 || screen.Cursor.Col != 2 {
		t.Fatalf("expected cursor 2,3")
	}
}

func TestEsctest2InsertDeleteErase(t *testing.T) {
	screen := NewScreen(5, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed("abcde"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "3" + EscCHA); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "1" + EscICH); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed("Z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "abZcd" {
		t.Fatalf("unexpected insert")
	}
	if err := stream.Feed(ControlCSI + "2" + EscDCH); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "abZ  " {
		t.Fatalf("unexpected delete")
	}
	if err := stream.Feed(ControlCSI + "2" + "b"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "abZZZ" {
		t.Fatalf("unexpected repeat")
	}
	if err := stream.Feed(ControlCSI + "1" + EscECH); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "abZZZ" {
		t.Fatalf("unexpected erase")
	}
}

func TestEsctest2EraseDisplayLine(t *testing.T) {
	screen := NewScreen(5, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("hello"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "1" + EscED); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "     " {
		t.Fatalf("expected cleared screen")
	}
	if err := stream.Feed("world"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "2" + EscEL); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "     " {
		t.Fatalf("expected cleared line")
	}
}

func TestEsctest2AlignmentDisplay(t *testing.T) {
	screen := NewScreen(4, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + "#" + EscDECALN); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "EEEE" || screen.Display()[1] != "EEEE" {
		t.Fatalf("expected alignment display")
	}
}

func TestEsctest2IndexReverseIndex(t *testing.T) {
	screen := NewScreen(3, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("ab\r\ncd"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlESC + EscIND); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "cd " {
		t.Fatalf("expected index scroll")
	}
	if err := stream.Feed(ControlCSI + EscCUP); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlESC + EscRI); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "   " {
		t.Fatalf("expected reverse index scroll")
	}
}

func TestEsctest2RIS(t *testing.T) {
	screen := NewScreen(3, 2)
	stream := NewStream(screen, false)
	if err := stream.Feed("abc"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlESC + EscRIS); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Display()[0] != "   " {
		t.Fatalf("expected reset")
	}
}

func TestEsctest2TabStops(t *testing.T) {
	screen := NewScreen(10, 1)
	stream := NewStream(screen, false)
	if err := stream.Feed(ControlESC + EscHTS); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if err := stream.Feed(ControlCSI + "3" + EscTBC); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if len(screen.TabStops) != 0 {
		t.Fatalf("expected tabstops cleared")
	}
	if err := stream.Feed(ControlCSI + "1" + "Z"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Cursor.Col != 0 {
		t.Fatalf("expected cursor back tab")
	}
}
