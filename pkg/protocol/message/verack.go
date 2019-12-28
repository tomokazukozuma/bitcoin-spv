package message

type Verack struct{}

func (v *Verack) Command() [12]byte {
	var commandName [12]byte
	copy(commandName[:], "verack")
	return commandName
}

func (v *Verack) Encode() []byte {
	return []byte{}
}
