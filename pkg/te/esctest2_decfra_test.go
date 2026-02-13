package te

import (
	"strings"
	"testing"
)

const esctestDecfraCharacter = "%"

func esctestDecfraCharacters(point esctestPoint, count int) string {
	return strings.Repeat(esctestDecfraCharacter, count)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_basic
func TestEsctestDecfraTestDecfraBasic(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleBasic(t, fixture, nil, fill, esctestDecfraCharacters)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_invalidRectDoesNothing
func TestEsctestDecfraTestDecfraInvalidRectDoesNothing(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleInvalidRectDoesNothing(t, fixture, nil, fill)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_defaultArgs
func TestEsctestDecfraTestDecfraDefaultArgs(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleDefaultArgs(t, fixture, nil, fill, esctestDecfraCharacters)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_respectsOriginMode
func TestEsctestDecfraTestDecfraRespectsOriginMode(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleRespectsOriginMode(t, fixture, nil, fill, esctestDecfraCharacters)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_overlyLargeSourceClippedToScreenSize
func TestEsctestDecfraTestDecfraOverlyLargeSourceClippedToScreenSize(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleOverlyLargeSourceClippedToScreenSize(t, fixture, nil, fill, esctestDecfraCharacters)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_cursorDoesNotMove
func TestEsctestDecfraTestDecfraCursorDoesNotMove(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleCursorDoesNotMove(t, fixture, nil, fill)
}

// From esctest2/esctest/tests/decfra.py::test_DECFRA_ignoresMargins
func TestEsctestDecfraTestDecfraIgnoresMargins(t *testing.T) {
	fixture := newEsctestFillRectangleFixture()
	fill := func(top, left, bottom, right *int) {
		esctestDECFRA(t, fixture.stream, int(esctestDecfraCharacter[0]), top, left, bottom, right)
	}
	esctestFillRectangleIgnoresMargins(t, fixture, nil, fill, esctestDecfraCharacters)
}
