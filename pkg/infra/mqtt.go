package infra

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	debugPrefixIdentifier = "prefix"
)

var (
	debugPrefix string
	client      mqtt.Client
)

type MessageListener func(message mqtt.Message)

func CreateMQTTClient(targetIP string) {
	getEnvironmentValues()
	mqtt.DEBUG = log.New(os.Stdout, debugPrefix, 0)
	mqtt.ERROR = log.New(os.Stdout, debugPrefix, 0)

	var defaultFunction mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Println(msg.Topic() + string(msg.Payload()))
	}

	options := mqtt.NewClientOptions().AddBroker(targetIP)
	options.SetKeepAlive(2 * time.Second)
	options.SetPingTimeout(2 * time.Second)
	options.SetDefaultPublishHandler(defaultFunction)
	client = mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func getEnvironmentValues() {
	debugPrefix = os.Getenv(debugPrefixIdentifier)
}

func Publish(topic string, message interface{}, qos byte) {
	client.Publish(topic, qos, false, message)
}

func Subscribe(topic string, listener MessageListener) {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		listener(msg)
	}
	if subscribeToken := client.Subscribe(topic, 0, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		log.Fatal(subscribeToken.Error())
	}
}

func CloseAll() {
	client.Disconnect(1000)
}
