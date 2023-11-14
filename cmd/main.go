package main

import (
	"Sakura-Hardware/pkg/api"
	"Sakura-Hardware/pkg/config"
	"Sakura-Hardware/pkg/device"
	"log"
)

func main() {
	config := config.LoadEnvironments()
	log.Print(config)
	device.InitializeServo()
	device.InitializePasori()
	api.InitializeApiServer()

	go listenForIDEvents()
}

func listenForIDEvents() {
	for id := range device.IDEventChannel {
		// IDイベントを受け取った際の処理
		log.Println("Received ID event:", id)
	}
}
