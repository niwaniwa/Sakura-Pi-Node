package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
)

type AESCCMDecryptor struct {
	key            []byte
	randomCode     []byte
	decryptCounter uint64
}

func NewAESCCMDecryptor(key, randomCode []byte) *AESCCMDecryptor {
	return &AESCCMDecryptor{
		key:        key,
		randomCode: randomCode,
	}
}

func (d *AESCCMDecryptor) Decrypt(data []byte) ([]byte, error) {
	// TODO: Implement this
}

func GenerateToken(key, data []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)[:4], nil // 最初の4バイトを使用
}
