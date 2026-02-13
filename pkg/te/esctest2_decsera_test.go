package te

import (
	"strings"
	"testing"
)

func esctestDecseraCharacters(alwaysBlank bool, data []string, point esctestPoint, count int) string {
	if alwaysBlank {
		return strings.Repeat(esctestBlank(), count)
	}
	var sb strings.Builder
	for i := 0; i < count; i++ {
		p := esctestPoint{X: point.X + i, Y: point.Y}
		if p.Y > len(data) {
			sb.WriteString(esctestBlank())
			continue
		}
		line := data[p.Y-1]
		if p.X > len(line) {
			sb.WriteString(esctestBlank())
			continue
		}
		if point.Y%2 == 1 {
			sb.WriteByte(line[p.X-1])
		} else {
			sb.WriteString(esctestBlank())
		}
	}
	return sb.String()
}

func esctestDecseraPrepare(t *testing.T, stream *Stream, data []string) {
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	protect := 1
	for _, line := range data {
		esctestDECSCA(t, stream, protect)
		esctestWrite(t, stream, line+ControlCR+ControlLF)
		protect = 1 - protect
	}
	esctestDECSCA(t, stream, 0)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_basic
func TestEsctestDecseraTestDecseraBasic(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	characters := func(point esctestPoint, count int) string {
		return esctestDecseraCharacters(false, fixture.data(), point, count)
	}
	esctestFillRectangleBasic(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill, characters)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_invalidRectDoesNothing
func TestEsctestDecseraTestDecseraInvalidRectDoesNothing(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleInvalidRectDoesNothing(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_defaultArgs
func TestEsctestDecseraTestDecseraDefaultArgs(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	characters := func(point esctestPoint, count int) string {
		return esctestDecseraCharacters(true, fixture.data(), point, count)
	}
	esctestFillRectangleDefaultArgs(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill, characters)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_respectsOriginMode
func TestEsctestDecseraTestDecseraRespectsOriginMode(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	characters := func(point esctestPoint, count int) string {
		return esctestDecseraCharacters(false, fixture.data(), point, count)
	}
	esctestFillRectangleRespectsOriginMode(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill, characters)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_overlyLargeSourceClippedToScreenSize
func TestEsctestDecseraTestDecseraOverlyLargeSourceClippedToScreenSize(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	characters := func(point esctestPoint, count int) string {
		return esctestDecseraCharacters(false, fixture.data(), point, count)
	}
	esctestFillRectangleOverlyLargeSourceClippedToScreenSize(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill, characters)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_cursorDoesNotMove
func TestEsctestDecseraTestDecseraCursorDoesNotMove(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleCursorDoesNotMove(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_ignoresMargins
func TestEsctestDecseraTestDecseraIgnoresMargins(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECSERAParams(t, fixture.stream, top, left, bottom, right)
	}
	characters := func(point esctestPoint, count int) string {
		return esctestDecseraCharacters(false, fixture.data(), point, count)
	}
	esctestFillRectangleIgnoresMargins(t, fixture, func() { esctestDecseraPrepare(t, fixture.stream, fixture.data()) }, fill, characters)
}

// From esctest2/esctest/tests/decsera.py::test_DECSERA_doesNotRespectISOProtect
func TestEsctestDecseraTestDecseraDoesNotRespectISOProtect(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"W")
	esctestDECSERAParams(t, stream, esctestIntPtr(1), esctestIntPtr(1), esctestIntPtr(1), esctestIntPtr(2))
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{strings.Repeat(esctestBlank(), 2)})
}
