package main

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/config"
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"Sakura-Pi-Node/pkg/usecase"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	subscribeEvents()

	go listenForIDEvents()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	go func() {
		<-signals
		fmt.Println("Ctrl+C pressed. Shutting down...")
		infra.CloseAll()
		os.Exit(0)
	}()

	select {}

}

func subscribeEvents() {
	infra.Subscribe(KeyStatePath, func(message mqtt.Message) {
		var key entity.KeyState
		err := json.Unmarshal(message.Payload(), &key)
		if err != nil {
			fmt.Println(err)
			return
		}
		usecase.KeyControl(key)
	})

	infra.Subscribe(DoorStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorState(os.Getenv(DoorStateResponcePath))
	})
}

func listenForIDEvents() {
	for id := range adapter.IDEventChannel {
		// IDイベントを受け取った際の処理
		log.Println("Received ID event:", id)
		usecase.PublishCard(id, os.Getenv(CardPath))
	}
}
