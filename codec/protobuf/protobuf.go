package protobuf

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

// Codec implements to serialize.ProtobufCodec interface
type Codec struct{}

// NewSerializer returns a new ProtobufCodec.
func NewSerializer() *Codec {
	return &Codec{}
}

// Encode Marshal returns the protobuf encoding of v.
func (s *Codec) Encode(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Decode Unmarshal parses the protobuf-encoded data and stores the result
// in the value pointed to by v.
func (s *Codec) Decode(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
