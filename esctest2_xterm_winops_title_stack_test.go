package te

import "testing"

func TestEsctest2XtermWinopsPushPopTitle(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	screen.Title = "one"
	screen.IconName = "icon1"
	if err := stream.Feed(ControlCSI + "22;0t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.Title = "two"
	screen.IconName = "icon2"
	if err := stream.Feed(ControlCSI + "23;0t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "one" || screen.IconName != "icon1" {
		t.Fatalf("unexpected titles: %q %q", screen.Title, screen.IconName)
	}
}

func TestEsctest2XtermWinopsPushPopWindowTitle(t *testing.T) {
	screen := NewScreen(10, 2)
	stream := NewStream(screen, false)
	screen.Title = "one"
	if err := stream.Feed(ControlCSI + "22;2t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	screen.Title = "two"
	if err := stream.Feed(ControlCSI + "23;2t"); err != nil {
		t.Fatalf("feed: %v", err)
	}
	if screen.Title != "one" {
		t.Fatalf("unexpected title: %q", screen.Title)
	}
}
