package i2c

import (
	"errors"
	"fmt"
	"log"
	"unsafe"

	"github.com/egpwg/bme280-driver/internal/common"
)

type Bus interface {
	Name() (name string)
	RdWr(addr uint16, w, r []byte) (err error)
}

type i2cBus struct {
	name  string
	path  string
	ioctl common.Ioctler
}

func (i *i2cBus) Name() (name string) {
	return i.name
}

type i2cRdwrIoctlData struct {
	msgs  uintptr
	nmsgs int
}

type i2cMsg struct {
	addr   uint16
	flags  uint16
	length uint16
	buf    uintptr
}

const (
	flagRD = 0x0001
	flagWR = 0x0000
)

const (
	ioctlRdwr = 0x707
)

func (i *i2cBus) RdWr(addr uint16, w, r []byte) (err error) {
	msg := []i2cMsg{}
	if len(w) != 0 {
		msg[0].addr = addr
		msg[0].flags = flagWR
		msg[0].length = uint16(len(w))
		msg[0].buf = uintptr(unsafe.Pointer(&w[0]))
	}
	if len(r) != 0 {
		msg[1].addr = addr
		msg[1].flags = flagRD
		msg[1].length = uint16(len(r))
		msg[1].buf = uintptr(unsafe.Pointer(&r[0]))
	}

	start := 2
	data := i2cRdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msg)),
		nmsgs: start,
	}

	ep := i.ioctl.Ioctl(uintptr(ioctlRdwr), uintptr(unsafe.Pointer(&data)))
	if ep != 0 {
		return ep
	}

	return nil
}

func Open(name string) (file *common.File, err error) {
	for _, b := range i2cDrv.bus {
		if b.Name() == name {
			err = file.Open(b.path)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			return file, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("no bus named %s", name))
}
