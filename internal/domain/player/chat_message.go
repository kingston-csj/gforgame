package player

type ChatMessage struct {
	Id string

	Channel int32

	SenderId string

	SenderHead int

	ReceiverId string

	Timestamp int64

	Content string
}
