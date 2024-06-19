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

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
	fmt.Println("published door state: ", jsonData)
}

func PublishDoorSwitchState(path string) {
	timestamp := time.Now()
	data := entity.DoorSwitchState{
		IsOpen:    adapter.GetDoorSwitchState(),
		Timestamp: timestamp,
		DeviceID:  os.Getenv(DeviceIDIdentifier),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}
	infra.Publish(path, jsonData, 0)
	fmt.Println("published door switch state: ", jsonData)
}

func PublishDoorSwitchStateCustom(path string, isOpen bool) {
	timestamp := time.Now()
	data := entity.DoorSwitchState{
		IsOpen:    isOpen,
		Timestamp: timestamp,
		DeviceID:  os.Getenv(DeviceIDIdentifier),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}
	infra.Publish(path, jsonData, 0)
	fmt.Println("published door switch state: ", jsonData)
}

func KeyControl(key entity.KeyState, publishPath string) {
	done := make(chan bool)
	if key.Open {
		go adapter.OpenKey(done)
	} else {
		go adapter.CloseKey(done)
	}
	result := <-done
	log.Println("Key process ", result)
	PublishDoorSwitchStateCustom(publishPath, key.Open)
}
