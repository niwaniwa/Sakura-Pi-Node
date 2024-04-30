package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	CardPathIdentifier                   = "card_publish_path"
	KeyStatePathIdentifier               = "key_state_publish_path"
	DoorStateRequestPathIdentifier       = "door_state_request_path"
	DoorSwitchStateRequestPathIdentifier = "door_switch_state_request_path"
	DoorStateResponsePathIdentifier      = "door_state_response_path"
	targetIPIdentifier                   = "target_ip"
	PasoriIntervalTimeIdentifier         = "pasori_interval_time"
	DebugPrefixIdentifier                = "debug_prefix"
	PwmPinIdentifier                     = "pwm_pin"
	SwPinIdentifier                      = "sw_pin"
	DoorSwPinIdentifier                  = "door_sw_pin"
	DoorReedSwitchIdentifier             = "door_reed_switch"
	RedLedIdentifier                     = "red_led"
	GreenLedIdentifier                   = "green_led"
)

type Config struct {
	PasoriIntervalTime         int
	CardPath                   string
	KeyStatePath               string
	DoorStateRequestPath       string
	DoorSwitchStateRequestPath string
	DoorStateResponsePath      string
	TargetIP                   string
	DebugPrefix                string
	PwmPin                     int
	SwPin                      int
	DoorSwPin                  int
	DoorReedSwitch             int
	RedLed                     int
	GreenLed                   int
}

func LoadEnvironments() *Config {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	pasoriIntervalTime, _ := strconv.Atoi(os.Getenv(PasoriIntervalTimeIdentifier))
	pwmPin, _ := strconv.Atoi(os.Getenv(PwmPinIdentifier))
	swPin, _ := strconv.Atoi(os.Getenv(SwPinIdentifier))
	doorSwPin, _ := strconv.Atoi(os.Getenv(DoorSwPinIdentifier))
	doorReedSwitch, _ := strconv.Atoi(os.Getenv(DoorReedSwitchIdentifier))
	redLed, _ := strconv.Atoi(os.Getenv(RedLedIdentifier))
	greenLed, _ := strconv.Atoi(os.Getenv(GreenLedIdentifier))

	return &Config{
		PasoriIntervalTime:         pasoriIntervalTime,
		CardPath:                   os.Getenv(CardPathIdentifier),
		KeyStatePath:               os.Getenv(KeyStatePathIdentifier),
		DoorStateRequestPath:       os.Getenv(DoorStateRequestPathIdentifier),
		DoorSwitchStateRequestPath: os.Getenv(DoorSwitchStateRequestPathIdentifier),
		DoorStateResponsePath:      os.Getenv(DoorStateResponsePathIdentifier),
		TargetIP:                   os.Getenv(targetIPIdentifier),
		DebugPrefix:                os.Getenv(DebugPrefixIdentifier),
		PwmPin:                     pwmPin,
		SwPin:                      swPin,
		DoorSwPin:                  doorSwPin,
		DoorReedSwitch:             doorReedSwitch,
		RedLed:                     redLed,
		GreenLed:                   greenLed,
	}

}
