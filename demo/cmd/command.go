package cmd

import (
	"fmt"
	"strings"

	"github.com/egpwg/bme280-driver/demo/model"
)

func Command(c string) (bool, error) {
	cList := strings.Split(c, " ")
	switch cList[0] {
	case "":

	case "exit":
		return true, nil
	case "init":
		if err := model.Init(); err != nil {
			return true, err
		}
	case "setmode":
		modErr := fmt.Errorf("Only support weather mode(input 1 or weather)")
		if len(cList) != 2 {
			return false, modErr
		}
		if cList[1] == "1" || cList[1] == "weather" {
			if err := model.SetUserMode(model.UMWeather); err != nil {
				return false, err
			}
		} else {
			return false, modErr
		}
	case "all":
		data, err := model.All()
		if err != nil {
			return false, err
		}
		allVal := data.All
		for _, v := range data.Sequence {
			if v == "Pressure" {
				fmt.Printf("%s : %f\n", v, allVal[v])
				continue
			}
			fmt.Printf("%s : %.2f\n", v, allVal[v])
		}
	case "t":
		data, err := model.Temperature()
		if err != nil {
			return false, err
		}
		fmt.Println(fmt.Sprintf("Temperature: %.2f", data))
	case "h":
		data, err := model.Humidity()
		if err != nil {
			return false, err
		}
		fmt.Println(fmt.Sprintf("Humidity: %.2f", data))
	case "p":
		data, err := model.Pressure()
		if err != nil {
			return false, err
		}
		fmt.Println(fmt.Sprintf("Pressure: %f", data))
	case "reset":
		err := model.Reset()
		if err != nil {
			return false, err
		}
	default:
		fmt.Println("This command is not supported")
	}

	return false, nil
}

func All() error {
	singleInit()
	data, err := model.All()
	if err != nil {
		return err
	}
	allVal := data.All
	for _, v := range data.Sequence {
		if v == "Pressure" {
			fmt.Printf("%s : %f\n", v, allVal[v])
			continue
		}
		fmt.Printf("%s : %.2f\n", v, allVal[v])
	}
	return nil
}

func Temperature() error {
	singleInit()
	t, err := model.Temperature()
	if err != nil {
		return err
	}
	fmt.Printf("Temperature: %.2f\n", t)
	return nil
}

func Humidity() error {
	singleInit()
	t, err := model.Humidity()
	if err != nil {
		return err
	}
	fmt.Printf("Humidity: %.2f\n", t)
	return nil
}

func Pressure() error {
	singleInit()
	t, err := model.Pressure()
	if err != nil {
		return err
	}
	fmt.Printf("Pressure: %f\n", t)
	return nil
}

func singleInit() {
	model.Init()
	// model.Reset()
	model.SetUserMode(model.UMWeather)
}
