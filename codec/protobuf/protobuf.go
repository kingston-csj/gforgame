package protobuf

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

// ProtobufCodec implements the serialize.ProtobufCodec interface
type ProtobufCodec struct{}

// NewSerializer returns a new ProtobufCodec.
func NewSerializer() *ProtobufCodec {
	return &ProtobufCodec{}
}

// Marshal returns the protobuf encoding of v.
func (s *ProtobufCodec) Encode(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal parses the protobuf-encoded data and stores the result
// in the value pointed to by v.
func (s *ProtobufCodec) Decode(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
