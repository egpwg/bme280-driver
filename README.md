# bme280-driver

## Struct

![arch](https://user-images.githubusercontent.com/21075532/117534305-574dea00-b023-11eb-8d25-84accb84062f.png)

## 接口调用实例

```go
package main

import (
    "fmt"
    "log"
    
	"github.com/egpwg/bme280-driver/pkg/driver"
    "github.com/egpwg/bme280-driver/pkg/device"
    "github.com/egpwg/bme280-driver/pkg/driver/i2c"
)

func main() {
    // 获取所有驱动信息，如i2c驱动名及其所有总线名
    // i2c驱动名固定为i2c-driver，总线名固定为i2c-X(如i2c-0、i2c-1)
    info, err := driver.GetDriverInfo()
    if err != nil {
        log.Fatal(err)
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
        log.Fatal(err)
    }
    defer bus.Close()
    
    // 通过文件描述符来创建bme280对象
    bme280, err := device.NewDevice(bus)
    if err != nil {
        log.Fatal(err)
    }
    
    // 设置用户模式：Weather、Indoor、HumiSensing、Gaming
    // 传入类型代表不同模式，UMWeather为Weather，UMHumiSensing为HumiSensing，UMIndoor为Indoor，UMGaming为Gaming
    err := bme280.SetUserMode(1)
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取传感器温度、压力、湿度
    value, err := bme280.GetSenseValue()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Temperature: %.2f, Pressure: %.2f, Humidity: %.2f", value.Temperature, value.Pressure, value.Humidity)
    
    // 获取传感器温度
    temperature, err := bme280.GetTemperatureValue()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Temperature: %.2f", temperature)
    
    // 获取传感器压力
    pressure, err := bme280.GetPressureValue()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Pressure: %.2f", pressure)
    
    // 获取传感器湿度
    humidity, err := bme280.GetHumidityValue()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Humidity: %.2f", humidity)
    
    // 重置传感器
    err = bme280.Reset()
    if err != nil {
        log.Fatal(err)
    }
}
```

