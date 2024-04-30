package adapter

import (
	"Sakura-Pi-Node/pkg/config"
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
)

var (
	redLedPin   rpio.Pin
	greenLedPin rpio.Pin
)

func InitializeLed(config config.Config) {
	err := rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	redLedPin = rpio.Pin(config.RedLed)
	redLedPin.Output()
	redLedPin.High()

	greenLedPin = rpio.Pin(config.GreenLed)
	greenLedPin.Output()
	greenLedPin.High()
}

func RedLedToggle() {
	redLedPin.Toggle()
}

func GreenLedToggle() {
	greenLedPin.Toggle()
}
