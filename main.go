package rpiws2811

import (
	"fmt"
	"os"
	"os/signal"
)

const (
// STRIP_TYPE            WS2811_STRIP_RGB		// WS2812/SK6812RGB integrated chip+leds
// STRIP_TYPE            WS2811_STRIP_GBR       // WS2812/SK6812RGB integrated chip+leds
// STRIP_TYPE            SK6812_STRIP_RGBW		// SK6812RGBW (NOT SK6812RGB)
)

var (

	// TODO @jmbarzee global variables to remove
	running       = true
	clear_on_exit = false
)

func init_handlers() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		running = false // TODO @jbarzee
		// change this to follow normal go conventions for passing messaging
	}()
}

func main() {
	init_handlers()

	// Create Channels
	c1, err := NewLEDStrandChannel(12, 30, 255, false, SK6812_STRIP_RGBW)
	if err != nil {
		panic(err)
	}
	c2, err := NewLEDStrandChannel(0, 0, 0, false, SK6812_STRIP_RGBW)
	if err != nil {
		panic(err)
	}

	// Create Strand
	_, err = NewLEDStrand(
		WS2811_TARGET_FREQ,
		10, // DMA
		true,
		c1,
		c2,
	)
	if err != nil {
		panic(err)
	}

	//TODO @jmbarzee use strand
}

func NewLEDStrand(freq uint32, dma int, clearOnExit bool, c1, c2 ws2811_channel_t) (ws2811_t, error) {
	strand := ws2811_t{}

	clear_on_exit = clearOnExit

	if dma < 14 {
		return strand, fmt.Errorf("invalid dma %v\n", dma)
	}
	strand.dmanum = dma

	if freq < 400000 || freq > 800000 {
		return strand, fmt.Errorf("invalid freq %v\n", freq)
	}
	strand.freq = freq

	strand.channel = [RPI_PWM_CHANNELS]ws2811_channel_t{
		c1,
		c2,
	}
	return strand, nil
}

func NewLEDStrandChannel(gpio, length, brightness int, invert bool, LEDType LEDType) (ws2811_channel_t, error) {

	/*              ====== GPIO ======
	PWM0, which can be set to use GPIOs 12, 18, 40, and 52.
	Only 12 (pin 32) and 18 (pin 12) are available on the B+/2B/3B
	PWM1 which can be set to use GPIOs 13, 19, 41, 45 and 53.
	Only 13 is available on the B+/2B/PiZero/3B, on pin 33
	PCM_DOUT, which can be set to use GPIOs 21 and 31.
	Only 21 is available on the B+/2B/PiZero/3B, on pin 40.
	SPI0-MOSI is available on GPIOs 10 and 38.
	Only GPIO 10 is available on all models.

	The library checks if the specified gpio is available
	on the specific model (from model B rev 1 till 3B)
	*/
	channel := ws2811_channel_t{}

	if length < 0 {
		return channel, fmt.Errorf("invalid length %v\n", length)
	}
	channel.count = length

	channel.gpionum = gpio
	channel.invert = invert
	channel.strip_type = LEDType

	return channel, nil
}
