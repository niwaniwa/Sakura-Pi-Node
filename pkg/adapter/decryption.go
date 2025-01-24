package adapter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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

func (e *AESCCMDecryptor) Decrypt(data []byte) ([]byte, error) {
	iv := make([]byte, 8)
	binary.LittleEndian.PutUint64(iv, e.decryptCounter)
	iv = append(iv, e.randomCode...)

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	cipherText, _ := base64.StdEncoding.DecodeString(string(data))
	e.decryptCounter += 1

	if len(cipherText) < aes.BlockSize {
		panic("cipher text must be longer than blocksize")
	} else if len(cipherText)%aes.BlockSize != 0 {
		panic("cipher text must be multiple of blocksize(128bit)")
	}

	cipherText = cipherText[aes.BlockSize:]
	plainText := make([]byte, len(cipherText))

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plainText, cipherText)
	return UnPadByPkcs7(plainText), nil
}

func UnPadByPkcs7(data []byte) []byte {
	padSize := int(data[len(data)-1])
	return data[:len(data)-padSize]
}

func GenerateToken(key, data []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)[:4], nil // 最初の4バイトを使用
}
