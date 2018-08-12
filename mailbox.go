package rpiws2811

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// Relevant source -> https://www.raspberrypi.org/app/uploads/2012/02/BCM2835-ARM-Peripherals.pdf

// **** <mailbox.h> ****

// TODO @jmbarzee
// #include <linux/ioctl.h>

const (
	MAJOR_NUM   = 100
	DEV_MEM     = "/dev/mem"
	DEV_GPIOMEM = "/dev/gpiomem"

	// #include <linux/ioccom.h>
	IOCPARM_MASK = 0x1fff

	IOC_OUT   = uint32(0x40000000)
	IOC_IN    = uint32(0x80000000)
	IOC_INOUT = (IOC_IN | IOC_OUT)
	// #include </linux/ioccom.h>
)

var ( // TODO @jmbarzee const
	IOCTL_MBOX_PROPERTY = _IOWR(MAJOR_NUM, 0, uint32(0))
)

// #include <linux/ioccom.h>
func _IOWR(g uint32, n uint32, t uint32) uint32 {
	return _IOC(IOC_INOUT, (g), (n), uint32(unsafe.Sizeof(t)))
}
func _IOC(inout uint32, group uint32, num uint32, len uint32) uint32 {
	return (inout | ((len & IOCPARM_MASK) << 16) | ((group) << 8) | (num))
}

// #include </linux/ioccom.h>

// **** </mailbox.h> ****
// **** <mailbox.c> ****

func mapmem(base uint32, memLength uintptr, mem_dev string) unsafe.Pointer {
	offsetmask := uint32(os.Getpagesize() - 1)
	pagemask := ^uint32(0) ^ offsetmask
	var mem_fd uintptr

	file, err := os.OpenFile(mem_dev, os.O_RDWR|os.O_SYNC, 0)
	defer file.Close()
	if err != nil {
		log.Printf("Can't open %v", mem_dev)
		return nil
	}
	mem_fd = file.Fd()

	bytes, err := syscall.Mmap(
		int(mem_fd),
		int64(base&pagemask),
		int(memLength),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		log.Printf("mmap error %v", err)
		return nil
	}

	// return a pointer to the new memory at an offset of (base & offsetmask) * sizeof byte
	offset := uintptr((base & offsetmask))
	return unsafe.Pointer(uintptr(unsafe.Pointer(&bytes[0])) + offset)
}

func unmapmem(addr unsafe.Pointer, size uintptr) {
	offsetmask := uint32(os.Getpagesize() - 1)
	pagemask := ^uint32(0) ^ offsetmask
	baseaddr := uintptr(uint32(uintptr(addr)) & pagemask)

	sh := &reflect.SliceHeader{
		Data: baseaddr,
		Len:  int(size),
		Cap:  int(size),
	}
	mem := *(*[]byte)(unsafe.Pointer(sh))

	err := syscall.Munmap(mem)
	if err != nil {
		log.Printf("mmap error %v", err)
	}
}

/*
 * use ioctl to send mbox property message
 */
// TODO @jmbarzee static
func mbox_property(file *os.File, buf unsafe.Pointer) error {

	if file == nil || file.Fd() < 0 {
		file, err := mbox_open()
		defer file.Close()
		if err != nil {
			log.Printf("mbox_property open file failed: %v\n", err)
			return err
		}
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(IOCTL_MBOX_PROPERTY), uintptr(buf))
	if err < 0 {
		return fmt.Errorf("ioctl_set_msg failed: %v\n", err)
	}
	return nil
}

func mem_alloc(file *os.File, size uint32, align uint32, flags uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x3000c // (the tag id)
	p[3] = 12      // (size of the buffer)
	p[4] = 12      // (size of the data)
	p[5] = size    // (num bytes? or pages?)
	p[6] = align   // (alignment)
	p[7] = flags   // (MEM_FLAG_L1_NONALLOCATING)

	p[8] = 0x00000000                      // end tag
	p[0] = 9 * uint32(unsafe.Sizeof(p[0])) // actual size
	err := mbox_property(file, unsafe.Pointer(&p[0]))
	if err != nil {
		return 0
	}
	return p[5] // TODO @jmbarzee why are these all returning numbers that the caller has? Are they being modified?
}

func mem_free(file *os.File, handle uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x3000f // (the tag id)
	p[3] = 4       // (size of the buffer)
	p[4] = 4       // (size of the data)
	p[5] = handle

	p[6] = 0x00000000                      // end tag
	p[0] = 7 * uint32(unsafe.Sizeof(p[0])) // actual size

	mbox_property(file, unsafe.Pointer(&p))
	// TODO @jmbarze error check

	return p[5]
}

