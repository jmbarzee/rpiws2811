package rpiws2811

// **** <clk.h> ****

const (
	CM_CLK_CTL_PASSWD = (0x5a << 24)
	// CM_CLK_CTL_MASH
	CM_CLK_CTL_FLIP        = (1 << 8)
	CM_CLK_CTL_BUSY        = (1 << 7)
	CM_CLK_CTL_KILL        = (1 << 5)
	CM_CLK_CTL_ENAB        = (1 << 4)
	CM_CLK_CTL_SRC_GND     = (0 << 0)
	CM_CLK_CTL_SRC_OSC     = (1 << 0)
	CM_CLK_CTL_SRC_TSTDBG0 = (2 << 0)
	CM_CLK_CTL_SRC_TSTDBG1 = (3 << 0)
	CM_CLK_CTL_SRC_PLLA    = (4 << 0)
	CM_CLK_CTL_SRC_PLLC    = (5 << 0)
	CM_CLK_CTL_SRC_PLLD    = (6 << 0)
	CM_CLK_CTL_SRC_HDMIAUX = (7 << 0)
)

func CM_CLK_CTL_MASH(val uint32) uint32 {
	return uint32((val & 0x3) << 9)
}

const (
	CM_CLK_DIV_PASSWD = (0x5a << 24)
	// CM_CLK_DIV_DIVI
	// CM_CLK_DIV_DIVF
)

func CM_CLK_DIV_DIVI(val uint32) uint32 {
	return (val & 0xfff) << 12
}
func CM_CLK_DIV_DIVF(val uint32) uint32 {
	return ((val & 0xfff) << 0)
}

type cm_clk_t struct {
	ctl uint32
	div uint32
} // TODO @jmbarzee __attribute__((packed, aligned(4)))

const (
	// PWM and PCM clock offsets from https://www.scribd.com/doc/127599939/BCM2835-Audio-clocks
	CM_PCM_OFFSET = (0x00101098)
	CM_PWM_OFFSET = (0x001010a0)
)

// **** <clk.h> ****
