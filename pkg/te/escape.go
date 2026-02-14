package te

const (
	// EscRIS is the RIS (reset to initial state) escape sequence (ESC c).
	EscRIS = "c"
	// EscIND is the IND (index) escape sequence (ESC D).
	EscIND = "D"
	// EscNEL is the NEL (next line) escape sequence (ESC E).
	EscNEL = "E"
	// EscHTS is the HTS (horizontal tab set) escape sequence (ESC H).
	EscHTS = "H"
	// EscRI is the RI (reverse index) escape sequence (ESC M).
	EscRI = "M"
	// EscDECSC is the DECSC (save cursor) escape sequence (ESC 7).
	EscDECSC = "7"
	// EscDECRC is the DECRC (restore cursor) escape sequence (ESC 8).
	EscDECRC = "8"

	// EscDECALN is the DECALN (screen alignment test) escape sequence (ESC # 8).
	EscDECALN = "8"

	// EscICH is the ICH (insert character) CSI final.
	EscICH = "@"
	// EscCUU is the CUU (cursor up) CSI final.
	EscCUU = "A"
	// EscCUD is the CUD (cursor down) CSI final.
	EscCUD = "B"
	// EscCUF is the CUF (cursor forward) CSI final.
	EscCUF = "C"
	// EscCUB is the CUB (cursor backward) CSI final.
	EscCUB = "D"
	// EscCNL is the CNL (cursor next line) CSI final.
	EscCNL = "E"
	// EscCPL is the CPL (cursor previous line) CSI final.
	EscCPL = "F"
	// EscCHA is the CHA (cursor horizontal absolute) CSI final.
	EscCHA = "G"
	// EscCUP is the CUP (cursor position) CSI final.
	EscCUP = "H"
	// EscED is the ED (erase in display) CSI final.
	EscED = "J"
	// EscEL is the EL (erase in line) CSI final.
	EscEL = "K"
	// EscIL is the IL (insert line) CSI final.
	EscIL = "L"
	// EscDL is the DL (delete line) CSI final.
	EscDL = "M"
	// EscDCH is the DCH (delete character) CSI final.
	EscDCH = "P"
	// EscECH is the ECH (erase character) CSI final.
	EscECH = "X"
	// EscHPR is the HPR (horizontal position relative) CSI final.
	EscHPR = "a"
	// EscDA is the DA (device attributes) CSI final.
	EscDA = "c"
	// EscVPA is the VPA (vertical position absolute) CSI final.
	EscVPA = "d"
	// EscVPR is the VPR (vertical position relative) CSI final.
	EscVPR = "e"
	// EscHVP is the HVP (horizontal and vertical position) CSI final.
	EscHVP = "f"
	// EscTBC is the TBC (tab clear) CSI final.
	EscTBC = "g"
	// EscSM is the SM (set mode) CSI final.
	EscSM = "h"
	// EscRM is the RM (reset mode) CSI final.
	EscRM = "l"
	// EscSGR is the SGR (select graphic rendition) CSI final.
	EscSGR = "m"
	// EscDSR is the DSR (device status report) CSI final.
	EscDSR = "n"
	// EscDECSTBM is the DECSTBM (set top/bottom margins) CSI final.
	EscDECSTBM = "r"
	// EscHPA is the HPA (horizontal position absolute) CSI final.
	EscHPA = "'"
)
