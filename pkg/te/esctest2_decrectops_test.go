package te

import (
	"strings"
	"testing"
)

type esctestFillRectangleFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestFillRectangleFixture() esctestFillRectangleFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestFillRectangleFixture{screen: screen, stream: stream}
}

func (f esctestFillRectangleFixture) data() []string {
	return []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	}
}

func (f esctestFillRectangleFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	for _, line := range f.data() {
		esctestWrite(t, f.stream, line+ControlCR+ControlLF)
	}
}

func (f esctestFillRectangleFixture) fillCharacters(point esctestPoint, count int) string {
	return strings.Repeat("!", count)
}

func esctestFillRectangleBasic(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int), characters func(point esctestPoint, count int) string) {
	if prepare != nil {
		prepare()
	} else {
		fixture.prepare(t)
	}
	top, left, bottom, right := 5, 5, 7, 7
	fill(&top, &left, &bottom, &right)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCD" + characters(esctestPoint{X: 5, Y: 5}, 3) + "H",
		"IJKL" + characters(esctestPoint{X: 5, Y: 6}, 3) + "P",
		"QRST" + characters(esctestPoint{X: 5, Y: 7}, 3) + "X",
		"YZ6789!@",
	})
}

func esctestFillRectangleInvalidRectDoesNothing(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int)) {
	if prepare != nil {
		prepare()
	} else {
		fixture.prepare(t)
	}
	top, left, bottom, right := 5, 5, 4, 4
	fill(&top, &left, &bottom, &right)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

func esctestFillRectangleDefaultArgs(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int), characters func(point esctestPoint, count int) string) {
	if prepare != nil {
		prepare()
	}
	size := esctestGetScreenSize(fixture.screen)
	points := []esctestPoint{
		{X: 1, Y: 1},
		{X: size.Width, Y: 1},
		{X: size.Width, Y: size.Height},
		{X: 1, Y: size.Height},
	}
	n := 1
	for _, point := range points {
		esctestCUP(t, fixture.stream, point)
		esctestWrite(t, fixture.stream, string(rune('0'+n)))
		n++
	}
	fill(nil, nil, nil, nil)
	for _, point := range points {
		esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: point.X, Top: point.Y, Right: point.X, Bottom: point.Y}, []string{characters(point, 1)})
	}
}

func esctestFillRectangleRespectsOriginMode(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int), characters func(point esctestPoint, count int) string) {
	if prepare != nil {
		prepare()
	} else {
		fixture.prepare(t)
	}
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 2, 9)
	esctestDECSTBM(t, fixture.stream, 2, 9)
	esctestDECSET(t, fixture.stream, esctestModeDECOM)
	top, left, bottom, right := 1, 1, 3, 3
	fill(&top, &left, &bottom, &right)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestDECRESET(t, fixture.stream, esctestModeDECOM)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"i" + characters(esctestPoint{X: 2, Y: 2}, 3) + "mnop",
		"q" + characters(esctestPoint{X: 2, Y: 3}, 3) + "uvwx",
		"y" + characters(esctestPoint{X: 2, Y: 4}, 3) + "2345",
		"ABCDEFGH",
		"IJKLMNOP",
		"QRSTUVWX",
		"YZ6789!@",
	})
}

func esctestFillRectangleOverlyLargeSourceClippedToScreenSize(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int), characters func(point esctestPoint, count int) string) {
	if prepare != nil {
		prepare()
	}
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: size.Height - 1})
	esctestWrite(t, fixture.stream, "ab")
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: size.Height})
	esctestWrite(t, fixture.stream, "cd")
	top, left, bottom, right := size.Height, size.Width, size.Height+10, size.Width+10
	fill(&top, &left, &bottom, &right)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: size.Width - 1, Top: size.Height - 1, Right: size.Width, Bottom: size.Height}, []string{
		"ab",
		"c" + characters(esctestPoint{X: size.Width, Y: size.Height}, 1),
	})
}

func esctestFillRectangleCursorDoesNotMove(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int)) {
	if prepare != nil {
		prepare()
	} else {
		fixture.prepare(t)
	}
	position := esctestPoint{X: 3, Y: 4}
	esctestCUP(t, fixture.stream, position)
	top, left, bottom, right := 2, 2, 4, 4
	fill(&top, &left, &bottom, &right)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), position)
}

func esctestFillRectangleIgnoresMargins(t *testing.T, fixture esctestFillRectangleFixture, prepare func(), fill func(top, left, bottom, right *int), characters func(point esctestPoint, count int) string) {
	if prepare != nil {
		prepare()
	} else {
		fixture.prepare(t)
	}
	esctestDECSET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSLRM(t, fixture.stream, 3, 6)
	esctestDECSTBM(t, fixture.stream, 3, 6)
	top, left, bottom, right := 5, 5, 7, 7
	fill(&top, &left, &bottom, &right)
	esctestDECRESET(t, fixture.stream, esctestModeDECLRMM)
	esctestDECSTBM(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 8}, []string{
		"abcdefgh",
		"ijklmnop",
		"qrstuvwx",
		"yz012345",
		"ABCD" + characters(esctestPoint{X: 5, Y: 5}, 3) + "H",
		"IJKL" + characters(esctestPoint{X: 5, Y: 6}, 3) + "P",
		"QRST" + characters(esctestPoint{X: 5, Y: 7}, 3) + "X",
		"YZ6789!@",
	})
}
