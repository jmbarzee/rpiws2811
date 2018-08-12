package rpiws2811

import (
	"fmt"
	"os"
	"time"
	"unsafe"
)

// **** <ws2811.h> ****

const (
	WS2811_TARGET_FREQ uint32 = 800000 // Can go as low as 400000

	// 4 color R, G, B and W ordering
	SK6812_STRIP_RGBW  LEDType = 0x18100800
	SK6812_STRIP_RBGW  LEDType = 0x18100008
	SK6812_STRIP_GRBW  LEDType = 0x18081000
	SK6812_STRIP_GBRW  LEDType = 0x18080010
	SK6812_STRIP_BRGW  LEDType = 0x18001008
	SK6812_STRIP_BGRW  LEDType = 0x18000810
	SK6812_SHIFT_WMASK LEDType = 0xf0000000

	// 3 color R, G and B ordering
	WS2811_STRIP_RGB LEDType = 0x00100800
	WS2811_STRIP_RBG LEDType = 0x00100008
	WS2811_STRIP_GRB LEDType = 0x00081000
	WS2811_STRIP_GBR LEDType = 0x00080010
	WS2811_STRIP_BRG LEDType = 0x00001008
	WS2811_STRIP_BGR LEDType = 0x00000810

	// predefined fixed LED types
	WS2812_STRIP  LEDType = WS2811_STRIP_GRB
	SK6812_STRIP  LEDType = WS2811_STRIP_GRB
	SK6812W_STRIP LEDType = SK6812_STRIP_GRBW
)

type (
	LEDType      uint32
	ws2811_led_t uint32

	ws2811_channel_t struct {
		gpionum    int           //< GPIO Pin with PWM alternate function, 0 if unused
		invert     bool          //< Invert output signal
		count      int           //< Number of LEDs, 0 if channel is unused
		strip_type LEDType       //< Strip color layout -- one of WS2811_STRIP_xxx constants
		leds       *ws2811_led_t //< LED buffers, allocated by driver based on count
		brightness byte          //< Brightness value between 0 and 255
		wshift     byte          //< White shift value
		rshift     byte          //< Red shift value
		gshift     byte          //< Green shift value
		bshift     byte          //< Blue shift value
		gamma      *byte         //< Gamma correction table
	}

	ws2811_t struct {
		render_wait_time uint64         //< time in µs before the next render can run
		device           *ws2811_device //< Private data for driver use
		rpi_hw           *rpi_hw_t      //< RPI Hardware Information
		freq             uint32         //< Required output frequency
		dmanum           int            //< DMA number _not_ already in use
		channel          [RPI_PWM_CHANNELS]ws2811_channel_t
	}

	ws2811_return_t int
)

const (
	WS2811_SUCCESS                = 0
	WS2811_ERROR_GENERIC          = -1
	WS2811_ERROR_OUT_OF_MEMORY    = -2
	WS2811_ERROR_HW_NOT_SUPPORTED = -3
	WS2811_ERROR_MEM_LOCK         = -4
	WS2811_ERROR_MMAP             = -5
	WS2811_ERROR_MAP_REGISTERS    = -6
	WS2811_ERROR_GPIO_INIT        = -7
	WS2811_ERROR_PWM_SETUP        = -8
	WS2811_ERROR_MAILBOX_DEVICE   = -9
	WS2811_ERROR_DMA              = -10
	WS2811_ERROR_ILLEGAL_GPIO     = -11
	WS2811_ERROR_PCM_SETUP        = -12
	WS2811_ERROR_SPI_SETUP        = -13
	WS2811_ERROR_SPI_TRANSFER     = -14
	WS2811_RETURN_STATE_COUNT     = -15 // I don't believe this is used anywhere...
)

func getWS2811ReturnMessage(returnCode ws2811_return_t) string {
	switch returnCode {
	case WS2811_SUCCESS:
		return "Success"
	case WS2811_ERROR_GENERIC:
		return "Generic failure"
	case WS2811_ERROR_OUT_OF_MEMORY:
		return "Out of memory"
	case WS2811_ERROR_HW_NOT_SUPPORTED:
		return "Hardware revision is not supported"
	case WS2811_ERROR_MEM_LOCK:
		return "Memory lock failed"
	case WS2811_ERROR_MMAP:
		return "mmap() failed"
	case WS2811_ERROR_MAP_REGISTERS:
		return "Unable to map registers into userspace"
	case WS2811_ERROR_GPIO_INIT:
		return "Unable to initialize GPIO"
	case WS2811_ERROR_PWM_SETUP:
		return "Unable to initialize PWM"
	case WS2811_ERROR_MAILBOX_DEVICE:
		return "Failed to create mailbox device"
	case WS2811_ERROR_DMA:
		return "DMA error"
	case WS2811_ERROR_ILLEGAL_GPIO:
		return "Selected GPIO not possible"
	case WS2811_ERROR_PCM_SETUP:
		return "Unable to initialize PCM"
	case WS2811_ERROR_SPI_SETUP:
		return "Unable to initialize SPI"
	case WS2811_ERROR_SPI_TRANSFER:
		return "SPI transfer error"
	}
	return ""
}

