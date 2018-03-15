package rpiws2811

import "fmt"

// **** <pwm.h> ****

const (
	RPI_PWM_CHANNELS = 2

	RPI_PWM_CTL_MSEN2 = 1 << 15
	RPI_PWM_CTL_USEF2 = 1 << 13
	RPI_PWM_CTL_POLA2 = 1 << 12
	RPI_PWM_CTL_SBIT2 = 1 << 11
	RPI_PWM_CTL_RPTL2 = 1 << 10
	RPI_PWM_CTL_MODE2 = 1 << 9
	RPI_PWM_CTL_PWEN2 = 1 << 8
	RPI_PWM_CTL_MSEN1 = 1 << 7
	RPI_PWM_CTL_CLRF1 = 1 << 6
	RPI_PWM_CTL_USEF1 = 1 << 5
	RPI_PWM_CTL_POLA1 = 1 << 4
	RPI_PWM_CTL_SBIT1 = 1 << 3
	RPI_PWM_CTL_RPTL1 = 1 << 2
	RPI_PWM_CTL_MODE1 = 1 << 1
	RPI_PWM_CTL_PWEN1 = 1 << 0

	RPI_PWM_STA_STA4  = 1 << 12
	RPI_PWM_STA_STA3  = 1 << 11
	RPI_PWM_STA_STA2  = 1 << 10
	RPI_PWM_STA_STA1  = 1 << 9
	RPI_PWM_STA_BERR  = 1 << 8
	RPI_PWM_STA_GAP04 = 1 << 7
	RPI_PWM_STA_GAP03 = 1 << 6
	RPI_PWM_STA_GAP02 = 1 << 5
	RPI_PWM_STA_GAP01 = 1 << 4
	RPI_PWM_STA_RERR1 = 1 << 3
	RPI_PWM_STA_WERR1 = 1 << 2
	RPI_PWM_STA_EMPT1 = 1 << 1
	RPI_PWM_STA_FULL1 = 1 << 0

	RPI_PWM_DMAC_ENAB = 1 << 31
)

func RPI_PWM_DMAC_PANIC(val int) int { return ((val & 0xff) << 8) }
func RPI_PWM_DMAC_DREQ(val int) int  { return ((val & 0xff) << 0) }

type pwm_t struct {
	ctl        uint32
	sta        uint32
	dmac       uint32
	resvd_0x0c uint32
	rng1       uint32
	dat1       uint32
	fif1       uint32
	resvd_0x1c uint32
	rng2       uint32
	dat2       uint32
} // TODO @jmbarzee __attribute__((packed, aligned(4)))

const (
	PWM_OFFSET      = 0x0020c000
	PWM_PERIPH_PHYS = 0x7e20c000
)

type pwm_pin_table_t struct {
	pinnum int
	altnum int
}

type pwm_pin_tables_t struct {
	count int
	pins  *pwm_pin_table_t
}

// **** </pwm.h> ****

// **** <pwm.c> ****

// Mapping of Pin to alternate function for PWM channel 0
var pwm_pin_chan0 = []pwm_pin_table_t{
	pwm_pin_table_t{
		pinnum: 12,
		altnum: 0},
	pwm_pin_table_t{
		pinnum: 18,
		altnum: 5},
	pwm_pin_table_t{
		pinnum: 40,
		altnum: 0},
}

// Mapping of Pin to alternate function for PWM channel 1
var pwm_pin_chan1 = []pwm_pin_table_t{
	pwm_pin_table_t{
		pinnum: 13,
		altnum: 0},
	pwm_pin_table_t{
		pinnum: 19,
		altnum: 5},
	pwm_pin_table_t{
		pinnum: 41,
		altnum: 0},
	pwm_pin_table_t{
		pinnum: 45,
		altnum: 0},
}

var pwm_pin_tables = [RPI_PWM_CHANNELS]pwm_pin_tables_t{
	pwm_pin_tables_t{
		pins:  pwm_pin_chan0,
		count: len(pwm_pin_chan0)},
	pwm_pin_tables_t{
		pins:  pwm_pin_chan1,
		count: len(pwm_pin_chan1)},
}

func pwm_pin_alt(channel int, pinnum int) (int, error) {
	pintable = &pwm_pin_tables[channel]

	for _, pin := range pintable {
		if pin.pinnum == pinnum {
			return pin.altnum, nil
		}
	}

	return 0, fmt.Errorf("no alternate pin found for pin: %v", pcmfun)
}

// **** </pwm.c> ****
