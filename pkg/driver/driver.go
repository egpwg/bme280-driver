package driver

import (
	"github.com/egpwg/bme280-driver/internal/driver"
)

type DriverInfo struct {
	Driver string
	Bus    []string
}

func GetDriverInfo() (info []DriverInfo, err error) {
	drivers := driver.GetDrivers()

	for k, v := range drivers {
		bus, err := v.Init()
		if err != nil {
			return nil, err
		}

		info = append(info, DriverInfo{
			Driver: k,
			Bus:    bus,
		})
	}

	return
}
