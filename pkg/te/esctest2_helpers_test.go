package te

import (
	"fmt"
	"reflect"
	"strconv"
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
	esctestModeLNM               = 20
	esctestModeDECLRMM           = 69
	esctestModeReverseWrapInline = 45
	esctestModeReverseWrapExtend = 1045
	esctestModeAllow80To132      = 40
	esctestModeAltBuf            = 47
	esctestModeOptAltBuf         = 1047
	esctestModeOptAltBufCursor   = 1049
	esctestModeSaveRestoreCursor = 1048
	esctestModeDECRLM            = 34
	esctestModeMoreFix           = 41
	esctestModeDECCKM            = 1
	esctestModeDECSCLM           = 4
	esctestModeDECARM            = 8
	esctestModeDECPFF            = 18
	esctestModeDECPEX            = 19
	esctestModeDECNRCM           = 42
	esctestModeDECNKM            = 66
	esctestModeDECBKM            = 67
	esctestModeDECKBUM           = 68
	esctestModeDECNCSM           = 95
	esctestModeDECOSCNM          = 106
	esctestModeDECHCCM           = 60
	esctestModeDECAAM            = 100
	esctestModeDECCANSM          = 101
	esctestModeDECNULM           = 102
	esctestModeDECHDPXM          = 103
	esctestModeDECESKM           = 104
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

func esctestDECFI(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+"9")
}

func esctestDECBI(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+"6")
}

func esctestDECALN(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlESC+"#"+EscDECALN)
}

func esctestWriteDCS(t *testing.T, stream *Stream, data string) {
	esctestWrite(t, stream, ControlESC+"P"+data+ControlST)
}

func esctestWriteAPC(t *testing.T, stream *Stream, data string) {
	esctestWrite(t, stream, ControlESC+"_"+data+ControlST)
}

func esctestWritePM(t *testing.T, stream *Stream, data string) {
	esctestWrite(t, stream, ControlESC+"^"+data+ControlST)
}

func esctestWriteSOS(t *testing.T, stream *Stream, data string) {
	esctestWrite(t, stream, ControlESC+"X"+data+ControlST)
}

func esctestChangeColor(t *testing.T, stream *Stream, params ...string) {
	esctestWrite(t, stream, ControlOSC+"4;"+strings.Join(params, ";")+ControlST)
}

func esctestXtermWinops(t *testing.T, stream *Stream, params ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s%st", ControlCSI, esctestJoinParams(params...)))
}

func esctestGetIndexedColors() int {
	return 16
}

func esctestChangeDynamicColor(t *testing.T, stream *Stream, params ...string) {
	esctestWrite(t, stream, ControlOSC+strings.Join(params, ";")+ControlST)
}

func esctestChangeSpecialColor(t *testing.T, stream *Stream, params ...string) {
	if len(params) > 0 {
		if index, err := strconv.Atoi(params[0]); err == nil && index >= 10 {
			esctestWrite(t, stream, ControlOSC+strings.Join(params, ";")+ControlST)
			return
		}
	}
	offset := esctestGetIndexedColors()
	parts := make([]string, len(params))
	copy(parts, params)
	for i := 0; i < len(parts); i += 2 {
		index, err := strconv.Atoi(parts[i])
		if err != nil {
			continue
		}
		parts[i] = fmt.Sprintf("%d", index+offset)
	}
	esctestWrite(t, stream, ControlOSC+"4;"+strings.Join(parts, ";")+ControlST)
}

func esctestChangeSpecialColor2(t *testing.T, stream *Stream, params ...string) {
	if len(params) > 0 {
		if index, err := strconv.Atoi(params[0]); err == nil && index >= 10 {
			esctestWrite(t, stream, ControlOSC+strings.Join(params, ";")+ControlST)
			return
		}
	}
	esctestWrite(t, stream, ControlOSC+"5;"+strings.Join(params, ";")+ControlST)
}

func esctestResetSpecialColor(t *testing.T, stream *Stream, params ...string) {
	esctestWrite(t, stream, ControlOSC+"105;"+strings.Join(params, ";")+ControlST)
}

func esctestResetDynamicColor(t *testing.T, stream *Stream, code string) {
	esctestWrite(t, stream, ControlOSC+code+ControlST)
}

