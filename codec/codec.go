package codec

type MessageCodec interface {

	// Encode 将bean序列化为byte数组
	Encode(v any) ([]byte, error)

	// Decode 将byte数组反序列化为bean
	Decode(data []byte, v any) error
}