// **** <\ws2811.h> ****

// **** <ws2811.c> ****

const (
	OSC_FREQ = 19200000 // crystal = frequency

	/* 4 colors (R, G, B + W), 8 bits per byte, 3 symbols per bit + 55uS low for reset signal */
	LED_COLOURS  = 4
	LED_RESET_uS = 55

	/* Minimum time to wait for reset to occur in microseconds. */
	LED_RESET_WAIT_TIME = 300

	// Symbol definitions
	SYMBOL_HIGH = 0x6 // 1 1 = 0
	SYMBOL_LOW  = 0x4 // 1 0 = 0

	// Symbol definitions for software inversion (PCM and SPI only)
	SYMBOL_HIGH_INV = 0x1 // 0 0 = 1
	SYMBOL_LOW_INV  = 0x3 // 0 1 = 1

	// Driver mode definitions
	NONE = 0
	PWM  = 1
	PCM  = 2
	SPI  = 3
)

func BUS_TO_PHYS(x uint32) uint32 {
	return ^(^x | 0xC0000000)
}

func LED_BIT_COUNT(leds int, freq uint32) int {
	first := leds * LED_COLOURS * 8 * 3
	second := (LED_RESET_uS * (freq * 3)) / 1000000
	return first + int(second)
}

// Pad out to the nearest uint32 + 32-bits for idle low/high times the number of channels
func PWM_BYTE_COUNT(leds int, freq uint32) uint32 {
	return uint32(((((LED_BIT_COUNT(leds, freq) >> 3) & ^0x7) + 4) + 4) * RPI_PWM_CHANNELS)
}
func PCM_BYTE_COUNT(leds int, freq uint32) uint32 {
	return uint32((((LED_BIT_COUNT(leds, freq) >> 3) & ^0x7) + 4) + 4)
}

// We use the mailbox interface to request memory from the VideoCore.
// This lets us request one physically contiguous chunk, find its
// physical address, and map it 'uncached' so that writes from this
// code are immediately visible to the DMA controller.  This struct
// holds data relevant to the mailbox interface.
type videocore_mbox_t struct {
	handle    os.File        /* From mbox_open() */
	mem_ref   uint32         /* From mem_alloc() */
	bus_addr  uintptr        /* From mem_lock() */
	size      uint32         /* Size of allocation */
	virt_addr unsafe.Pointer /* From mapmem() */
}

type ws2811_device struct {
	driver_mode int
	pxl_raw     *uint8 // TODO @jmbarzee volatile
	dma         *dma_t // TODO @jmbarzee volatile
	pwm         *pwm_t // TODO @jmbarzee volatile
	pcm         *pcm_t // TODO @jmbarzee volatile
	spi_fd      int
	dma_cb      *dma_cb_t // TODO @jmbarzee volatile
	dma_cb_addr uint32
	gpio        *gpio_t   // TODO @jmbarzee volatile
	cm_clk      *cm_clk_t // TODO @jmbarzee volatile
	mbox        videocore_mbox_t
	max_count   int
}

//============================ PLACE HOLDER ============================

// **** </ws2811.c> ****

func get_microsecond_timestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func max_channel_led_count(strand *ws2811_t) int {
	max := 0
	for _, channel := range strand.channel {
		if channel.count > max {
			max = channel.count
		}
	}
	return max
}