func esctestResetColor(t *testing.T, stream *Stream, params ...string) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlOSC+"104"+ControlST)
		return
	}
	esctestWrite(t, stream, ControlOSC+"104;"+strings.Join(params, ";")+ControlST)
}

func esctestXtermSave(t *testing.T, stream *Stream, params ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s?%ss", ControlCSI, esctestJoinParams(params...)))
}

func esctestXtermRestore(t *testing.T, stream *Stream, params ...int) {
	esctestWrite(t, stream, fmt.Sprintf("%s?%sr", ControlCSI, esctestJoinParams(params...)))
}

func esctestCaptureResponse(screen *Screen, fn func()) string {
	var response string
	prev := screen.WriteProcessInput
	screen.WriteProcessInput = func(data string) { response = data }
	fn()
	screen.WriteProcessInput = prev
	return response
}

func esctestParseCSI(t *testing.T, response string, prefix rune) []int {
	return esctestReadCSI(t, response, 'c', prefix)
}

func esctestReadOSC(t *testing.T, response string, prefix string) string {
	if !strings.HasPrefix(response, ControlOSC+prefix) {
		t.Fatalf("expected OSC %s response, got %q", prefix, response)
	}
	payload := strings.TrimPrefix(response, ControlOSC+prefix)
	payload = strings.TrimSuffix(payload, ControlST)
	return payload
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

func esctestSGR(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+EscSGR)
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%s", ControlCSI, esctestJoinParams(params...), EscSGR))
}

func esctestReadCSI(t *testing.T, response string, expectedFinal rune, expectedPrefix rune) []int {
	if !strings.HasPrefix(response, ControlCSI) {
		t.Fatalf("expected CSI response, got %q", response)
	}
	payload := strings.TrimPrefix(response, ControlCSI)
	if expectedPrefix != 0 {
		if len(payload) == 0 || rune(payload[0]) != expectedPrefix {
			t.Fatalf("expected CSI prefix %q, got %q", string(expectedPrefix), response)
		}
		payload = payload[1:]
	}
	if !strings.HasSuffix(payload, string(expectedFinal)) {
		t.Fatalf("expected CSI final %q, got %q", string(expectedFinal), response)
	}
	payload = strings.TrimSuffix(payload, string(expectedFinal))
	if payload == "" {
		return nil
	}
	parts := strings.Split(payload, ";")
	params := make([]int, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			t.Fatalf("invalid CSI param %q", part)
		}
		params = append(params, value)
	}
	return params
}

func esctestDECIC(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"'}")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s'}", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECDC(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"'~")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s'~", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECCRA(t *testing.T, stream *Stream, srcTop, srcLeft, srcBottom, srcRight, srcPage, dstTop, dstLeft, dstPage *int) {
	params := esctestJoinOptionalParams(srcTop, srcLeft, srcBottom, srcRight, srcPage, dstTop, dstLeft, dstPage)
	esctestWrite(t, stream, ControlCSI+params+"$v")
}

func esctestDECSTR(t *testing.T, stream *Stream) {
	esctestWrite(t, stream, ControlCSI+"!p")
}

func esctestDECRQM(t *testing.T, stream *Stream, mode int, dec bool) {
	prefix := ""
	if dec {
		prefix = "?"
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s%d$p", ControlCSI, prefix, mode))
}

func esctestDECRQSS(t *testing.T, stream *Stream, query string) {
	esctestWriteDCS(t, stream, "$q"+query)
}

func esctestDECELF(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"+q")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s+q", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECLFKC(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"*}")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s*}", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSACE(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"*x")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s*x", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSCL(t *testing.T, stream *Stream, level int, sevenBit *int) {
	if sevenBit == nil {
		esctestWrite(t, stream, fmt.Sprintf("%s%d\"p", ControlCSI, level))
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%d;%d\"p", ControlCSI, level, *sevenBit))
}

func esctestDECSMKR(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"+r")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s+r", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSNLS(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"*|")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s*|", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSSDT(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"$~")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s$~", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECSCUSR(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+" q")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s q", ControlCSI, esctestJoinParams(params...)))
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

