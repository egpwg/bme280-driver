package device

import (
	"fmt"
	"log"
	"testing"

	intI2c "github.com/egpwg/bme280-driver/internal/driver/i2c"
	pkgI2c "github.com/egpwg/bme280-driver/pkg/driver/i2c"
)

func TestSetUserMode(t *testing.T) {
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

	var ctlHum [1]byte
	err = dev.readFromRegister(regAddrCtrlHum, ctlHum[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ctlHum: ", ctlHum)

	var config [1]byte
	err = dev.readFromRegister(regAddrConfig, config[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config: ", config)

	var ctlMeas [1]byte
	err = dev.readFromRegister(regAddrCtrlMeas, ctlMeas[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ctlMeas: ", ctlMeas)
}

func TestGetSenseValue(t *testing.T) {
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

	data, err := dev.GetSenseValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature: %f, Pressure: %f, Humidity: %f", data.Temperature, data.Pressure, data.Humidity)
}

func TestGetTemperatureValue(t *testing.T) {
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

	temp, err := dev.GetTemperatureValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature: %.2f", temp)
}

func TestGetPressureValue(t *testing.T) {
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

	pres, err := dev.GetPressureValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pressure: %.2f", pres)
}

func TestGetHumidityValue(t *testing.T) {
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

	humi, err := dev.GetHumidityValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Humitidy: %.2f", humi)
}

func TestReset(t *testing.T) {
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

	err = dev.Reset()
	if err != nil {
		log.Fatal(err)
	}

	var ctlHum [1]byte
	err = dev.readFromRegister(regAddrCtrlHum, ctlHum[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ctlHum: ", ctlHum)

	var config [1]byte
	err = dev.readFromRegister(regAddrConfig, config[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config: ", config)

	var ctlMeas [1]byte
	err = dev.readFromRegister(regAddrCtrlMeas, ctlMeas[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ctlMeas: ", ctlMeas)
}
