package protocol

type Packet struct {
	Header MessageHeader

	Data []byte
}

type RequestDataFrame struct {
	Header MessageHeader

	Msg any
}
