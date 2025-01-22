package sesami

type Token []byte

type RandomCode []byte

type Command struct {
	OpCode    byte
	ItemCode  byte
	Payload   []byte
	IsEncrypt bool
}
