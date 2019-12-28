package message

type Message interface {
	Command() [12]byte
	Encode() []byte
}
