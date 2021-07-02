package i2c

import (
	"errors"
	"fmt"
	"log"
	"unsafe"

	"github.com/egpwg/bme280-driver/internal/driver/util"
)

type Bus interface {
	Name() (name string)
	RdWr(addr uint16, w, r []byte) (err error)
	Close() (err error)
}

type i2cBus struct {
	name string
	path string
	file util.File
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
	start := 1
	msg := [2]i2cMsg{}
	if len(w) != 0 {
		msg[0].addr = addr
		msg[0].flags = flagWR
		msg[0].length = uint16(len(w))
		msg[0].buf = uintptr(unsafe.Pointer(&w[0]))
	}
	if len(r) != 0 {
		start = 2
		msg[1].addr = addr
		msg[1].flags = flagRD
		msg[1].length = uint16(len(r))
		msg[1].buf = uintptr(unsafe.Pointer(&r[0]))
	}

	data := i2cRdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msg)),
		nmsgs: start,
	}

	ep := i.file.Ioctl(uintptr(ioctlRdwr), uintptr(unsafe.Pointer(&data)))
	if ep != 0 {
		return ep
	}

	return nil
}

func (i *i2cBus) Close() (err error) {
	return i.file.Close()
}

func Open(name string) (bus *i2cBus, err error) {
	for i, b := range i2cDrv.bus {
		if b.Name() == name {
			err = i2cDrv.bus[i].file.Open(b.path)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			return &i2cDrv.bus[i], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("no bus named %s", name))
}
