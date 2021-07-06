package i2c

import (
	"fmt"
	"log"
	"testing"
)

const (
	bme280 = 0x77
	chipId = 0xD0
)

func TestI2cRead(t *testing.T) {
	_, err := i2cDrv.Init()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len: ", len(i2cDrv.bus))

	bus, err := Open(i2cDrv.bus[0].name)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.file.Close()

	var chip [1]byte
	err = i2cDrv.bus[0].RdWr(bme280, []byte{chipId}, chip[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("chip id: ", chip[0])
}

func TestI2cWrite(t *testing.T) {
	_, err := i2cDrv.Init()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len: ", len(i2cDrv.bus))

	bus, err := Open(i2cDrv.bus[0].name)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.file.Close()

	ctrlMesg := []byte{0xF2, byte(0x01)}
	err = i2cDrv.bus[0].RdWr(0x77, ctrlMesg, nil)
	if err != nil {
		log.Fatal(err)
	}
}
