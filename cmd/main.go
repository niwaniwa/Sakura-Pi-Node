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

var (
	environments *config.Config
)

func main() {
	environments = config.LoadEnvironments()
	log.Print(environments)
	adapter.InitializeServo(*environments)
	adapter.InitializePasori()
	adapter.InitializeLed(*environments)
	infra.CreateMQTTClient(environments.TargetIP)

	subscribeEvents()

	go listenForIDEvents()

	start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	go func() {
		<-signals
		log.Println("\nCtrl+C pressed. Shutting down...")
		infra.CloseAll()
		os.Exit(0)
	}()

	select {}

}

func subscribeEvents() {
	infra.Subscribe(environments.KeyStatePath, func(message mqtt.Message) {
		var key entity.KeyState
		err := json.Unmarshal(message.Payload(), &key)
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Println("Received key event. Key State:", key.Open)
		usecase.KeyControl(key)
	})

	infra.Subscribe(environments.DoorStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorState(environments.DoorStateResponsePath)
	})

	infra.Subscribe(environments.DoorSwitchStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorSwitchState(environments.DoorSwitchStateRequestPath)
	})

}

func listenForIDEvents() {
	for id := range adapter.IDEventChannel {
		// IDイベントを受け取った際の処理
		log.Println("Received ID event:", id)
		usecase.PublishCard(id, environments.CardPath)
	}
}

func start() {
	adapter.StartReading()
}
