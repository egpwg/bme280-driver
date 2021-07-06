package driver

import (
	"fmt"
	"log"
	"testing"

	"github.com/egpwg/bme280-driver/internal/driver"
	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

func TestGetDriverInfo(t *testing.T) {
	var tDrv i2c.I2cDriver
	driver.Register(&tDrv)

	info, err := GetDriverInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("driver_info: %v", info)
}
