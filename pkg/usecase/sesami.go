package usecase

import (
	"Sakura-Pi-Node/pkg/adapter"
	"Sakura-Pi-Node/pkg/entity/sesami"
	"errors"
	"fmt"

	"tinygo.org/x/bluetooth"
)

type BluetoothUsecase struct {
	Peripheral    adapter.Peripheral
	Encryptor     sesami.Encryptor
	Decryptor     sesami.Decryptor
	Token         sesami.Token
	RandomCode    sesami.RandomCode
	notifyHandler func([]byte)
	buffer        []byte
}

var (
	WriteCharacteristicUUID bluetooth.UUID
)

func NewBluetoothUsecase(peripheral adapter.Peripheral, encryptor sesami.Encryptor, decryptor sesami.Decryptor) *BluetoothUsecase {
	WriteCharacteristicUUID, _ = bluetooth.ParseUUID("16860002-a5ae-9856-b6d3-dbb4c676993e")
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
	return nil
}

func (u *BluetoothUsecase) InitializeNotify() error {
	//data := []byte{0x01, 0x00}
	handle, err := bluetooth.ParseUUID("16860003-a5ae-9856-b6d3-dbb4c676993e")
	if err != nil {
		return err
	}
	err = u.Peripheral.WriteNotification(handle)
	if err != nil {
		return err
	}
	fmt.Printf("end on connected\n")
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
	fmt.Printf("token %x\n", u.Token)
	loginCommand := []byte{0x02, u.Token[0], u.Token[1], u.Token[2], u.Token[3]}
	fmt.Printf("Sending login command: %x\n", loginCommand)
	//return u.Peripheral.WriteCharacteristic(WriteCharacteristicUUID, loginCommand)
	return u.Send(WriteCharacteristicUUID, loginCommand, false)
}

func (u *BluetoothUsecase) SendUnlockCommand(unlockMessage string) error {
	unlockTag := []byte(unlockMessage)
	unlockCommand := append([]byte{0x53, byte(len(unlockTag))}, unlockTag...)
	fmt.Printf("Sending unlock command: %x\n", unlockCommand)
	return u.Peripheral.WriteCharacteristic(WriteCharacteristicUUID, unlockCommand)
}

func (u *BluetoothUsecase) SendLockCommand(lockMessage string) error {
	lockTag := []byte(lockMessage)
	unlockCommand := append([]byte{0x52, byte(len(lockTag))}, lockTag...)
	fmt.Printf("Sending unlock command: %x\n", unlockCommand)
	uuid := bluetooth.New16BitUUID(uint16(0x000D))
	return u.Peripheral.WriteCharacteristic(uuid, unlockCommand)
}

func (u *BluetoothUsecase) Send(uuid bluetooth.UUID, sendData []byte, isEncrypt bool) error {
	if isEncrypt {
		if len(u.Token) <= 0 {
			println("Error: Token is not set")
			return errors.New("error: Token is not set")
		}
		var err error
		sendData, err = u.Encryptor.Encrypt(sendData)
		if err != nil {
			return fmt.Errorf("Error: %e\n", err)
		}
	}

	remain := len(sendData)
	offset := 0

	for remain > 0 {
		header := byte(0)
		if offset == 0 {
			header += 1 // 最初のパケット
		}

		var buffer []byte
		if remain <= 19 {
			buffer = sendData[offset:]
			remain = 0
			if isEncrypt {
				header += 4
			} else {
				header += 2
			}
		} else {
			buffer = sendData[offset : offset+19]
			offset += 19
			remain -= 19
		}

		buffer = append([]byte{header}, buffer...)

		fmt.Printf("rawdata:\t%x\n", sendData)
		fmt.Printf("bufdata:\t%x\n", buffer)

		err := u.Peripheral.WriteCharacteristic(uuid, buffer)

		if err != nil {
			return err
		}
	}
	return nil

}

func (u *BluetoothUsecase) SetNotifyHandler(handler func([]byte)) {
	u.notifyHandler = handler
}

func (u *BluetoothUsecase) StartNotifyHandler() {
	// Notificationsチャネルを取得
	notifications := u.Peripheral.Notifications()

	// 通知を待ち受けるゴルーチンを開始
	go func() {
		for data := range notifications {
			// fmt.Printf("Received notification: %x\n", data)
			u.HandleNotification(data)
		}
	}()
}

// 通知ハンドラーの設定
func (u *BluetoothUsecase) HandleNotification(data []byte) {
	fmt.Printf("Handling notification data: %x\n", data)
	// データの解析と処理をここに実装
	// 例: ランダムコードの取得など
	if len(data) < 2 {
		fmt.Println("Invalid data length.")
		return
	}

	if len(data) <= 0 {
		return
	}

	if (data[0] & 1) != 0 {
		u.buffer = make([]byte, 0)
		fmt.Println("buffer reset")
	}

	u.buffer = append(u.buffer, data[1:]...)
	headerType := data[0] >> 1

	if headerType == 0 {
		println("Continuation of data.")
		return
	}

	var decryptedData []byte
	if headerType == 2 {
		decryptedData, err := u.Decryptor.Decrypt(u.buffer)

		if err != nil {
			fmt.Errorf("Decryption failed : %e\n", err)
			return
		}

		fmt.Printf("Decrypted data: %x\n", decryptedData)
	} else {
		decryptedData = u.buffer
		fmt.Printf("Received data: %x\n", decryptedData)
	}

	if len(decryptedData) < 2 {
		fmt.Println("Invalid data length.")
		return
	}

	opCode := decryptedData[0]
	itemCode := decryptedData[1]

	fmt.Printf("Op Code: %d, Item Code: %d\n", opCode, itemCode)

	if itemCode == 0x0E { // 例: ランダムコード取得のオペコード
		if len(decryptedData) < 6 {
			fmt.Println("Insufficient data for random code.")
			return
		}
		u.RandomCode = decryptedData[2:6]
		fmt.Printf("Random Code Updated: %x\n", u.RandomCode)
	}
}
