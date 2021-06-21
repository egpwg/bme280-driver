package i2c

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/egpwg/bme280-driver/internal/driver"
)

type i2cDriver struct {
	bus []i2cBus
}

func (i *i2cDriver) Name() (name string) {
	return "i2c-driver"
}

func (i *i2cDriver) Init() (err error) {
	path := "/dev/i2c-*"
	matches, err := filepath.Glob(path)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(matches) == 0 {
		err = errors.New("no i2c bus exist")
		log.Println(err)
		return err
	}

	for _, m := range matches {
		if err != nil {
			continue
		}
		ib := i2cBus{
			name: m[len("/dev/"):],
			path: m,
		}
		i.bus = append(i.bus, ib)
	}

	return nil
}

var i2cDrv i2cDriver

func init() {
	err := driver.Register(&i2cDrv)
	if err != nil {
		panic(err)
	}
}
