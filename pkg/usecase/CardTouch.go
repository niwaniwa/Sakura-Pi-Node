package usecase

import (
	"Sakura-Pi-Node/pkg/entity"
	"Sakura-Pi-Node/pkg/infra"
	"encoding/json"
	"fmt"
	"time"
)

func PublishCard(id []byte, path string) {
	timestamp := time.Now()
	data := entity.Card{
		ID:        id,
		Timestamp: timestamp,
	}

	// Jsonにしているが基本的に何でもよい
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	infra.Publish(path, jsonData, 0)
}
