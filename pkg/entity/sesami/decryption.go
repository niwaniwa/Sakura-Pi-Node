package sesami

type Decryptor interface {
	Decrypt(data []byte) ([]byte, error)
}
