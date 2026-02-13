package te

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type esctestPoint struct {
	X int
	Y int
}

type esctestSize struct {
	Width  int
	Height int
}

type esctestRect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

const (
	esctestModeDECOM   = 6
	esctestModeDECAWM  = 7
	esctestModeDECLRMM = 69
)

func esctestCUP(t *testing.T, stream *Stream, point esctestPoint) {
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d%s", ControlCSI, point.Y, point.X, EscCUP))
}

func esctestDECSET(t *testing.T, stream *Stream, modes ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s?%sh", ControlCSI, esctestJoinParams(modes...)))
}

func esctestDECRESET(t *testing.T, stream *Stream, modes ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s?%sl", ControlCSI, esctestJoinParams(modes...)))
}

func esctestDECSLRM(t *testing.T, stream *Stream, left, right int) {
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d%s", ControlCSI, left, right, "s"))
}

func esctestDECSTBM(t *testing.T, stream *Stream, top, bottom int) {
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d%s", ControlCSI, top, bottom, EscDECSTBM))
}

func esctestIND(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscIND)
}

func esctestNEL(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscNEL)
}

func esctestWrite(t *testing.T, stream *Stream, data string) {
	if err := stream.Feed(data); err != nil {
		t.Fatalf("feed: %v", err)
	}
}

func esctestJoinParams(params ...int) string {
	parts := make([]string, len(params))
	for i, param := range params {
		parts[i] = fmt.Sprintf("%d", param)
	}
	return strings.Join(parts, ";")
}

func esctestGetCursorPosition(screen *Screen) esctestPoint {
	return esctestPoint{X: screen.Cursor.Col + 1, Y: screen.Cursor.Row + 1}
}

func esctestGetScreenSize(screen *Screen) esctestSize {
	return esctestSize{Width: screen.Columns, Height: screen.Lines}
}

func esctestAssertEQ(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func esctestEmpty() string {
	return " "
}

func esctestAssertScreenCharsInRectEqual(t *testing.T, screen *Screen, rect esctestRect, expected []string) {
	rows := rect.Bottom - rect.Top + 1
	if rows != len(expected) {
		t.Fatalf("expected %d rows, got %d", rows, len(expected))
	}
	for row := rect.Top; row <= rect.Bottom; row++ {
		line := screen.Buffer[row-1]
		var sb strings.Builder
		for col := rect.Left; col <= rect.Right; col++ {
			idx := col - 1
			if idx >= 0 && idx < len(line) {
				sb.WriteString(line[idx].Data)
			} else {
				sb.WriteString(" ")
			}
		}
		got := sb.String()
		want := expected[row-rect.Top]
		if got != want {
			t.Fatalf("row %d: expected %q, got %q", row, want, got)
		}
	}
}