// TODO @jmbarzee deal with strange error handling
func mem_lock(file *os.File, handle uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x3000d // (the tag id)
	p[3] = 4       // (size of the buffer)
	p[4] = 4       // (size of the data)
	p[5] = handle

	p[6] = 0x00000000                      // end tag
	p[0] = 7 * uint32(unsafe.Sizeof(p[0])) // actual size

	err := mbox_property(file, unsafe.Pointer(&p))
	if err != nil {
		return ^uint32(0)
		// TODO @jmbarze wtf is this return for
	}
	return p[5]
}

func mem_unlock(file *os.File, handle uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x3000e // (the tag id)
	p[3] = 4       // (size of the buffer)
	p[4] = 4       // (size of the data)
	p[5] = handle

	p[6] = 0x00000000                      // end tag
	p[0] = 7 * uint32(unsafe.Sizeof(p[0])) // actual size

	mbox_property(file, unsafe.Pointer(&p))
	// TODO @jmbarze error check

	return p[5]
}

// TODO @jmbarzee tripple check this shit. Its crazy
func execute_code(file *os.File, code uint32, r0 uint32, r1 uint32,
	r2 uint32, r3 uint32, r4 uint32, r5 uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x30010 // (the tag id)
	p[3] = 28      // (size of the buffer)
	p[4] = 28      // (size of the data)
	p[5] = code
	p[6] = r0
	p[7] = r1
	p[8] = r2
	p[9] = r3
	p[10] = r4
	p[11] = r5

	p[12] = 0x00000000                      // end tag
	p[0] = 13 * uint32(unsafe.Sizeof(p[0])) // actual size

	mbox_property(file, unsafe.Pointer(&p))
	// TODO @jmbarze error check

	return p[5]
}

func qpu_enable(file *os.File, enable uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request

	p[2] = 0x30012 // (the tag id)
	p[3] = 4       // (size of the buffer)
	p[4] = 4       // (size of the data)
	p[5] = enable

	p[6] = 0x00000000                      // end tag
	p[0] = 7 * uint32(unsafe.Sizeof(p[0])) // actual size

	mbox_property(file, unsafe.Pointer(&p))
	// TODO @jmbarze error check

	return p[5]
}

func execute_qpu(file *os.File, num_qpus uint32, control uint32,
	noflush uint32, timeout uint32) uint32 {
	p := make([]uint32, 32)

	p[0] = 0          // size
	p[1] = 0x00000000 // process request
	p[2] = 0x30011    // (the tag id)
	p[3] = 16         // (size of the buffer)
	p[4] = 16         // (size of the data)
	p[5] = num_qpus
	p[6] = control
	p[7] = noflush
	p[8] = timeout // ms

	p[9] = 0x00000000                       // end tag
	p[0] = 10 * uint32(unsafe.Sizeof(p[0])) // actual size

	mbox_property(file, unsafe.Pointer(&p))
	// TODO @jmbarze error check

	return p[5]
}

// **** <makedev.c> ****
// REF: https://github.com/lattera/glibc/blob/master/sysdeps/unix/sysv/linux/makedev.c
type dev_t uint64

func gnu_dev_major(dev dev_t) uint32 {
	return (uint32(dev>>8) & 0xfff) | (uint32(dev>>32) & ^uint32(0xfff))
}

func gnu_dev_minor(dev dev_t) uint32 {
	return uint32(dev&0xff) | (uint32(dev>>12) & ^uint32(0xff))
}

func gnu_dev_makedev(major uint32, minor uint32) dev_t {
	first := uint64(minor&0xff) | (uint64(major&0xfff) << 8)
	second := uint64(minor & ^uint32(0xff)) << 12
	third := uint64(major & ^uint32(0xfff)) << 32
	return dev_t(first | second | third)
}

// **** </makedev.c> ****

func mbox_open() (*os.File, error) {

	file, err := os.OpenFile("/dev/vcio", 0, 0)
	if file.Fd() >= 0 {
		return file, nil
	}

	// open a char device file used for communicating with kernel mbox driver
	filename := fmt.Sprintf("/tmp/mailbox-%d", os.Getpid())
	syscall.Unlink(filename)
	err = syscall.Mknod(filename, syscall.S_IFCHR|0600, int(gnu_dev_makedev(100, 0)))
	if err != nil {
		log.Printf("Failed to create mailbox device\n")
		return nil, err
	}
	file, err = os.OpenFile(filename, 0, 0)
	if err != nil {
		log.Printf("Can't open device file: %v", err)
		syscall.Unlink(filename)
		return nil, err
	}
	syscall.Unlink(filename)

	return file, nil
}

// **** </mailbox.c> ****
