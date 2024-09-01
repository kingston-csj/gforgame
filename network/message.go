package network

type Packet struct {
	Header MessageHeader

	Data []byte
}

// create a NewPacket Packet instance.
func NewPacket() *Packet {
	return &Packet{}
}

type RequestDataFrame struct {
	Header MessageHeader

	Msg any
}