func map_registers(strand *ws2811_t) error {
	device := strand.device
	rpi_hw := strand.rpi_hw // const
	base := strand.rpi_hw.periph_base
	var dma_addr uint32
	offset := uint32(0)

	dma_addr = dmanum_to_offset(strand.dmanum)
	if dma_addr == 0 {
		return fmt.Errorf("invalid dma_addr %v", dma_addr)
	}
	dma_addr += rpi_hw.periph_base

	device.dma = (*dma_t)(mapmem(dma_addr, unsafe.Sizeof(dma_t{}), DEV_MEM))
	if device.dma == nil {
		return fmt.Errorf("invalid device.dma nil")
	}

	switch device.driver_mode {
	case PWM:
		device.pwm = (*pwm_t)(mapmem(PWM_OFFSET+base, unsafe.Sizeof(pwm_t{}), DEV_MEM))
		if device.pwm == nil {
			return fmt.Errorf("invalid device.pwm nil")
		}
		break

	case PCM:
		device.pcm = (*pcm_t)(mapmem(PCM_OFFSET+base, unsafe.Sizeof(pcm_t{}), DEV_MEM))
		if device.pcm == nil {
			return fmt.Errorf("invalid device.pcm nil")
		}
		break
	}

	/*
	 * The below call can potentially work with /dev/gpiomem instead.
	 * However, it used /dev/mem before, so I'm leaving it as such.
	 */

	device.gpio = (*gpio_t)(mapmem(GPIO_OFFSET+base, unsafe.Sizeof(gpio_t{}), DEV_MEM))
	if device.gpio == nil {
		return fmt.Errorf("invalid device.gpio nil")
	}

	switch device.driver_mode {
	case PWM:
		offset = CM_PWM_OFFSET
		break
	case PCM:
		offset = CM_PCM_OFFSET
		break
	}
	device.cm_clk = (*cm_clk_t)(mapmem(offset+base, unsafe.Sizeof(cm_clk_t{}), DEV_MEM))
	if device.cm_clk == nil {
		return fmt.Errorf("invalid device.cm_clk nil")
	}

	return nil
}

func unmap_registers(strand *ws2811_t) {
	device := strand.device

	if device.dma != nil {
		unmapmem(unsafe.Pointer(device.dma), unsafe.Sizeof(dma_t{}))
	}

	if device.pwm != nil {
		unmapmem(unsafe.Pointer(device.pwm), unsafe.Sizeof(pwm_t{}))
	}

	if device.pcm != nil {
		unmapmem(unsafe.Pointer(device.pcm), unsafe.Sizeof(pcm_t{}))
	}

	if device.cm_clk != nil {
		unmapmem(unsafe.Pointer(device.cm_clk), unsafe.Sizeof(cm_clk_t{}))
	}

	if device.gpio != nil {
		unmapmem(unsafe.Pointer(device.gpio), unsafe.Sizeof(gpio_t{}))
	}

}

/**
 * Given a userspace address pointer, return the matching bus address used by DMA.
 *     Note: The bus address is not the same as the CPU physical address.
 *
 * @param    addr   Userspace virtual address pointer.
 *
 * @returns  Bus address for use by DMA.
 */
func addr_to_bus(device *ws2811_device, virt uintptr) uintptr {
	mbox := &device.mbox

	var sizeExample *uint8

	offset := uintptr(virt - uintptr(unsafe.Pointer(mbox.virt_addr))*unsafe.Sizeof(sizeExample))

	return mbox.bus_addr + offset
}

