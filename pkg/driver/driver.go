package driver

import (
	"errors"
	"fmt"
	"log"

	"github.com/egpwg/bme280-driver/internal/driver"
)

// Init input driver name then init，such as i2c-driver
func Init(name string) (err error) {
	drivers := driver.GetDrivers()
	if _, ok := drivers[name]; !ok {
		err = errors.New(fmt.Sprintf("The driver %s is not exist", name))
		log.Println(err)
		return err
	}

	err = drivers[name].Init()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
