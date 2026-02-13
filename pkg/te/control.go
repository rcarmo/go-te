package te

const (
	ControlSP  = " "
	ControlNUL = "\x00"
	ControlBEL = "\x07"
	ControlBS  = "\x08"
	ControlHT  = "\t"
	ControlLF  = "\n"
	ControlVT  = "\x0b"
	ControlFF  = "\x0c"
	ControlCR  = "\r"
	ControlSO  = "\x0e"
	ControlSI  = "\x0f"
	ControlCAN = "\x18"
	ControlSUB = "\x1a"
	ControlESC = "\x1b"
	ControlDEL = "\x7f"

	ControlCSIC0 = "\x1b["
	ControlCSIC1 = "\x9b"
	ControlCSI   = ControlCSIC0

	ControlSTC0 = "\x1b\\"
	ControlSTC1 = "\x9c"
	ControlST   = ControlSTC0

	ControlOSCC0 = "\x1b]"
	ControlOSCC1 = "\x9d"
	ControlOSC   = ControlOSCC0
	ControlDCS   = ControlESC + "P"
)