/**
 * Stop the PWM controller.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func stop_pwm(strand *ws2811_t) {
	pwm := strand.device.pwm
	cm_clk := strand.device.cm_clk

	// Turn off the PWM in case already running
	pwm.ctl = 0
	// TODO @jmbarzee discover what is going on here... waiting for writes to go through? 4 total
	time.Sleep(time.Microsecond * 10)

	// Kill the clock if it was already running
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_KILL
	time.Sleep(time.Microsecond * 10)

	// TODO @jmbarzee dear god, this may not work
	for (cm_clk.ctl & CM_CLK_CTL_BUSY) != 0 {
	}

}

/**
 * Stop the PCM controller.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func stop_pcm(strand *ws2811_t) {
	pcm := strand.device.pcm
	cm_clk := strand.device.cm_clk

	// Turn off the PCM in case already running
	pcm.cs = 0
	time.Sleep(time.Microsecond * 10)

	// Kill the clock if it was already running
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_KILL
	time.Sleep(time.Microsecond * 10)
	for (cm_clk.ctl & CM_CLK_CTL_BUSY) != 0 {
	}

}

/**
 * Setup the PWM controller in serial mode on both channels using DMA to feed the PWM FIFO.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func setup_pwm(strand *ws2811_t) {
	dma := strand.device.dma
	dma_cb := strand.device.dma_cb
	pwm := strand.device.pwm
	cm_clk := strand.device.cm_clk
	maxcount := strand.device.max_count
	freq := strand.freq
	var byte_count uint32

	stop_pwm(strand)

	// Setup the Clock - Use OSC @ 19.2Mhz w/ 3 clocks/tick
	cm_clk.div = CM_CLK_DIV_PASSWD | CM_CLK_DIV_DIVI(OSC_FREQ/(3*freq))
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_SRC_OSC
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_SRC_OSC | CM_CLK_CTL_ENAB
	time.Sleep(time.Microsecond * 10)
	for (cm_clk.ctl & CM_CLK_CTL_BUSY) == 0 {

	}

	// Setup the PWM, use delays as the block is rumored to lock up without them.  Make
	// sure to use a high enough priority to avoid any FIFO underruns, especially if
	// the CPU is busy doing lots of memory accesses, or another DMA controller is
	// busy.  The FIFO will clock out data at a much slower rate (2.6Mhz max), so
	// the odds of a DMA priority boost are extremely low.

	pwm.rng1 = 32 // 32-bits per word to serialize
	time.Sleep(time.Microsecond * 10)
	pwm.ctl = RPI_PWM_CTL_CLRF1
	time.Sleep(time.Microsecond * 10)
	pwm.dmac = RPI_PWM_DMAC_ENAB | RPI_PWM_DMAC_PANIC(7) | RPI_PWM_DMAC_DREQ(3)
	time.Sleep(time.Microsecond * 10)
	pwm.ctl = RPI_PWM_CTL_USEF1 | RPI_PWM_CTL_MODE1 |
		RPI_PWM_CTL_USEF2 | RPI_PWM_CTL_MODE2
	if strand.channel[0].invert {
		pwm.ctl |= RPI_PWM_CTL_POLA1
	}
	if strand.channel[1].invert {
		pwm.ctl |= RPI_PWM_CTL_POLA2
	}
	time.Sleep(time.Microsecond * 10)
	pwm.ctl |= RPI_PWM_CTL_PWEN1 | RPI_PWM_CTL_PWEN2

	// Initialize the DMA control block
	byte_count = PWM_BYTE_COUNT(maxcount, freq)
	dma_cb.ti = RPI_DMA_TI_NO_WIDE_BURSTS | // 32-bit transfers
		RPI_DMA_TI_WAIT_RESP | // wait for write complete
		RPI_DMA_TI_DEST_DREQ | // user peripheral flow control
		RPI_DMA_TI_PERMAP(5) | // PWM peripheral
		RPI_DMA_TI_SRC_INC // Increment src addr

	dma_cb.source_ad = addr_to_bus(strand.device, uintptr(unsafe.Pointer(strand.device.pxl_raw)))

	dma_cb.dest_ad = uintptr(unsafe.Pointer(&(*pwm_t)(unsafe.Pointer(uintptr(PWM_PERIPH_PHYS))).fif1))
	dma_cb.txfr_len = byte_count
	dma_cb.stride = 0
	dma_cb.nextconbk = 0

	dma.cs = 0
	dma.txfr_len = 0
}

/**
 * Setup the PCM controller with one 32-bit channel in a 32-bit frame using DMA to feed the PCM FIFO.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func setup_pcm(strand *ws2811_t) {
	dma := strand.device.dma
	dma_cb := strand.device.dma_cb
	pcm := strand.device.pcm
	cm_clk := strand.device.cm_clk
	//int maxcount := max_channel_led_count(ws2811)
	maxcount := strand.device.max_count
	freq := strand.freq
	var byte_count uint32

	stop_pcm(strand)

	// Setup the PCM Clock - Use OSC @ 19.2Mhz w/ 3 clocks/tick
	cm_clk.div = CM_CLK_DIV_PASSWD | CM_CLK_DIV_DIVI(OSC_FREQ/(3*freq))
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_SRC_OSC
	cm_clk.ctl = CM_CLK_CTL_PASSWD | CM_CLK_CTL_SRC_OSC | CM_CLK_CTL_ENAB
	time.Sleep(time.Microsecond * 10)
	for (cm_clk.ctl & CM_CLK_CTL_BUSY) == 0 {

	}

	// Setup the PCM, use delays as the block is rumored to lock up without them.  Make
	// sure to use a high enough priority to avoid any FIFO underruns, especially if
	// the CPU is busy doing lots of memory accesses, or another DMA controller is
	// busy.  The FIFO will clock out data at a much slower rate (2.6Mhz max), so
	// the odds of a DMA priority boost are extremely low.

	pcm.cs = RPI_PCM_CS_EN // Enable PCM hardware
	pcm.mode = (RPI_PCM_MODE_FLEN(31) | RPI_PCM_MODE_FSLEN(1))
	// Framelength 32, clock enabled, frame sync pulse
	pcm.txc = RPI_PCM_TXC_CH1WEX | RPI_PCM_TXC_CH1EN | RPI_PCM_TXC_CH1POS(0) | RPI_PCM_TXC_CH1WID(8)
	// Single 32-bit channel
	pcm.cs |= RPI_PCM_CS_TXCLR // Reset transmit fifo
	time.Sleep(time.Microsecond * 10)
	pcm.cs |= RPI_PCM_CS_DMAEN                                       // Enable DMA DREQ
	pcm.dreq = (RPI_PCM_DREQ_TX(0x3F) | RPI_PCM_DREQ_TX_PANIC(0x10)) // Set FIFO tresholds

	// Initialize the DMA control block
	byte_count = PCM_BYTE_COUNT(maxcount, freq)
	dma_cb.ti = RPI_DMA_TI_NO_WIDE_BURSTS | // 32-bit transfers
		RPI_DMA_TI_WAIT_RESP | // wait for write complete
		RPI_DMA_TI_DEST_DREQ | // user peripheral flow control
		RPI_DMA_TI_PERMAP(2) | // PCM TX peripheral
		RPI_DMA_TI_SRC_INC // Increment src addr

	dma_cb.source_ad = addr_to_bus(strand.device, uintptr(unsafe.Pointer(strand.device.pxl_raw)))
	dma_cb.dest_ad = uintptr(unsafe.Pointer(&(*pcm_t)(unsafe.Pointer(uintptr(PCM_PERIPH_PHYS))).fifo))
	dma_cb.txfr_len = byte_count
	dma_cb.stride = 0
	dma_cb.nextconbk = 0

	dma.cs = 0
	dma.txfr_len = 0
}

/**
 * Start the DMA feeding the PWM FIFO.  This will stream the entire DMA buffer out of both
 * PWM channels.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func dma_start(strand *ws2811_t) {
	dma := strand.device.dma
	pcm := strand.device.pcm
	dma_cb_addr := strand.device.dma_cb_addr

	dma.cs = RPI_DMA_CS_RESET
	time.Sleep(time.Microsecond * 10)

	dma.cs = RPI_DMA_CS_INT | RPI_DMA_CS_END
	time.Sleep(time.Microsecond * 10)

	dma.conblk_ad = dma_cb_addr
	dma.debug = 7 // clear debug error flags
	dma.cs = RPI_DMA_CS_WAIT_OUTSTANDING_WRITES |
		RPI_DMA_CS_PANIC_PRIORITY(15) |
		RPI_DMA_CS_PRIORITY(15) |
		RPI_DMA_CS_ACTIVE

	if strand.device.driver_mode == PCM {
		pcm.cs |= RPI_PCM_CS_TXON // Start transmission
	}
}

/**
 * Initialize the application selected GPIO pins for PWM/PCM operation.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  0 on success, -1 on unsupported pin
 */
