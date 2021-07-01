package device

import (
	"errors"
	"log"

	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

type Bme280 struct {
	DevAddr    uint16
	UserMode   UserMode
	SensorMode SensorMode
	Filter     FilterCoef
	Bus        i2c.Bus
	calib      Calibration
}

// NewDevice new bme280 by bus file descriptor
func NewDevice(bus i2c.Bus) (bme280 *Bme280, err error) {
	testB := Bme280{DevAddr: DEVICE_ADDR}
	if err = testB.checkChipID(); err != nil {
		return nil, err
	}

	return &Bme280{
		DevAddr: DEVICE_ADDR,
		Bus:     bus,
	}, nil
}

// checkChipID check bme280 chip id
func (b *Bme280) checkChipID() (err error) {
	var chipId [1]byte
	err = b.readFromRegister(b.DevAddr, 0xD0, chipId[:])
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

func (b *Bme280) SetUserMode(mode UserMode) (err error) {
	var (
		tph [0xA2 - 0x88]byte
		h   [0xE8 - 0xE1]byte
	)
	err = b.readFromRegister(b.DevAddr, 0x88, tph[:])
	if err != nil {
		log.Println(err)
		return err
	}
	err = b.readFromRegister(b.DevAddr, 0xE1, h[:])
	if err != nil {
		log.Println(err)
		return err
	}
	b.calib = newCalibration(tph[:], h[:])

	b.UserMode = mode
	os := b.UserMode.GetOversampling()
	ctrlMesg := []byte{0xF4, byte(os["Temperature"])<<5 | byte(os["Pressure"])<<2 | byte(b.SensorMode)}
	err = b.writeToRegister(b.DevAddr, ctrlMesg)
	if err != nil {
		log.Println(err)
		return err
	}
	ctrlHum := []byte{0xF2, byte(os["Humidity"])}
	err = b.writeToRegister(b.DevAddr, ctrlHum)
	if err != nil {
		log.Println(err)
		return err
	}
	config := []byte{0xF5, byte(TSb1000)<<5 | byte(FilterCoefOff)<<2}
	err = b.writeToRegister(b.DevAddr, config)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type sensorData struct {
	Temperature float32
	Pressure    float32
	Humidity    float32
}

func (b *Bme280) GetSenseData() (data *sensorData, err error) {
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
	if err = b.readFromRegister(b.DevAddr, 0xF7, buf[:]); err != nil {
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

func (b *Bme280) Reset() (err error) {
	reset := []byte{0xE0, 0xB6}
	if err = b.writeToRegister(b.DevAddr, reset); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) readFromRegister(devAddr uint16, regAddr uint8, data []byte) (err error) {
	err = b.Bus.RdWr(devAddr, []byte{regAddr}, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) writeToRegister(devAddr uint16, regData []byte) (err error) {
	err = b.Bus.RdWr(devAddr, regData, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (b *Bme280) checkStatus() (status bool, err error) {
	var s [1]byte
	err = b.readFromRegister(b.DevAddr, 0xF3, s[:])
	if err != nil {
		log.Println(err)
		return false, err
	}

	return s[0]&0x08 == 0, nil
}
