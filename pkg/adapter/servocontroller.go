package adapter

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"time"
)

var (
	managePWMPin    rpio.Pin
	manageDoorSwPin rpio.Pin
	manageSwPin     rpio.Pin
	manageReedPin   rpio.Pin
	isOpen          bool = true
	motorRunning    bool = false
)

const (
	PwmPin           = 13
	SwPin            = 18
	DoorSwPin        = 21
	DoorReedSwitch   = 26
	StopPosition     = 1520 // サーボモーターを停止させるPWMパルス幅(マイクロ秒)
	ForwardPosition  = 800  // サーボモーターを正転させるPWMパルス幅(マイクロ秒)
	ReversePosition  = 2200 // サーボモーターを反転させるPWMパルス幅(マイクロ秒)
	IgnoreSwitchTime = 500  // スイッチ判定を無視する時間 (ミリ秒)
	timeout          = 2500 // 応答がなかった場合にタイムアウトして処理を終了する時間 (ミリ秒)
)

func InitializeServo() {
	err := rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	managePWMPin = rpio.Pin(PwmPin) // PWM setup
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 100)
	managePWMPin.Low()

	manageSwPin = rpio.Pin(SwPin)
	manageSwPin.Input()
	manageSwPin.PullUp()

	manageDoorSwPin = rpio.Pin(DoorSwPin)
	manageDoorSwPin.Input()
	manageDoorSwPin.PullUp()

	manageReedPin = rpio.Pin(DoorReedSwitch)
	manageReedPin.Input()
	manageReedPin.PullUp()
}

func OpenKey(done chan<- bool) {
	managePWMPin.High()

	motorStartTime := time.Now()
	motorRunning = true
	setServo(managePWMPin, float64(ForwardPosition))

	for {
		if motorRunning && time.Since(motorStartTime) > IgnoreSwitchTime*time.Millisecond {
			if manageSwPin.Read() == rpio.Low {
				position := StopPosition
				motorRunning = false
				setServo(managePWMPin, float64(position))
				break
			}
		}

		// timeout
		if motorRunning && time.Since(motorStartTime) > timeout*time.Millisecond {
			position := StopPosition
			motorRunning = false
			setServo(managePWMPin, float64(position))
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	managePWMPin.Low()
	RedLedToggle()
	GreenLedToggle()
	isOpen = true
	done <- true
}

func CloseKey(done chan<- bool) {
	managePWMPin.High()
	motorStartTime := time.Now()
	motorRunning = true
	setServo(managePWMPin, float64(ReversePosition))

	for {
		if motorRunning && time.Since(motorStartTime) > IgnoreSwitchTime*time.Millisecond {
			if manageSwPin.Read() == rpio.Low {
				position := StopPosition
				motorRunning = false
				setServo(managePWMPin, float64(position))
				break
			}
		}

		// timeout
		if motorRunning && time.Since(motorStartTime) > timeout*time.Millisecond {
			position := StopPosition
			motorRunning = false
			setServo(managePWMPin, float64(position))
			break
		}

		time.Sleep(10 * time.Millisecond)
	}
	managePWMPin.Low()
	RedLedToggle()
	GreenLedToggle()
	isOpen = false
	done <- true
}

func GetKeyState() bool {
	return isOpen
}

// true = close, false = open
func GetDoorState() bool {
	return manageReedPin.Read() == 0
}

func GetDoorSwitchState() bool {
	return manageDoorSwPin.Read() == 0
}

// 指定したパルス幅でサーボモーターを制御
func setServo(pin rpio.Pin, pulseWidthMicroSeconds float64) {
	pulseWidthFraction := pulseWidthMicroSeconds / 20000
	dutyCycle := uint32(pulseWidthFraction * 1000)

	pin.DutyCycle(dutyCycle, 1000)
}