func esctestDECSEL(t *testing.T, stream *Stream, params ...int) {
	if len(params) == 0 {
		esctestWrite(t, stream, ControlCSI+"?K")
		return
	}
	esctestWrite(t, stream, fmt.Sprintf("%s?%sK", ControlCSI, esctestJoinParams(params...)))
}

func esctestDECERA(t *testing.T, stream *Stream, top, left, bottom, right *int) {
	params := esctestJoinOptionalParams(top, left, bottom, right)
	esctestWrite(t, stream, fmt.Sprintf("%s%s$z", ControlCSI, params))
}

func esctestDECFRA(t *testing.T, stream *Stream, ch int, top, left, bottom, right *int) {
	params := []string{fmt.Sprintf("%d", ch)}
	optional := []*int{top, left, bottom, right}
	for _, value := range optional {
		if value == nil {
			params = append(params, "")
			continue
		}
		params = append(params, fmt.Sprintf("%d", *value))
	}
	for len(params) > 0 && params[len(params)-1] == "" {
		params = params[:len(params)-1]
	}
	esctestWrite(t, stream, fmt.Sprintf("%s%s$x", ControlCSI, strings.Join(params, ";")))
}

func esctestDECSERAParams(t *testing.T, stream *Stream, top, left, bottom, right *int) {
	params := esctestJoinOptionalParams(top, left, bottom, right)
	esctestWrite(t, stream, fmt.Sprintf("%s%s${", ControlCSI, params))
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

func esctestJoinOptionalParams(values ...*int) string {
	parts := make([]string, len(values))
	for i, value := range values {
		if value == nil {
			parts[i] = ""
			continue
		}
		parts[i] = fmt.Sprintf("%d", *value)
	}
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, ";")
}

func esctestIntPtr(value int) *int {
	return &value
}

func esctestGetCursorPosition(screen *Screen) esctestPoint {
	var response string
	prev := screen.WriteProcessInput
	screen.WriteProcessInput = func(data string) { response = data }
	screen.ReportDeviceStatus(6, false, 0)
	screen.WriteProcessInput = prev
	if response == "" {
		return esctestPoint{X: screen.Cursor.Col + 1, Y: screen.Cursor.Row + 1}
	}
	response = strings.TrimPrefix(response, ControlCSI)
	response = strings.TrimSuffix(response, "R")
	parts := strings.Split(response, ";")
	if len(parts) < 2 {
		return esctestPoint{X: screen.Cursor.Col + 1, Y: screen.Cursor.Row + 1}
	}
	y, errY := strconv.Atoi(parts[0])
	x, errX := strconv.Atoi(parts[1])
	if errY != nil || errX != nil {
		return esctestPoint{X: screen.Cursor.Col + 1, Y: screen.Cursor.Row + 1}
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

func esctestAtoi(value string) int {
	out := 0
	for _, r := range value {
		if r < '0' || r > '9' {
			return out
		}
		out = out*10 + int(r-'0')
	}
	return out
}

func esctestItoa(value int) string {
	if value == 0 {
		return "0"
	}
	buf := make([]byte, 0, 6)
	for value > 0 {
		d := value % 10
		buf = append([]byte{byte('0' + d)}, buf...)
		value /= 10
	}
	return string(buf)
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

func esctestAssertCharHasSGR(t *testing.T, screen *Screen, point esctestPoint, missing, present []int) {
	cell := screen.Buffer[point.Y-1][point.X-1]
	for _, attr := range missing {
		if esctestCharHasSGR(cell, attr) {
			t.Fatalf("expected attr %d to be unset", attr)
		}
	}
	for _, attr := range present {
		if !esctestCharHasSGR(cell, attr) {
			t.Fatalf("expected attr %d to be set", attr)
		}
	}
}

func esctestCharHasSGR(cell Cell, attr int) bool {
	switch attr {
	case 1:
		return cell.Attr.Bold
	case 22:
		return !cell.Attr.Bold
	case 30:
		return cell.Attr.Fg.Mode == ColorANSI16 && cell.Attr.Fg.Index == 0
	case 39:
		return cell.Attr.Fg.Mode == ColorDefault
	case 40:
		return cell.Attr.Bg.Mode == ColorANSI16 && cell.Attr.Bg.Index == 0
	case 49:
		return cell.Attr.Bg.Mode == ColorDefault
	default:
		return false
	}
}
