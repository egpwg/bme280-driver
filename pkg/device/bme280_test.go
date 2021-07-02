package device

import (
	"fmt"
	"log"
	"testing"

	intI2c "github.com/egpwg/bme280-driver/internal/driver/i2c"
	pkgI2c "github.com/egpwg/bme280-driver/pkg/driver/i2c"
)

func TestNewDevice(t *testing.T) {
	i2cDrv := intI2c.GetI2cDriver()
	_, err := i2cDrv.Init()
	if err != nil {
		log.Fatal(err)
	}

	bus, err := pkgI2c.Open("i2c-1")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	dev, err := NewDevice(bus)
	if err != nil {
		log.Fatal(err)
	}

	err = dev.SetUserMode(1)
	if err != nil {
		log.Fatal(err)
	}

	data, err := dev.GetSenseValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature: %.2f, Pressure: %.2f, Humidity: %.2f", data.Temperature, data.Pressure, data.Humidity)
}