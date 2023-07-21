package tuntap

import (
	"syscall"
	"unsafe"
)

type DeviceType int

const (
	_ DeviceType = iota
	TUN
	TAP

	IFF_MULTI_QUEUE = 0x0100
)

type Config struct {
	Type       DeviceType
	Name       string
	MultiQueue bool
	Persist    bool
}

type IfReq struct {
	Name  [16]byte
	Flags uint16
	pad   [22]byte
}

type Device struct {
	Fd int
	C  Config
}

func New(c Config) (*Device, error) {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}

	err = setup(fd, c)
	if err != nil {
		return nil, err
	}

	return &Device{
		Fd: fd,
		C:  c,
	}, nil
}

func setup(fd int, c Config) error {
	ifreq := IfReq{}

	copy(ifreq.Name[:], c.Name)
	ifreq.Flags = syscall.IFF_NO_PI
	if c.Type == TUN {
		ifreq.Flags = syscall.IFF_TUN
	} else {
		ifreq.Flags = syscall.IFF_TAP
	}
	if c.MultiQueue {
		ifreq.Flags |= IFF_MULTI_QUEUE
	}
	err := ioctl(uintptr(fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq)))
	if err != nil {
		return err
	}

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(socket)

	ifreq.Flags |= syscall.IFF_UP | syscall.IFF_RUNNING
	err = ioctl(uintptr(socket), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq)))
	if err != nil {
		return err
	}

	value := 0
	if c.Persist {
		value = 1
	}
	return ioctl(uintptr(fd), syscall.TUNSETPERSIST, uintptr(value))
}

func ioctl(fd, request, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	if errno != 0 {
		return errno
	}
	return nil
}

func (d *Device) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return syscall.Read(d.Fd, p)
}

func (d *Device) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return syscall.Write(d.Fd, p)
}

func (d *Device) Close() error {
	return syscall.Close(d.Fd)
}
