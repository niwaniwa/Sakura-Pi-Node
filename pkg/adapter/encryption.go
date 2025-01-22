package adapter

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

type AESCCMEncryptor struct {
	key            []byte
	randomCode     []byte
	encryptCounter uint64
}

func NewAESCCMEncryptor(key, randomCode []byte) *AESCCMEncryptor {
	return &AESCCMEncryptor{
		key:        key,
		randomCode: randomCode,
	}
}

func (e *AESCCMEncryptor) Encrypt(data []byte) ([]byte, error) {

	iv := encryptCounter.toBytes()

	// uint64をbyte配列に変換
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, encryptCounter)

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	encrypt = make([]byte, len(data))
	cbcEncrypter := cipher.NewCBCEncrypter(block, e.randomCode)
}
