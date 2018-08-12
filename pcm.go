package rpiws2811

import "fmt"

// **** <pcm.h> ****

/*
 *
 * Pin mapping of alternate pin configuration for PCM
 *
 * GPIO   ALT PCM_CLK   ALT PCM-FS   ALT PCM_DIN   ALT PCM_DOUT
 *
 *  18        0
 *  19                      0
 *  20                                   0
 *  21                                                 0
 *  28        2
 *  29                      2
 *  30                                   2
 *  31                                                 2
 *
 */
const (
	RPI_PCM_CS_STBY   = 1 << 25
	RPI_PCM_CS_SYNC   = 1 << 24
	RPI_PCM_CS_RXSEX  = 1 << 23
	RPI_PCM_CS_RXF    = 1 << 22
	RPI_PCM_CS_TXE    = 1 << 21
	RPI_PCM_CS_RXD    = 1 << 20
	RPI_PCM_CS_TXD    = 1 << 19
	RPI_PCM_CS_RXR    = 1 << 18
	RPI_PCM_CS_TXW    = 1 << 17
	RPI_PCM_CS_RXERR  = 1 << 16
	RPI_PCM_CS_TXERR  = 1 << 15
	RPI_PCM_CS_RXSYNC = 1 << 14
	RPI_PCM_CS_TXSYNC = 1 << 13
	RPI_PCM_CS_DMAEN  = 1 << 9
	// RPI_PCM_CS_RXTHR
	// RPI_PCM_CS_TXTHR
	RPI_PCM_CS_RXCLR = 1 << 4
	RPI_PCM_CS_TXCLR = 1 << 3
	RPI_PCM_CS_TXON  = 1 << 2
	RPI_PCM_CS_RXON  = 1 << 1
	RPI_PCM_CS_EN    = 1 << 0
)

func RPI_PCM_CS_RXTHR(val int) int { return ((val & 0x03) << 7) }
func RPI_PCM_CS_TXTHR(val int) int { return ((val & 0x03) << 5) }

const (
	RPI_PCM_MODE_CLK_DIS = 1 << 28
	RPI_PCM_MODE_PDMN    = 1 << 27
	RPI_PCM_MODE_PDME    = 1 << 26
	RPI_PCM_MODE_FRXP    = 1 << 25
	RPI_PCM_MODE_FTXP    = 1 << 24
	RPI_PCM_MODE_CLKM    = 1 << 23
	RPI_PCM_MODE_CLKI    = 1 << 22
	RPI_PCM_MODE_FSM     = 1 << 21
	RPI_PCM_MODE_FSI     = 1 << 20
	// RPI_PCM_MODE_FLEN
	// RPI_PCM_MODE_FSLEN
)

func RPI_PCM_MODE_FLEN(val int) uint32  { return uint32((val & 0x3ff) << 10) }
func RPI_PCM_MODE_FSLEN(val int) uint32 { return uint32((val & 0x3ff) << 0) }

const (
	RPI_PCM_RXC_CH1WEX = 1 << 31
	RPI_PCM_RXC_CH1EN  = 1 << 30
	// RPI_PCM_RXC_CH1POS
	// RPI_PCM_RXC_CH1WID
	RPI_PCM_RXC_CH2WEX = 1 << 15
	RPI_PCM_RXC_CH2EN  = 1 << 14
	// RPI_PCM_RXC_CH2POS
	// RPI_PCM_RXC_CH2WID
)

func RPI_PCM_RXC_CH1POS(val int) int { return ((val & 0x3ff) << 20) }
func RPI_PCM_RXC_CH1WID(val int) int { return ((val & 0x0f) << 16) }
func RPI_PCM_RXC_CH2POS(val int) int { return ((val & 0x3ff) << 4) }
func RPI_PCM_RXC_CH2WID(val int) int { return ((val & 0x0f) << 0) }

const (
	RPI_PCM_TXC_CH1WEX = uint32(1 << 31)
	RPI_PCM_TXC_CH1EN  = uint32(1 << 30)
	// RPI_PCM_TXC_CH1POS
	// RPI_PCM_TXC_CH1WID
	RPI_PCM_TXC_CH2WEX = uint32(1 << 15)
	RPI_PCM_TXC_CH2EN  = uint32(1 << 14)
	// RPI_PCM_TXC_CH2POS
	// RPI_PCM_TXC_CH2WID
)

func RPI_PCM_TXC_CH1POS(val uint32) uint32 { return (val & 0x3ff) << 20 }
func RPI_PCM_TXC_CH1WID(val uint32) uint32 { return (val & 0x0f) << 16 }
func RPI_PCM_TXC_CH2POS(val uint32) uint32 { return (val & 0x3ff) << 4 }
func RPI_PCM_TXC_CH2WID(val uint32) uint32 { return (val & 0x0f) << 0 }

