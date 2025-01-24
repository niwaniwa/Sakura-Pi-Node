package adapter

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
	iv := make([]byte, 8)
	binary.LittleEndian.PutUint64(iv, e.encryptCounter)
	iv = append(iv, e.randomCode...)

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(data)) // cipher text must be larger than plaintext
	cbcEncrypt := cipher.NewCBCEncrypter(block, iv)
	e.encryptCounter += 1
	cbcEncrypt.CryptBlocks(cipherText[aes.BlockSize:], data)
	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherText)
	return []byte(cipherTextBase64), nil
}
