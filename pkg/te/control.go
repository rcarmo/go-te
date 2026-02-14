package te

const (
	// ControlSP is the space character.
	ControlSP = " "
	// ControlNUL is the NUL control (0x00).
	ControlNUL = "\x00"
	// ControlBEL rings the terminal bell (BEL).
	ControlBEL = "\x07"
	// ControlBS is backspace (BS).
	ControlBS = "\x08"
	// ControlHT is horizontal tab (HT).
	ControlHT = "\t"
	// ControlLF is line feed (LF).
	ControlLF = "\n"
	// ControlVT is vertical tab (VT).
	ControlVT = "\x0b"
	// ControlFF is form feed (FF).
	ControlFF = "\x0c"
	// ControlCR is carriage return (CR).
	ControlCR = "\r"
	// ControlSO is shift out (SO) to select G1.
	ControlSO = "\x0e"
	// ControlSI is shift in (SI) to select G0.
	ControlSI = "\x0f"
	// ControlCAN is cancel (CAN).
	ControlCAN = "\x18"
	// ControlSUB is substitute (SUB).
	ControlSUB = "\x1a"
	// ControlESC is escape (ESC).
	ControlESC = "\x1b"
	// ControlDEL is delete (DEL).
	ControlDEL = "\x7f"

	// ControlCSIC0 is the 7-bit CSI introducer (ESC [).
	ControlCSIC0 = "\x1b["
	// ControlCSIC1 is the 8-bit CSI introducer (0x9b).
	ControlCSIC1 = "\x9b"
	// ControlCSI is the canonical CSI introducer (7-bit form).
	ControlCSI = ControlCSIC0

	// ControlSTC0 is the 7-bit string terminator (ESC \).
	ControlSTC0 = "\x1b\\"
	// ControlSTC1 is the 8-bit string terminator (0x9c).
	ControlSTC1 = "\x9c"
	// ControlST is the canonical string terminator (7-bit form).
	ControlST = ControlSTC0

	// ControlOSCC0 is the 7-bit OSC introducer (ESC ]).
	ControlOSCC0 = "\x1b]"
	// ControlOSCC1 is the 8-bit OSC introducer (0x9d).
	ControlOSCC1 = "\x9d"
	// ControlOSC is the canonical OSC introducer (7-bit form).
	ControlOSC = ControlOSCC0
	// ControlDCS is the device control string introducer (ESC P).
	ControlDCS = ControlESC + "P"
)
