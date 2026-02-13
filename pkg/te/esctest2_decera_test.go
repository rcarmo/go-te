package te

import (
	"strings"
	"testing"
)

func esctestDeceraCharacters(point esctestPoint, count int) string {
	return strings.Repeat(esctestBlank(), count)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_basic
func TestEsctestDeceraTestDeceraBasic(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleBasic(t, fixture, nil, fill, esctestDeceraCharacters)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_invalidRectDoesNothing
func TestEsctestDeceraTestDeceraInvalidRectDoesNothing(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleInvalidRectDoesNothing(t, fixture, nil, fill)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_defaultArgs
func TestEsctestDeceraTestDeceraDefaultArgs(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleDefaultArgs(t, fixture, nil, fill, esctestDeceraCharacters)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_respectsOriginMode
func TestEsctestDeceraTestDeceraRespectsOriginMode(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleRespectsOriginMode(t, fixture, nil, fill, esctestDeceraCharacters)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_overlyLargeSourceClippedToScreenSize
func TestEsctestDeceraTestDeceraOverlyLargeSourceClippedToScreenSize(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleOverlyLargeSourceClippedToScreenSize(t, fixture, nil, fill, esctestDeceraCharacters)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_cursorDoesNotMove
func TestEsctestDeceraTestDeceraCursorDoesNotMove(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleCursorDoesNotMove(t, fixture, nil, fill)
}

// From esctest2/esctest/tests/decera.py::test_DECERA_ignoresMargins
func TestEsctestDeceraTestDeceraIgnoresMargins(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECERA(t, fixture.stream, top, left, bottom, right)
	}
	esctestFillRectangleIgnoresMargins(t, fixture, nil, fill, esctestDeceraCharacters)
}
