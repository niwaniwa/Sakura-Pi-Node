package adapter

import (
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

type Peripheral interface {
	Connect(address string) error
	Disconnect() error
	WriteCharacteristic(handle bluetooth.UUID, data []byte) error
	WriteNotification(handle bluetooth.UUID) error
	WaitForNotifications(timeout float64) ([]byte, error)
	Scan(callback func(*bluetooth.Adapter, bluetooth.ScanResult)) error
	Notifications() <-chan []byte
}

type BLEPeripheral struct {
	adapter  *bluetooth.Adapter
	device   bluetooth.Device
	notifyCh chan []byte
}

func NewBLEPeripheral() Peripheral {
	return &BLEPeripheral{
		notifyCh: make(chan []byte, 10), // Buffered channel to prevent blocking
	}
}

func (p *BLEPeripheral) Scan(callback func(*bluetooth.Adapter, bluetooth.ScanResult)) error {
	p.adapter = bluetooth.DefaultAdapter
	p.must("enable BLE stack", p.adapter.Enable())
	err := p.adapter.Scan(callback)
	if err != nil {
		return err
	}
	return nil
}

func (p *BLEPeripheral) Connect(address string) error {

	p.adapter = bluetooth.DefaultAdapter

	rawAddress, err := bluetooth.ParseMAC(address)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	targetAddress := bluetooth.Address{
		MACAddress: bluetooth.MACAddress{
			MAC: rawAddress,
		},
	}
	targetAddress.SetRandom(true)

	// Enable BLE interface.
	p.must("enable BLE stack", p.adapter.Enable())

	fmt.Println("address " + targetAddress.MAC.String())

	p.device, err = p.adapter.Connect(targetAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	fmt.Println("connected to ", p.device.Address.String())

	p.onConnected(p.device)

	return nil
}

func (p *BLEPeripheral) Disconnect() error {
	if &p.device == nil {
		return errors.New("no device connected")
	}
	fmt.Println("Disconnecting from peripheral...")
	if err := p.device.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	fmt.Println("Disconnected from peripheral.")
	return nil
}

func (p *BLEPeripheral) WriteCharacteristic(uuid bluetooth.UUID, data []byte) error {
	characteristic, err := p.getCharacteristicByHandle(uuid)
	if err != nil {
		return err
	}

	fmt.Println("Write without response...")

	_, err = characteristic.WriteWithoutResponse(data)

	if err != nil {
		return fmt.Errorf("failed to write characteristic: %w", err)
	}

	return nil
}

func (p *BLEPeripheral) WriteNotification(handle bluetooth.UUID) error {
	characteristic, err := p.getCharacteristicByHandle(handle)
	if err != nil {
		return err
	}
	err = characteristic.EnableNotifications(p.callback)
	if err != nil {
		return fmt.Errorf("failed to Enable Notifications characteristic: %w", err)
	}
	return nil
}

func (p *BLEPeripheral) WaitForNotifications(timeout float64) ([]byte, error) {
	select {
	case data := <-p.notifyCh:
		return data, nil
	case <-time.After(time.Duration(timeout * float64(time.Second))):
		return nil, errors.New("timeout waiting for notifications")
	}
}

// onConnected コールバック
func (p *BLEPeripheral) onConnected(peripheral bluetooth.Device) {
	fmt.Println("Peripheral connecting:", p.device.Address.String())
}

// onDisconnected コールバック
func (p *BLEPeripheral) onDisconnected(peripheral bluetooth.Device, reason uint8) {
	fmt.Println("Peripheral disconnected:", peripheral.Address.String(), "Reason:", reason)
	close(p.notifyCh)
}

// getCharacteristicByHandle ハンドルから特性を取得
func (p *BLEPeripheral) getCharacteristicByHandle(handle bluetooth.UUID) (*bluetooth.DeviceCharacteristic, error) {

	services, err := p.device.DiscoverServices([]bluetooth.UUID{bluetooth.New16BitUUID(uint16(0xfd81))})
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	fmt.Printf("service length %v\n", len(services))

	for _, service := range services {
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to discover characteristics for service %s: %w", service.UUID().String(), err)
		}
		//fmt.Println("=============================")
		//fmt.Printf("index %v,Now Service\t%v\ntarget service\t\t%v\n", i, service.UUID().String(), bluetooth.New16BitUUID(uint16(0xfd81)))
		//fmt.Printf("target chara\t%v\n", handle)
		//
		//for _, char := range characteristics {
		//	fmt.Printf("       chara\t%v\n", char.UUID().String())
		//}

		for _, char := range characteristics {
			if char.UUID() == handle {

				//err := char.EnableNotifications(callback)
				//// If error, bail out
				//if err != nil {
				//	fmt.Printf("error enabling notifications: %v\n", err)
				//	return nil, err
				//}

				return &char, nil
			}
		}
	}

	return nil, fmt.Errorf("characteristic with handle %v not found", handle)
}

func (p *BLEPeripheral) Notifications() <-chan []byte {
	return p.notifyCh
}

func (p *BLEPeripheral) callback(data []byte) {
	// Non-blocking send to prevent blocking the BLE stack
	select {
	case p.notifyCh <- data:
		// Data sent successfully
	default:
		fmt.Println("Notification channel is full, discarding data")
	}
}

func (p *BLEPeripheral) must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
