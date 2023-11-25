package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	CardPathIdentifier              = "card_publish_path"
	KeyStatePathIdentifier          = "key_state_publish_path"
	DoorStateRequestPathIdentifier  = "door_state_request_path"
	DoorStateResponcePathIdentifier = "door_state_response_path"
	targetIPIdentifier              = "target_ip"
	PasoriIntervalTimeIdentifier    = "Pasori_Interval_Time"
	DebugPrefixIdentifier           = "debug_prefix"
)

type Config struct {
	PasoriIntervalTime    int
	CardPath              string
	KeyStatePath          string
	DoorStateRequestPath  string
	DoorStateResponcePath string
	TargetIP              string
	DebugPrefix           string
}

func LoadEnvironments() *Config {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	pasoriIntervalTime, _ := strconv.Atoi(os.Getenv(PasoriIntervalTimeIdentifier))

	return &Config{
		PasoriIntervalTime:    pasoriIntervalTime,
		CardPath:              os.Getenv(CardPathIdentifier),
		KeyStatePath:          os.Getenv(KeyStatePathIdentifier),
		DoorStateRequestPath:  os.Getenv(DoorStateRequestPathIdentifier),
		DoorStateResponcePath: os.Getenv(DoorStateResponcePathIdentifier),
		TargetIP:              os.Getenv(targetIPIdentifier),
		DebugPrefix:           os.Getenv(DebugPrefixIdentifier),
	}

}
