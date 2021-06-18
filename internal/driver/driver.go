package driver

import (
	"errors"
	"fmt"
)

type Driver interface {
	Name() (name string)
	Init() (err error)
}

var drivers = map[string]Driver{}

func Register(d Driver) (err error) {
	n := d.Name()
	if _, ok := drivers[n]; ok {
		err = errors.New(fmt.Sprintf("The driver %s was registerd", n))
		return err
	}

	drivers[n] = d

	return nil
}

func GetDrivers() (drv map[string]Driver) {
	return drivers
}
