package te

type DebugScreen struct {
	*Screen
	log []string
}

func NewDebugScreen(cols, lines int) *DebugScreen {
	return &DebugScreen{
		Screen: NewScreen(cols, lines),
		log:    []string{},
	}
}

func (s *DebugScreen) PutRune(ch rune) {
	s.log = append(s.log, string(ch))
	s.Screen.PutRune(ch)
}

func (s *DebugScreen) Log() []string {
	return append([]string(nil), s.log...)
}
