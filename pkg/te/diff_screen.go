package te

// DiffScreen wraps Screen and tracks dirty line updates.
type DiffScreen struct {
	*Screen
}

// NewDiffScreen creates a DiffScreen with the given dimensions.
func NewDiffScreen(cols, lines int) *DiffScreen {
	return &DiffScreen{Screen: NewScreen(cols, lines)}
}
