package main

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/config"
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"Sakura-Pi-Node/pkg/usecase"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tinygo.org/x/bluetooth"
)

var (
	environments *config.Config
)

func main() {
	environments = config.LoadEnvironments()
	log.Print(environments)
	adapter.InitializeServo(*environments)
	//adapter.InitializePasori(*environments)
	adapter.InitializeLed(*environments)
	println("Initialized")
	infra.CreateMQTTClient(environments.TargetIP, func(c mqtt.Client) { subscribeEvents() })
	println("Initialized 2")
	go InitializeSesami()
	println("Initialized 3")
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

	infra.Subscribe(environments.DoorStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorState(environments.DoorStateResponsePath)
	})

	infra.Subscribe(environments.DoorSwitchStateRequestPath, func(message mqtt.Message) {
		usecase.PublishDoorSwitchState(environments.DoorSwitchStateResponsePath)
	})

	infra.Subscribe(environments.KeyStatePath, func(message mqtt.Message) {
		var key entity.KeyState
		err := json.Unmarshal(message.Payload(), &key)
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Println("Received key event. Key State:", key.Open)
		if key.DeviceID != os.Getenv(usecase.DeviceIDIdentifier) {
			return
		}

		usecase.KeyControl(key, environments.DoorSwitchStateResponsePath)
	})

}

func InitializeSesami() {
	// 必要な変数の定義
	privateKey := "xxxx" // テスト用
	//privateKey := "xxxx"  // 家に付着してる方
	//sesame5Address := "xxxx"
	sesame5Address := "xxxx"
	initialRandomCode := []byte{0x01, 0x02, 0x03, 0x04}
	// インフラストラクチャの初期化
	peripheral := adapter.NewBLEPeripheral()

	// ユースケースの初期化
	uc := usecase.NewBluetoothUsecase(peripheral, nil, nil)

	uc.Peripheral.Scan(func(b *bluetooth.Adapter, result bluetooth.ScanResult) {
		// RSSIフィルタ
		if result.RSSI < -60 {
			return
		}

		// ManufacturerDataの確認
		manufacturerData := result.ManufacturerData()
		for _, element := range manufacturerData {
			if element.CompanyID == 0x055A { // SesamiのCompanyID
				if sesame5Address == result.Address.String() {
					fmt.Printf("add %v\n", result.Address.String())
					b.StopScan()
				}

				//sesame5Address = result.Address.String()
				//fmt.Printf("add %v\n", result.Address.String())
				//b.StopScan()
			}
		}
	})

	fmt.Printf("ADDR %v\n", sesame5Address)

	// 接続
	err := uc.Connect(sesame5Address)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer peripheral.Disconnect()

	//// 通知の初期化
	err = uc.InitializeNotify()
	if err != nil {
		fmt.Printf("Failed to initialize notifications: %v\n", err)
		os.Exit(1)
	}
	//
	// ランダムコードの取得（仮）
	uc.RandomCode = initialRandomCode // 実際には通知ハンドラーで更新する必要があります
	//
	// 暗号化/復号化のインスタンスを作成
	encryptor := adapter.NewAESCCMEncryptor([]byte(privateKey), uc.RandomCode)
	decryptor := adapter.NewAESCCMDecryptor([]byte(privateKey), uc.RandomCode)
	uc.Encryptor = encryptor
	uc.Decryptor = decryptor
	uc.StartNotifyHandler()

	time.Sleep(1 * time.Second)

	//
	// トークン生成
	err = uc.GenerateToken(privateKey)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	// ログインコマンドの送信
	err = uc.SendLoginCommand()
	if err != nil {
		fmt.Printf("Failed to send login command: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)

	err = uc.SendUnlockCommand("test")
	if err != nil {
		fmt.Printf("Failed to send Unlock command: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(20 * time.Second)

}

func listenForIDEvents() {
	//for id := range adapter.IDEventChannel {
	//	// IDイベントを受け取った際の処理
	//	log.Println("Received ID event:", id)
	//	usecase.PublishCard(id, environments.CardPath)
	//}
	//for {
	//	fmt.Println("please input: ")
	//	var str string
	//	_, err := fmt.Scan(&str)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(str)
	//}
}

func start() {
	adapter.StartReading()
}
