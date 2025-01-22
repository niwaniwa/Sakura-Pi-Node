package usecase

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/entity/sesami"
	"errors"
	"fmt"
)

type BluetoothUsecase struct {
	Peripheral adapter.Peripheral
	Encryptor  sesami.Encryptor
	Decryptor  sesami.Decryptor
	Token      sesami.Token
	RandomCode sesami.RandomCode
}

func NewBluetoothUsecase(peripheral adapter.Peripheral, encryptor sesami.Encryptor, decryptor sesami.Decryptor) *BluetoothUsecase {
	return &BluetoothUsecase{
		Peripheral: peripheral,
		Encryptor:  encryptor,
		Decryptor:  decryptor,
	}
}

func (u *BluetoothUsecase) Connect(address string) error {
	fmt.Println("Connecting to peripheral...")
	err := u.Peripheral.Connect(address)
	if err != nil {
		return err
	}
	fmt.Println("Connected.")
	return nil
}

func (u *BluetoothUsecase) InitializeNotify() error {
	enableNotifyData := []byte{0x01, 0x00}
	fmt.Printf("Writing enable notify data: %x\n", enableNotifyData)
	err := u.Peripheral.WriteCharacteristic(0x0010, enableNotifyData, true)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for random code from notifications...")
	notified, err := u.Peripheral.WaitForNotifications(5.0)
	if err != nil {
		return err
	}
	if !notified {
		return errors.New("random code not received")
	}
	fmt.Println("Random code received.")
	return nil
}

func (u *BluetoothUsecase) GenerateToken(privateKey string) error {
	if len(u.RandomCode) == 0 {
		return errors.New("random code is not set")
	}
	token, err := adapter.GenerateToken([]byte(privateKey), u.RandomCode)
	if err != nil {
		return err
	}
	u.Token = token
	return nil
}

func (u *BluetoothUsecase) SendLoginCommand() error {
	if len(u.Token) < 4 {
		return errors.New("invalid token length")
	}
	loginCommand := []byte{0x02, u.Token[0], u.Token[1], u.Token[2], u.Token[3]}
	fmt.Printf("Sending login command: %x\n", loginCommand)
	return u.Peripheral.WriteCharacteristic(0x000D, loginCommand, false)
}

func (u *BluetoothUsecase) SendUnlockCommand(unlockMessage string) error {
	unlockTag := []byte(unlockMessage)
	unlockCommand := append([]byte{0x53, byte(len(unlockTag))}, unlockTag...)
	fmt.Printf("Sending unlock command: %x\n", unlockCommand)
	return u.Peripheral.WriteCharacteristic(0x000D, unlockCommand, true)
}

func (u *BluetoothUsecase) SendLockCommand(unlockMessage string) error {
	unlockTag := []byte(unlockMessage)
	unlockCommand := append([]byte{0x52, byte(len(unlockTag))}, unlockTag...)
	fmt.Printf("Sending unlock command: %x\n", unlockCommand)
	return u.Peripheral.WriteCharacteristic(0x000D, unlockCommand, true)
}
