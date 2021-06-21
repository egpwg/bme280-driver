package i2c

import (
	"log"

	"github.com/egpwg/bme280-driver/internal/common"
	"github.com/egpwg/bme280-driver/internal/driver/i2c"
)

func Open(name string) (file *common.File, err error) {
	file, err = i2c.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}
