package main

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/config"
	"Sakura-Pi-Node/pkg/infra"
	"Sakura-Pi-Node/pkg/usecase"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

const (
	CardPath              = "card_publish_path"
	KeyStatePath          = "key_state_publish_path"
	DoorStateRequestPath  = "door_state_request_path"
	DoorStateResponcePath = "door_state_response_path"
)

func main() {
	environments := config.LoadEnvironments()
	log.Print(environments)
	adapter.InitializeServo()
	adapter.InitializePasori()
	infra.InitializeMQTT(environments)

	infra.Subscribe(KeyStatePath, func(message mqtt.Message) {
		usecase.KeyControl(message)
	})
	infra.Subscribe(DoorStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorState(os.Getenv(DoorStateResponcePath))
	})

	go listenForIDEvents()
}

func listenForIDEvents() {
	for id := range adapter.IDEventChannel {
		// IDイベントを受け取った際の処理
		log.Println("Received ID event:", id)
		usecase.PublishCard(id, os.Getenv(CardPath))
	}
}
