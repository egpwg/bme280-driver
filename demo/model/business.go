package model

import (
	"fmt"

	"github.com/egpwg/bme280-driver/pkg/device"
	"github.com/egpwg/bme280-driver/pkg/driver"
	"github.com/egpwg/bme280-driver/pkg/driver/i2c"
)

var globalBME280 *device.Bme280
var globalBus i2c.Bus

type SensorValue struct {
	Sequence []string
	All      map[string]float32
}

func Init() error {
	// 获取所有驱动信息，如i2c驱动名及其所有总线名
	// i2c驱动名固定为i2c-driver，总线名固定为i2c-X(如i2c-0、i2c-1)
	info, err := driver.GetDriverInfo()
	if err != nil {
		return err
	}

	var name string
	for _, v := range info {
		if v.Driver == "i2c-driver" {
			name = v.Bus[0]
		}
	}

	// 选择某个总线并打开其文件获取文件描述符，用完需关闭描述符
	bus, err := i2c.Open(name)
	if err != nil {
		return err
	}
	globalBus = bus
	// defer bus.Close()

	// 通过文件描述符来创建bme280对象
	bme280, err := device.NewDevice(bus)
	if err != nil {
		return err
	}
	setGlobalBME280(bme280)
	return nil
}

func SetUserMode(m int) error {
	// 设置用户模式：Weather、Indoor、HumiSensing、Gaming
	// 目前只支持Weather模式
	// 传入数字代表不同模式，1为Weather，2为HumiSensing，3为Indoor，4为Gaming
	if m < 1 || m > 4 {
		return fmt.Errorf("Please set mode in[1:Weather,2:Humisensing,3:Indoor,4:Gaming]!")
	}
	err := GlobalBME280().SetUserMode(device.UserMode(m))
	if err != nil {
		return err
	}
	return nil
}

func setGlobalBME280(g *device.Bme280) {
	globalBME280 = g
}

func GlobalBME280() *device.Bme280 {
	return globalBME280
}

func All() (*SensorValue, error) {
	// 获取传感器温度、压力、湿度
	value, err := GlobalBME280().GetSenseValue()
	if err != nil {
		return nil, err
	}

	seq := [3]string{"Temperature", "Humidity", "Pressure"}

	r := make(map[string]float32)
	r["Temperature"] = value.Temperature
	r["Pressure"] = value.Pressure
	r["Humidity"] = value.Humidity

	s := &SensorValue{seq[:], r}

	return s, nil
}

func Temperature() (float32, error) {
	// 获取传感器温度
	t, err := GlobalBME280().GetTemperatureValue()
	if err != nil {
		return 0, err
	}
	return t, nil
}

func Pressure() (float32, error) {
	// 获取传感器温度
	p, err := GlobalBME280().GetPressureValue()
	if err != nil {
		return 0, err
	}
	return p, nil
}

func Humidity() (float32, error) {
	// 获取传感器温度
	h, err := GlobalBME280().GetHumidityValue()
	if err != nil {
		return 0, err
	}
	return h, nil
}

func Reset() error {
	// 重置传感器
	err := GlobalBME280().Reset()
	if err != nil {
		return err
	}
	return nil
}

func CloseBus() error {
	if globalBus != nil {
		return globalBus.Close()
	}
	return nil
}
