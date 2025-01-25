package adapter

import (
	"Sakura-Pi-Node/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Sesame interface {
	OpenKey(done chan<- bool)
	CloseKey(done chan<- bool)
	GetKeyState() bool
}

type sesameImpl struct {
	config     *config.Config
	keyState   bool
	keyStateCh chan bool
}

func InitializeSesame(cfg *config.Config) Sesame {
	return &sesameImpl{
		config:     cfg,
		keyState:   false,
		keyStateCh: make(chan bool, 1),
	}
}

func (s *sesameImpl) OpenKey(done chan<- bool) {
	go func() {
		err := s.sendCommand("open")
		if err != nil {
			fmt.Printf("Failed to open key: %v\n", err)
			done <- false
			return
		}
		done <- true
	}()
}

func (s *sesameImpl) CloseKey(done chan<- bool) {
	go func() {
		err := s.sendCommand("close")
		if err != nil {
			fmt.Printf("Failed to close key: %v\n", err)
			done <- false
			return
		}
		done <- true
	}()
}

func (s *sesameImpl) GetKeyState() bool {
	return s.keyState
}

func (s *sesameImpl) sendCommand(command string) error {
	url := fmt.Sprintf("%s/%s", s.config.SesamePythonIp, command)

	var jsonData = []byte("{}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (s *sesameImpl) monitorKeyState() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state, err := s.fetchKeyState()
			if err != nil {
				fmt.Printf("Failed to fetch key state: %v\n", err)
				continue
			}
			s.keyState = state
			select {
			case s.keyStateCh <- state:
			default:
			}
		}
	}
}

func (s *sesameImpl) fetchKeyState() (bool, error) {
	url := fmt.Sprintf("%s/status", s.config.SesamePythonIp)

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	// Parse the JSON response
	var result struct {
		Status  string                 `json:"status"`
		Data    map[string]interface{} `json:"data"`
		Message string                 `json:"message,omitempty"`
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return false, err
	}

	if result.Status != "success" {
		return false, fmt.Errorf("API error: %s", result.Message)
	}

	isLockRange, ok1 := result.Data["is_lock_range"].(bool)
	isUnlockRange, ok2 := result.Data["is_unlock_range"].(bool)
	if !ok1 || !ok2 {
		return false, fmt.Errorf("invalid status data")
	}

	// ロック範囲にある場合は閉錠中、アンロック範囲にある場合は開錠中
	if isLockRange {
		return false, nil // 閉錠中
	}
	if isUnlockRange {
		return true, nil // 開錠中
	}

	// どちらでもない場合は不明として閉錠中と仮定
	return false, nil
}
