package i2c

import (
	"log"

	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

func Open(name string) (bus i2c.Bus, err error) {
	bus, err = i2c.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return
}
