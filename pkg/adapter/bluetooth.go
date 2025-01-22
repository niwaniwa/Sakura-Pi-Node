package adapter

import (
	"fmt"
	"time"
)

type Peripheral interface {
	Connect(address string) error
	Disconnect() error
	WriteCharacteristic(handle uint16, data []byte, withResponse bool) error
	WaitForNotifications(timeout float64) (bool, error)
}

type BLEPeripheral struct {
	// 実際の接続情報を保持
}

func NewBLEPeripheral() Peripheral {
	return &BLEPeripheral{}
}

func (p *BLEPeripheral) Connect(address string) error {
	// 実際の接続処理を実装
	fmt.Printf("Connecting to %s...\n", address)
	time.Sleep(2 * time.Second) // 擬似的な遅延
	fmt.Println("Connected.")
	return nil
}

func (p *BLEPeripheral) Disconnect() error {
	// 実際の切断処理を実装
	fmt.Println("Disconnected from peripheral.")
	return nil
}

func (p *BLEPeripheral) WriteCharacteristic(handle uint16, data []byte, withResponse bool) error {
	// 実際の書き込み処理を実装
	fmt.Printf("Writing to handle 0x%04X: %x (withResponse=%v)\n", handle, data, withResponse)
	return nil
}

func (p *BLEPeripheral) WaitForNotifications(timeout float64) (bool, error) {
	// 実際の通知待機処理を実装
	fmt.Printf("Waiting for notifications for %.1f seconds...\n", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)
	// 擬似的に通知を受け取ったと仮定
	return true, nil
}
