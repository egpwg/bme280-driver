package i2c

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/egpwg/bme280-driver/internal/driver"
)

type I2cDriver struct {
	bus []i2cBus
}

func (i *I2cDriver) Name() (name string) {
	return "i2c-driver"
}

func (i *I2cDriver) Init() (bus []string, err error) {
	path := "/dev/i2c-*"
	matches, err := filepath.Glob(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(matches) == 0 {
		err = errors.New("no i2c bus exist")
		log.Println(err)
		return nil, err
	}

	for _, m := range matches {
		if err != nil {
			continue
		}
		n := m[len("/dev/"):]
		bus = append(bus, n)
		ib := i2cBus{
			name: n,
			path: m,
		}
		i.bus = append(i.bus, ib)
	}

	return
}

var i2cDrv I2cDriver

func GetI2cDriver() (drv *I2cDriver) {
	return &i2cDrv
}

func init() {
	err := driver.Register(&i2cDrv)
	if err != nil {
		panic(err)
	}
}
