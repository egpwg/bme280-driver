package device

import (
	"errors"
	"log"

	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

type Bme280 struct {
	devAddr    uint16
	userMode   UserMode
	sensorMode SensorMode
	filter     FilterCoef
	bus        i2c.Bus
	calib      Calibration
}

// NewDevice new bme280 by bus file descriptor
func NewDevice(bus i2c.Bus) (bme280 *Bme280, err error) {
	bme280 = &Bme280{
		devAddr: DEVICE_ADDR,
		bus:     bus,
	}

	if err = bme280.checkChipID(); err != nil {
		return nil, err
	}

	var (
		tph [0xA2 - 0x88]byte
		h   [0xE8 - 0xE1]byte
	)
	err = bme280.readFromRegister(regAddrCalib1, tph[:])
	if err != nil {
		return nil, err
	}
	err = bme280.readFromRegister(regAddrCalib2, h[:])
	if err != nil {
		return nil, err
	}
	bme280.calib = newCalibration(tph[:], h[:])

	return
}

// checkChipID check bme280 chip id
func (b *Bme280) checkChipID() (err error) {
	var chipId [1]byte
	err = b.readFromRegister(regAddrChipID, chipId[:])
	if err != nil {
		log.Println(err)
		return err
	}

	switch chipId[0] {
	case 0x60:
	default:
		return errors.New("The device is not bme280")
	}

	return nil
}

// SetUserMode set user mode such weather, indoor
func (b *Bme280) SetUserMode(mode int) (err error) {
	b.userMode = UserMode(mode)
	os := b.userMode.GetOversampling()
	ctrlMesg := []byte{regAddrCtrlMeas, byte(os["Temperature"])<<5 | byte(os["Pressure"])<<2 | byte(b.sensorMode)}
	err = b.writeToRegister(ctrlMesg)
	if err != nil {
		log.Println(err)
		return err
	}
	ctrlHum := []byte{regAddrCtrlHum, byte(os["Humidity"])}
	err = b.writeToRegister(ctrlHum)
	if err != nil {
		log.Println(err)
		return err
	}
	config := []byte{regAddrConfig, byte(TSb1000)<<5 | byte(FilterCoefOff)<<2}
	err = b.writeToRegister(config)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type sensorValue struct {
	Temperature float32
	Pressure    float32
	Humidity    float32
}

// GetSenseData get sense all value: temperature, pressure, humidity
func (b *Bme280) GetSenseValue() (data *sensorValue, err error) {
	for {
		status, err := b.checkStatus()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if status {
			break
		}
	}

	var buf [8]byte
	if err = b.readFromRegister(regAddrPress, buf[:]); err != nil {
		log.Println(err)
		return nil, err
	}

	tRaw := int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4
	tFine, t := b.calib.CompensateTemperatureInt32(tRaw)
	data.Temperature = float32(t / 100)

	pRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	p := b.calib.CompensatePressureInt64(tFine, pRaw)
	data.Pressure = float32(p / 256)

	hRaw := int32(buf[6])<<8 | int32(buf[7])
	h := b.calib.CompensateHumidityInt32(tFine, hRaw)
	data.Humidity = float32(h / 1024)

	return
}

// GetTemperatureValue get temperature from sensor register
func (b *Bme280) GetTemperatureValue() (temperature float32, err error) {
	for {
		status, err := b.checkStatus()
		if err != nil {
			log.Println(err)
			return 0, err
		}
		if status {
			break
		}
	}

	var buf [3]byte
	if err = b.readFromRegister(regAddrTemp, buf[:]); err != nil {
		log.Println(err)
		return 0, err
	}

	tRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	_, t := b.calib.CompensateTemperatureInt32(tRaw)
	temperature = float32(t / 100)

	return
}

// GetPressureValue get pressure from sensor register
func (b *Bme280) GetPressureValue() (pressure float32, err error) {
	for {
		status, err := b.checkStatus()
		if err != nil {
			log.Println(err)
			return 0, err
		}
		if status {
			break
		}
	}

	var buf [6]byte
	if err = b.readFromRegister(regAddrPress, buf[:]); err != nil {
		log.Println(err)
		return 0, err
	}

	tRaw := int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4
	tFine, _ := b.calib.CompensateTemperatureInt32(tRaw)

	pRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	p := b.calib.CompensatePressureInt64(tFine, pRaw)
	pressure = float32(p / 256)

	return
}

// GetHumidityValue get humidity from sensor register
func (b *Bme280) GetHumidityValue() (humidity float32, err error) {
	for {
		status, err := b.checkStatus()
		if err != nil {
			log.Println(err)
			return 0, err
		}
		if status {
			break
		}
	}

	var buf [5]byte
	if err = b.readFromRegister(regAddrTemp, buf[:]); err != nil {
		log.Println(err)
		return 0, err
	}

	tRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	tFine, _ := b.calib.CompensateTemperatureInt32(tRaw)

	hRaw := int32(buf[3])<<8 | int32(buf[4])
	h := b.calib.CompensateHumidityInt32(tFine, hRaw)
	humidity = float32(h / 1024)

	return
}

// Reset reset sensor all mode
func (b *Bme280) Reset() (err error) {
	reset := []byte{regAddrReset, 0xB6}
	if err = b.writeToRegister(reset); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) readFromRegister(regAddr uint8, data []byte) (err error) {
	err = b.bus.RdWr(b.devAddr, []byte{regAddr}, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) writeToRegister(regData []byte) (err error) {
	err = b.bus.RdWr(b.devAddr, regData, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) checkStatus() (status bool, err error) {
	var s [1]byte
	err = b.readFromRegister(regAddrStatus, s[:])
	if err != nil {
		log.Println(err)
		return false, err
	}

	return s[0]&0x08 == 0, nil
}
