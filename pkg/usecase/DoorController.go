package usecase

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const DeviceIDIdentifier = "device_id"

func PublishDoorState(path string) {
	timestamp := time.Now()
	data := entity.DoorState{
		IsOpen:    adapter.GetDoorState(),
		Timestamp: timestamp,
		DeviceID:  os.Getenv(DeviceIDIdentifier),
	}

	// Jsonにしているが基本的に何でもよい
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
}

func PublishDoorSwitchState(path string) {
	timestamp := time.Now()
	data := entity.DoorSwitchState{
		IsOpen:    adapter.GetDoorSwitchState(),
		Timestamp: timestamp,
		DeviceID:  os.Getenv(DeviceIDIdentifier),
	}

	// Jsonにしているが基本的に何でもよい
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
}

func KeyControl(key entity.KeyState) {
	done := make(chan bool)
	if key.Open {
		go adapter.OpenKey(done)
	} else {
		go adapter.CloseKey(done)
	}
	<-done
	log.Println("Key process ", done)
}
