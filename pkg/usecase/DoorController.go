package usecase

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func PublishKeyState(id []byte, path string) {
	timestamp := time.Now()
	data := entity.KeyState{
		Open:      adapter.GetKeyState(),
		Timestamp: timestamp,
	}

	// Jsonにしているが基本的に何でもよい
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
}

func PublishDoorState(path string) {
	timestamp := time.Now()
	data := entity.DoorState{
		IsOpen:    adapter.GetDoorState(),
		Timestamp: timestamp,
	}

	// Jsonにしているが基本的に何でもよい
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
}

func KeyControl(state mqtt.Message) {
	var key entity.KeyState

	err := json.Unmarshal([]byte(state.Payload()), &key)
	if err != nil {
		fmt.Println(err)
		return
	}

	adapter.OpeningCurrent()
	done := make(chan bool)
	if key.Open {
		go adapter.OpenKey(done)
	} else {
		go adapter.CloseKey(done)
	}
	<-done
	adapter.BlockCurrent()

}
