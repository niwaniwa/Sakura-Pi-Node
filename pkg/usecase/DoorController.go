package usecase

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/config"
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const DeviceIDIdentifier = "device_id"

// リードSwitchを使用したもの。ドアが物理的に開いているか閉まっているかを記述
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

// Deprecated: 鍵ボックス内にあるSwitchがオンかどうか(現在使用不可)
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

// Deprecated: 鍵ボックス内にあるSwitchがオンかどうか(現在使用不可)
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

	isChanged, _ := strconv.ParseBool(os.Getenv(config.ChangedKeyDirection))

	open := key.Open

	if isChanged {
		open = !open
	}

	if open {
		go adapter.OpenKey(done)
	} else {
		go adapter.CloseKey(done)
	}
	result := <-done
	log.Println("Key process ", result)
	PublishDoorSwitchStateCustom(publishPath, key.Open)
}

// Sesameによる鍵制御。
func KeyControlBySesame(key entity.KeyState, publishPath string, sesame adapter.Sesame) {
	done := make(chan bool)
	open := key.Open

	if open {
		sesame.OpenKey(done)
	} else {
		sesame.CloseKey(done)
	}
	result := <-done
	log.Println("Key process ", result)
}
