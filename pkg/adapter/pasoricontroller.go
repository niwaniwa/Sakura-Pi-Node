package adapter

import (
	"Sakura-Pi-Node/pkg/config"
	"github.com/bamchoh/pasori"
	"time"
)

var (
	allowRead      = false
	IDEventChannel = make(chan []byte)
)

const (
	VID uint16 = 0x054C // SONY
	PID uint16 = 0x06C1 // RC-S380
)

func InitializePasori(config config.Config) {
	go continuouslyReadID(config.PasoriIntervalTime)
}

func GetID() ([]byte, error) {
	return pasori.GetID(VID, PID)
}

func StartReading() {
	allowRead = true
}

func EndReading() {
	allowRead = false
}

func PublishID(id []byte) {
	IDEventChannel <- id
}

func continuouslyReadID(interval int) {
	for {
		if allowRead {
			id, err := GetID()
			if err != nil {
				continue
			}

			IDEventChannel <- id
		}

		// ここで少し待機することでCPUの負荷を下げる
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
