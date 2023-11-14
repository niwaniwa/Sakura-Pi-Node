package device

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"time"
)

var (
	managePWMPin rpio.Pin
	manageMosPin rpio.Pin
	manageSwPin  rpio.Pin
	isOpen       bool = true
	lock         bool = false
)

const (
	PwmPin = 13
	MosPin = 17
	SwPin  = 18
)

func InitializeServo() {
	err := rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manageMosPin = rpio.Pin(MosPin) // MOS SEIGYO OUT PUT PIN
	manageMosPin.Output()
	manageMosPin.Low()
	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 100)
	managePWMPin.Low()

	manageSwPin = rpio.Pin(SwPin)
	manageSwPin.Input()
	manageSwPin.PullUp()
}

func BlockCurrent() {
	manageMosPin.High()
}

func OpeningCurrent() {
	manageMosPin.Low()
}

func OpenKey() {
	managePWMPin.High()
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	managePWMPin.Low()
	isOpen = true
}

func CloseKey() {
	managePWMPin.High()
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(50-i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	managePWMPin.Low()
	isOpen = false
}

func GetKeyState() bool {
	return isOpen
}
