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
	esctestModeDECOM             = 6
	esctestModeDECAWM            = 7
	esctestModeIRM               = 4
	esctestModeDECLRMM           = 69
	esctestModeReverseWrapInline = 45
	esctestModeReverseWrapExtend = 1045
	esctestModeAllow80To132      = 40
	esctestModeAltBuf            = 47
	esctestModeSaveRestoreCursor = 1048
	esctestModeDECRLM            = 34
	esctestXtermReverseWrap      = 383

	esctestTitleSetHex    = 0
	esctestTitleQueryHex  = 1
	esctestTitleSetUTF8   = 2
	esctestTitleQueryUTF8 = 3
)

func esctestCUP(t *testing.T, stream *Stream, point esctestPoint) {
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d%s", ControlCSI, point.Y, point.X, EscCUP))
}

func esctestCUPParams(t *testing.T, stream *Stream, row, col *int) {
	if row == nil && col == nil {
		esctestWrite(t, stream, ControlCSI+EscCUP)
		return
	}
	params := []string{}
	if row == nil {
		params = append(params, "")
	} else {
		params = append(params, fmt.Sprintf("%d", *row))
	}
	if col != nil {
		params = append(params, fmt.Sprintf("%d", *col))
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, strings.Join(params, ";"), EscCUP))
}

func esctestHVPParams(t *testing.T, stream *Stream, row, col *int) {
	if row == nil && col == nil {
		esctestWrite(t, stream, ControlCSI+EscHVP)
		return
	}
	params := []string{}
	if row == nil {
		params = append(params, "")
	} else {
		params = append(params, fmt.Sprintf("%d", *row))
	}
	if col != nil {
		params = append(params, fmt.Sprintf("%d", *col))
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, strings.Join(params, ";"), EscHVP))
}

func esctestCHT(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"I")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sI", ControlCSI, esctestJoinParams(params...)))
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

func esctestDECSTBM(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscDECSTBM)
		return
	}
	if len(params) == 1 {
		esctestWrite(t, stream, fmt.Sprintf("%s%d%s", ControlCSI, params[0], EscDECSTBM))
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d%s", ControlCSI, params[0], params[1], EscDECSTBM))
}

func esctestIND(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscIND)
}

func esctestNEL(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscNEL)
}

func esctestRI(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscRI)
}

func esctestRIS(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscRIS)
}

func esctestDECSC(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscDECSC)
}

func esctestDECRC(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscDECRC)
}

func esctestSCOSC(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlCSI+"s")
}

func esctestSCORC(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlCSI+"u")
}

func esctestReverseWraparoundMode() int {
	if esctestXtermReverseWrap >= 383 {
		return esctestModeReverseWrapExtend
	}
	return esctestModeReverseWrapInline
}

func esctestED(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscED)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscED))
}

func esctestEL(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscEL)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscEL))
}

func esctestECH(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscECH)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscECH))
}

func esctestREP(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"b")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sb", ControlCSI, esctestJoinParams(params...)))
}

func esctestRMTitle(t *testing.T, stream *Stream, params ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s>%sT", ControlCSI, esctestJoinParams(params...)))
}

func esctestSMTitle(t *testing.T, stream *Stream, params ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s>%st", ControlCSI, esctestJoinParams(params...)))
}

func esctestSM(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"h")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sh", ControlCSI, esctestJoinParams(params...)))
}

func esctestRM(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"l")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sl", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSTR(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlCSI+"!p")
}

func esctestDECSASD(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"$}")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s$}", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSCA(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"\"q")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s\"q", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSED(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"?J")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s?%sJ", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSERA(t *testing.T, stream *Stream, top, left, bottom, right int) {
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d;%d;%d${", ControlCSI, top, left, bottom, right))
}

func esctestChangeWindowTitle(t *testing.T, stream *Stream, title string) {
	esctestWrite(t, stream, fmt.Sprintf("%s2;%s%s", ControlOSC, title, ControlST))
}

func esctestChangeIconTitle(t *testing.T, stream *Stream, title string) {
	esctestWrite(t, stream, fmt.Sprintf("%s1;%s%s", ControlOSC, title, ControlST))
}

func esctestGetWindowTitle(screen *Screen) string {
	return screen.Title
}

func esctestGetIconTitle(screen *Screen) string {
	return screen.IconName
}

func esctestSU(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"S")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sS", ControlCSI, esctestJoinParams(params...)))
}

func esctestSD(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"T")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sT", ControlCSI, esctestJoinParams(params...)))
}

func esctestDCH(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscDCH)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscDCH))
}

func esctestICH(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscICH)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscICH))
}

func esctestHTS(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+EscHTS)
}

func esctestTBC(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscTBC)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscTBC))
}

func esctestCBT(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"Z")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%sZ", ControlCSI, esctestJoinParams(params...)))
}

func esctestCHA(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCHA)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCHA))
}

func esctestVPA(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscVPA)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscVPA))
}

func esctestVPR(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscVPR)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscVPR))
}

func esctestHPA(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscHPA)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscHPA))
}

func esctestHPR(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscHPR)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscHPR))
}

func esctestCUU(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCUU)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCUU))
}

func esctestCUD(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCUD)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCUD))
}

func esctestCUF(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCUF)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCUF))
}

func esctestCUB(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCUB)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCUB))
}

func esctestCNL(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCNL)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCNL))
}

func esctestCPL(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscCPL)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscCPL))
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
	x := screen.Cursor.Col + 1
	y := screen.Cursor.Row + 1
	if screen.isModeSet(ModeDECOM) && screen.Margins != nil {
		y -= screen.Margins.Top
	}
	if screen.isModeSet(ModeDECOM) && screen.isModeSet(ModeDECLRMM) {
		x -= screen.leftMargin
	}
	return esctestPoint{X: x, Y: y}
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

func esctestBlank() string {
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
