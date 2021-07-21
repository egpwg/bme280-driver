package i2c

import (
	"log"

	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

type Bus interface {
	Name() (name string)
	RdWr(addr uint16, w, r []byte) (err error)
	Close() (err error)
}

func Open(name string) (bus Bus, err error) {
	bus, err = i2c.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return
}
