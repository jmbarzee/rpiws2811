package rpiws2811

// **** <gpio.h> ****

type gpio_t struct {
	fsel       [6]uint32 // GPIO Function Select
	resvd_0x18 uint32
	set        [2]uint32 // GPIO Pin Output Set
	resvd_0x24 uint32
	clr        [2]uint32 // GPIO Pin Output Clear
	resvd_0x30 uint32
	lev        [2]uint32 // GPIO Pin Level
	resvd_0x3c uint32
	eds        [2]uint32 // GPIO Pin Event Detect Status
	resvd_0x48 uint32
	ren        [2]uint32 // GPIO Pin Rising Edge Detect Enable
	resvd_0x54 uint32
	fen        [2]uint32 // GPIO Pin Falling Edge Detect Enable
	resvd_0x60 uint32
	hen        [2]uint32 // GPIO Pin High Detect Enable
	resvd_0x6c uint32
	len        [2]uint32 // GPIO Pin Low Detect Enable
	resvd_0x78 uint32
	aren       [2]uint32 // GPIO Pin Async Rising Edge Detect
	resvd_0x84 uint32
	afen       [2]uint32 // GPIO Pin Async Falling Edge Detect
	resvd_0x90 uint32
	pud        uint32    // GPIO Pin Pull up/down Enable
	pudclk     [2]uint32 // GPIO Pin Pull up/down Enable Clock
	resvd_0xa0 [4]uint32
	test       uint32
} // TODO @jmbarzee __attribute__((packed, aligned(4))) gpio_t;

const (
	GPIO_OFFSET = 0x00200000
)

func gpio_function_set(gpio *gpio_t, pin uint8, function uint8) {
	regnum := int(pin / 10)
	offset := int((pin % 10) * 3)
	funcmap = []uint8{4, 5, 6, 7, 3, 2} // See datasheet for mapping

	if function > 5 {
		return
	}

	gpio.fsel[regnum] &= ^(0x7 << offset)
	gpio.fsel[regnum] |= ((funcmap[function]) << offset)
}

func gpio_level_set(gpio *gpio_t, pin uint8, level uint8) {
	regnum = int(pin >> 5)
	offset = int(pin & 0x1f)

	if level {
		gpio.set[regnum] = (1 << offset)
	} else {
		gpio.clr[regnum] = (1 << offset)
	}
}

func gpio_output_set(gpio *gpio_t, pin uint8, output uint8) {
	regnum = int(pin / 10)
	offset = int((pin % 10) * 3)
	function := uint8(0)
	if output {
		function = 1 // See datasheet for mapping
	}

	gpio.fsel[regnum] &= ^(0x7 << offset)
	gpio.fsel[regnum] |= ((function & 0x7) << offset)
}
