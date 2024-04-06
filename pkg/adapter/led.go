package adapter

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
)

const (
	RedLed   = 13
	GreenLed = 18
)

var (
	redLedPin   rpio.Pin
	greenLedPin rpio.Pin
)

func InitializeLed() {
	err := rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	redLedPin = rpio.Pin(RedLed)
	redLedPin.Output()
	redLedPin.High()

	greenLedPin = rpio.Pin(GreenLed)
	greenLedPin.Output()
	greenLedPin.High()
}

func RedLedToggle() {
	redLedPin.Toggle()
}

func GreenLedToggle() {
	greenLedPin.Toggle()
}
