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
	err := i2cDrv.Init()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len: ", len(i2cDrv.bus))

	file, err := Open(i2cDrv.bus[0].name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var chip [1]byte
	err = i2cDrv.bus[0].RdWr(bme280, []byte{chipId}, chip[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("chip id: ", chip[0])
}