func gpio_init(strand ws2811_t) error {
	gpio := strand.device.gpio

	for i, channel := range strand.channel {
		pinnum := channel.gpionum

		if pinnum != 0 {
			var altnum int
			var err error
			switch strand.device.driver_mode {
			case PWM:
				altnum, err = pwm_pin_alt(i, pinnum)
				if err != nil {
					return err
				}
				break
			case PCM:
				altnum, err = pcm_pin_alt(PCMFUN_DOUT, pinnum)
				if err != nil {
					return err
				}
				break
			default:
				return fmt.Errorf("Unrecognized driver_mode %v", strand.device.driver_mode)
			}

			gpio_function_set(gpio, pinnum, altnum)
		}
	}
	return nil
}

/**
 * Initialize the PWM DMA buffer with all zeros, inverted operation will be
 * handled by hardware.  The DMA buffer length is assumed to be a word
 * multiple.
 *
 * @param    ws2811  ws2811 instance pointer.
 *
 * @returns  None
 */
func pwm_raw_init(strand *ws2811_t) {
	pxl_raw := strand.device.pxl_raw
	maxcount := strand.device.max_count
	wordcount := (PWM_BYTE_COUNT(maxcount, strand.freq) / uint32(unsafe.Sizeof(uint32(0)))) / RPI_PWM_CHANNELS

	for i, channel := range strand.channel {
		wordpos := i

		for i := uint32(0); i < wordcount; i++ {
			pxl_raw[wordpos] = 0x0
			wordpos += 2
		}
	}
}
