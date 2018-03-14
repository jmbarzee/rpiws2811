package rpiws2811

// **** <dma.h> ****

/*
 * DMA Control Block in Main Memory
 *
 * Note: Must start at a 256 byte aligned address.
 *       Use corresponding register field definitions.
 */
type dma_cb_t struct {
	ti         uint32
	source_ad  uint32
	dest_ad    uint32
	txfr_len   uint32
	stride     uint32
	nextconbk  uint32
	resvd_0x18 [2]uint32
} // TODO @jmbarzee __attribute__((packed, aligned(4)))

const (
	RPI_DMA_CS_RESET                   = (1 << 31)
	RPI_DMA_CS_ABORT                   = (1 << 30)
	RPI_DMA_CS_DISDEBUG                = (1 << 29)
	RPI_DMA_CS_WAIT_OUTSTANDING_WRITES = (1 << 28)
	// RPI_DMA_CS_PANIC_PRIORITY
	// RPI_DMA_CS_PRIORITY
	RPI_DMA_CS_ERROR                      = (1 << 8)
	RPI_DMA_CS_WAITING_OUTSTANDING_WRITES = (1 << 6)
	RPI_DMA_CS_DREQ_STOPS_DMA             = (1 << 5)
	RPI_DMA_CS_PAUSED                     = (1 << 4)
	RPI_DMA_CS_DREQ                       = (1 << 3)
	RPI_DMA_CS_INT                        = (1 << 2)
	RPI_DMA_CS_END                        = (1 << 1)
	RPI_DMA_CS_ACTIVE                     = (1 << 0)
)

func RPI_DMA_CS_PANIC_PRIORITY(val int) int {
	return (val & 0xf) << 20
}
func RPI_DMA_CS_PRIORITY(val int) int {
	return (val & 0xf) << 16
}

const (
	RPI_DMA_TI_NO_WIDE_BURSTS = (1 << 26)
	// RPI_DMA_TI_WAITS
	// RPI_DMA_TI_PERMAP
	// RPI_DMA_TI_BURST_LENGTH
	RPI_DMA_TI_SRC_IGNORE  = (1 << 11)
	RPI_DMA_TI_SRC_DREQ    = (1 << 10)
	RPI_DMA_TI_SRC_WIDTH   = (1 << 9)
	RPI_DMA_TI_SRC_INC     = (1 << 8)
	RPI_DMA_TI_DEST_IGNORE = (1 << 7)
	RPI_DMA_TI_DEST_DREQ   = (1 << 6)
	RPI_DMA_TI_DEST_WIDTH  = (1 << 5)
	RPI_DMA_TI_DEST_INC    = (1 << 4)
	RPI_DMA_TI_WAIT_RESP   = (1 << 3)
	RPI_DMA_TI_TDMODE      = (1 << 1)
	RPI_DMA_TI_INTEN       = (1 << 0)
)

func RPI_DMA_TI_WAITS(val int) int {
	return (val & 0x1f) << 21
}
func RPI_DMA_TI_PERMAP(val int) int {
	return (val & 0x1f) << 16
}
func RPI_DMA_TI_BURST_LENGTH(val int) int {
	return (val & 0xf) << 12
}

/*
 * DMA register set
 */
type dma_t struct {
	cs        uint32
	conblk_ad uint32
	ti        uint32
	source_ad uint32
	dest_ad   uint32
	txfr_len  uint32
	// RPI_DMA_TXFR_LEN_YLENGTH
	// RPI_DMA_TXFR_LEN_XLENGTH
	stride uint32
	// RPI_DMA_STRIDE_D_STRIDE
	// RPI_DMA_STRIDE_S_STRIDE
	nextconbk uint32
	debug     uint32
} // TODO @jmbarzee  __attribute__((packed, aligned(4)))

func RPI_DMA_TXFR_LEN_YLENGTH(val int) int {
	return (val & 0xffff) << 16
}
func RPI_DMA_TXFR_LEN_XLENGTH(val int) int {
	return (val & 0xffff) << 0
}

func RPI_DMA_STRIDE_D_STRIDE(val int) int {
	return (val & 0xffff) << 16
}
func RPI_DMA_STRIDE_S_STRIDE(val int) int {
	return ((val & 0xffff) << 0)
}

const (
	DMA0_OFFSET  = (0x00007000)
	DMA1_OFFSET  = (0x00007100)
	DMA2_OFFSET  = (0x00007200)
	DMA3_OFFSET  = (0x00007300)
	DMA4_OFFSET  = (0x00007400)
	DMA5_OFFSET  = (0x00007500)
	DMA6_OFFSET  = (0x00007600)
	DMA7_OFFSET  = (0x00007700)
	DMA8_OFFSET  = (0x00007800)
	DMA9_OFFSET  = (0x00007900)
	DMA10_OFFSET = (0x00007a00)
	DMA11_OFFSET = (0x00007b00)
	DMA12_OFFSET = (0x00007c00)
	DMA13_OFFSET = (0x00007d00)
	DMA14_OFFSET = (0x00007e00)
	DMA15_OFFSET = (0x00e05000)

	PAGE_SIZE = (1 << 12)
	PAGE_MASK = (^(PAGE_SIZE - 1))
)

func PAGE_OFFSET(page int) int {
	return (page & (PAGE_SIZE - 1))
}

// **** </dma.h> ****
// **** <dma.c> ****
var dma_offset = []uint32{
	DMA0_OFFSET,
	DMA1_OFFSET,
	DMA2_OFFSET,
	DMA3_OFFSET,
	DMA4_OFFSET,
	DMA5_OFFSET,
	DMA6_OFFSET,
	DMA7_OFFSET,
	DMA8_OFFSET,
	DMA9_OFFSET,
	DMA10_OFFSET,
	DMA11_OFFSET,
	DMA12_OFFSET,
	DMA13_OFFSET,
	DMA14_OFFSET,
	DMA15_OFFSET,
}

func dmanum_to_offset(dmanum int) uint32 {
	if dmanum >= len(dma_offset) {
		return 0
	}
	return dma_offset[dmanum]
}

// **** </dma.h> ****