func RPI_PCM_DREQ_TX_PANIC(val uint32) uint32 { return (val & 0x7f) << 24 }
func RPI_PCM_DREQ_RX_PANIC(val uint32) uint32 { return (val & 0x7f) << 16 }
func RPI_PCM_DREQ_TX(val uint32) uint32       { return (val & 0x7f) << 8 }
func RPI_PCM_DREQ_RX(val uint32) uint32       { return (val & 0x7f) << 0 }

const (
	RPI_PCM_INTEN_RXERR = uint32(1 << 3)
	RPI_PCM_INTEN_TXERR = uint32(1 << 2)
	RPI_PCM_INTEN_RXR   = uint32(1 << 1)
	RPI_PCM_INTEN_TXW   = uint32(1 << 0)

	RPI_PCM_INTSTC_RXERR = uint32(1 << 3)
	RPI_PCM_INTSTC_TXERR = uint32(1 << 2)
	RPI_PCM_INTSTC_RXR   = uint32(1 << 1)
	RPI_PCM_INTSTC_TXW   = uint32(1 << 0)

	// RPI_PCM_GRAY_RXFIFOLEVEL
	// RPI_PCM_GRAY_FLUSHED
	// RPI_PCM_GRAY_RXLEVEL
	RPI_PCM_GRAY_FLUSH = uint32(1 << 2)
	RPI_PCM_GRAY_CLR   = uint32(1 << 1)
	RPI_PCM_GRAY_EN    = uint32(1 << 0)
)

func RPI_PCM_GRAY_RXFIFOLEVEL(val int) int { return (val & 0x3f) << 16 }
func RPI_PCM_GRAY_FLUSHED(val int) int     { return (val & 0x3f) << 10 }
func RPI_PCM_GRAY_RXLEVEL(val int) int     { return (val & 0x3f) << 4 }

type pcm_t struct {
	cs     uint32
	fifo   uint32
	mode   uint32
	rxc    uint32
	txc    uint32
	dreq   uint32
	inten  uint32
	intstc uint32
	gray   uint32
} // TODD @jmbarzee __attribute__((packed, aligned(4)))

const (
	PCM_OFFSET      = 0x00203000
	PCM_PERIPH_PHYS = 0x7e203000

	NUM_PCMFUNS = 4
	PCMFUN_CLK  = 0
	PCMFUN_FS   = 1
	PCMFUN_DIN  = 2
	PCMFUN_DOUT = 3
)

type pcm_pin_table_t struct {
	pinnum int
	altnum int
}

type pcm_pin_tables_t struct {
	count int
	pins  []pcm_pin_table_t
}

// **** </pcm.h> ****

// **** <pcm.c> ****

// Mapping of Pin to alternate function for PCM_CLK
var pcm_pin_clk = []pcm_pin_table_t{
	pcm_pin_table_t{
		pinnum: 18,
		altnum: 0},
	pcm_pin_table_t{
		pinnum: 28,
		altnum: 2},
}

// Mapping of Pin to alternate function for PCM_FS
var pcm_pin_fs = []pcm_pin_table_t{
	pcm_pin_table_t{
		pinnum: 19,
		altnum: 0},
	pcm_pin_table_t{
		pinnum: 29,
		altnum: 2},
}

// Mapping of Pin to alternate function for PCM_DIN
var pcm_pin_din = []pcm_pin_table_t{
	pcm_pin_table_t{
		pinnum: 20,
		altnum: 0},
	pcm_pin_table_t{
		pinnum: 30,
		altnum: 2},
}

// Mapping of Pin to alternate function for PCM_DOUT
var pcm_pin_dout = []pcm_pin_table_t{
	pcm_pin_table_t{
		pinnum: 21,
		altnum: 0},
	pcm_pin_table_t{
		pinnum: 31,
		altnum: 2},
}

var pcm_pin_tables = [NUM_PCMFUNS]pcm_pin_tables_t{
	pcm_pin_tables_t{
		pins:  pcm_pin_clk,
		count: len(pcm_pin_clk)},
	pcm_pin_tables_t{
		pins:  pcm_pin_fs,
		count: len(pcm_pin_fs)},
	pcm_pin_tables_t{
		pins:  pcm_pin_din,
		count: len(pcm_pin_din)},
	pcm_pin_tables_t{
		pins:  pcm_pin_dout,
		count: len(pcm_pin_dout)},
}

func pcm_pin_alt(pcmfun int, pinnum int) (int, error) {
	if pcmfun < 0 || pcmfun > 3 {
		return 0, fmt.Errorf("pcmfun out of acceptable range: %v", pcmfun)
	}
	pins := pcm_pin_tables[pcmfun].pins

	for _, pin := range pins {
		if pin.pinnum == pinnum {
			return pin.altnum, nil
		}
	}

	return 0, fmt.Errorf("no alternate pin found for pin: %v", pcmfun)
}

// **** </pcm.c> ****
