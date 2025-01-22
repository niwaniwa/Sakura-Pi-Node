package sesami

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
}
