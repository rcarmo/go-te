package te

type DiffScreen struct {
	*Screen
}

func NewDiffScreen(cols, lines int) *DiffScreen {
	return &DiffScreen{Screen: NewScreen(cols, lines)}
}
