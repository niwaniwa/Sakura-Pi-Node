package infra

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

var targetIP = "mosquitto:1883"

func TestConnect(t *testing.T) {
	CreateMQTTClient(targetIP)
	CloseAll()
}

func TestMQTTPublishAndSubscribe(t *testing.T) {
	CreateMQTTClient(targetIP)

	defer CloseAll()

	// Subscribe
	receivedMessage := make(chan string)

	Subscribe("test/topic", func(message mqtt.Message) {
		receivedMessage <- string(message.Payload())
	})

	// Publish
	testMessage := "Hello MQTT"
	Publish("test/topic", testMessage, 0)

	select {
	case message := <-receivedMessage:
		if message != testMessage {
			t.Errorf("Expected message %q but got %q", testMessage, message)
		}
	case <-time.After(time.Second * 5):
		t.Error("Timeout waiting for the message")
	}
}
