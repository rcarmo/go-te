package te

const (
	ModeLNM = 20
	ModeIRM = 4

	ModeDECTCEM         = 25 << 5
	ModeDECSCNM         = 5 << 5
	ModeDECOM           = 6 << 5
	ModeDECAWM          = 7 << 5
	ModeDECCOLM         = 3 << 5
	ModeDECLRMM           = 69 << 5
	ModeReverseWrapInline = 45 << 5
	ModeReverseWrapExtend = 1045 << 5
	ModeDECSaveCursor     = 1048 << 5
	ModeAltBuf            = 47 << 5
	ModeAltBufOpt         = 1047 << 5
	ModeAltBufCursor      = 1049 << 5
	ModeAllow80To132      = 40 << 5
	ModeMoreFix           = 41 << 5
)
