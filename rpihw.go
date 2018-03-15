package rpiws2811

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// **** <rpihw.h> ****

const (
	RPI_HWVER_TYPE_UNKNOWN = 0
	RPI_HWVER_TYPE_PI1     = 1
	RPI_HWVER_TYPE_PI2     = 2
)

type rpi_hw_t struct {
	typeNum        uint32
	hwver          uint32
	periph_base    uint32
	videocore_base uint32
	desc           string
}

// **** </rpihw.h> ****
// **** <rpihw.c> ****

const (
	LINE_WIDTH_MAX = 80
	HW_VER_STRING  = "Revision"

	PERIPH_BASE_RPI  = 0x20000000
	PERIPH_BASE_RPI2 = 0x3f000000

	VIDEOCORE_BASE_RPI  = 0x40000000
	VIDEOCORE_BASE_RPI2 = 0xc0000000

	RPI_MANUFACTURER_MASK = (0xf << 16)
	RPI_WARRANTY_MASK     = (0x3 << 24)
)

var rpi_hw_info = []rpi_hw_t{
	//
	// Model B Rev 1.0
	//
	rpi_hw_t{
		hwver:          0x02,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},
	rpi_hw_t{
		hwver:          0x03,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},

	//
	// Model B Rev 2.0
	//
	rpi_hw_t{
		hwver:          0x04,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},
	rpi_hw_t{
		hwver:          0x05,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},
	rpi_hw_t{
		hwver:          0x06,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},

	//
	// Model A
	//
	rpi_hw_t{
		hwver:          0x07,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A"},
	rpi_hw_t{
		hwver:          0x08,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A"},
	rpi_hw_t{
		hwver:          0x09,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A"},

	//
	// Model B
	//
	rpi_hw_t{
		hwver:          0x0d,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},
	rpi_hw_t{
		hwver:          0x0e,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},
	rpi_hw_t{
		hwver:          0x0f,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B"},

	//
	// Model B+
	//
	rpi_hw_t{
		hwver:          0x10,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B+"},
	rpi_hw_t{
		hwver:          0x13,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B+"},
	rpi_hw_t{
		hwver:          0x900032,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model B+"},

	//
	// Compute Module
	//
	rpi_hw_t{
		hwver:          0x11,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Compute Module 1"},
	rpi_hw_t{
		hwver:          0x14,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Compute Module 1"},

	//
	// Pi Zero
	//
	rpi_hw_t{
		hwver:          0x900092,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Pi Zero v1.2"},
	rpi_hw_t{
		hwver:          0x900093,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Pi Zero v1.3"},
	rpi_hw_t{
		hwver:          0x920093,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Pi Zero v1.3"},
	rpi_hw_t{
		hwver:          0x9200c1,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Pi Zero W v1.1"},
	rpi_hw_t{
		hwver:          0x9000c1,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Pi Zero W v1.1"},

	//
	// Model A+
	//
	rpi_hw_t{
		hwver:          0x12,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A+"},
	rpi_hw_t{
		hwver:          0x15,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A+"},
	rpi_hw_t{
		hwver:          0x900021,
		typeNum:        RPI_HWVER_TYPE_PI1,
		periph_base:    PERIPH_BASE_RPI,
		videocore_base: VIDEOCORE_BASE_RPI,
		desc:           "Model A+"},

	//
	// Pi 2 Model B
	//
	rpi_hw_t{
		hwver:          0xa01041,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 2"},
	rpi_hw_t{
		hwver:          0xa01040,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 2"},
	rpi_hw_t{
		hwver:          0xa21041,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 2"},
	//
	// Pi 2 with BCM2837
	//
	rpi_hw_t{
		hwver:          0xa22042,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 2"},
	//
	// Pi 3 Model B
	//
	rpi_hw_t{
		hwver:          0xa02082,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 3"},
	rpi_hw_t{
		hwver:          0xa22082,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Pi 3"},
	//
	// Pi Compute Module 3
	//
	rpi_hw_t{
		hwver:          0xa020a0,
		typeNum:        RPI_HWVER_TYPE_PI2,
		periph_base:    PERIPH_BASE_RPI2,
		videocore_base: VIDEOCORE_BASE_RPI2,
		desc:           "Compute Module 3/L3"},
}

func rpi_hw_detect() (*rpi_hw_t, error) {
	cpuiInfoPath := "/proc/cpuinfo"
	file, err := os.OpenFile(cpuiInfoPath, os.O_RD, 0)
	defer file.Close()
	if err != nil {
		log.Printf("Can't open %v", cpuiInfoPath)
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Can't read all from %v", cpuiInfoPath)
		return nil, err
	}
	var regex = regexp.MustCompile(`Revision.*: (.*)`)

	all := regex.FindSubmatch([]byte("Revision  : 1ab246"))
	if len(all) < 2 {
		err = fmt.Errorf("Can't find revsion number %v", cpuiInfoPath)
		return nil, err
	}
	revString := all[1]
	rev := string(revString[:len(revString)])

	for _, rpi := range rpi_hw_info {
		hwver := rpi_hw_info[i].hwver

		// Take out warranty and manufacturer bits
		hwver &= ^(RPI_WARRANTY_MASK | RPI_MANUFACTURER_MASK)
		rev &= ^(RPI_WARRANTY_MASK | RPI_MANUFACTURER_MASK)

		if rev == hwver {
			return &rpi_hw_info[i], nil
		}
	}
	return nil, fmt.Errorf("couldn't find matching revision for %v in rpi_hw_info", rev)
}
