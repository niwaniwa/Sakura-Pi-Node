package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	PasoriIntervalTime int
	GearAngle          int
}

func LoadEnvironments() *Config {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	pasoriIntervalTime, _ := strconv.Atoi(os.Getenv("PasoriIntervalTime"))
	gearAngle, _ := strconv.Atoi(os.Getenv("GearAngle"))

	return &Config{
		PasoriIntervalTime: pasoriIntervalTime,
		GearAngle:          gearAngle,
	}

}
