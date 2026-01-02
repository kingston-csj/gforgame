package protocol

import (
	"errors"
	"io/github/gforgame/common"
)

// Protocol constants.
const (
	HeadLength    = 12
	MaxPacketSize = 64 * 1024
)

type MessageHeader struct {
	Cmd   int32 //消息类型
	Index int32 //客户端消息索引，由客户端维护自增长
	Size  int32 //消息长度
}

// ErrPacketSizeExceed is the error used for encode/decode.
var ErrPacketSizeExceed = errors.New("protocol: packet size exceed")

type Protocol struct {
	buf *common.ByteBuffer
}

// NewDecoder returns a new decoder that used for decode network bytes slice.
func NewDecoder() *Protocol {
	return &Protocol{
		buf: common.NewByteBuffer(4096, MaxPacketSize),
	}
}

func (c *Protocol) readHeader() (*MessageHeader, error) {
	buff, err := c.buf.Next(HeadLength)
	if err != nil {
		return nil, err
	}
	size := bytesToInt32(buff[0:4])
	index := bytesToInt32(buff[4:8])
	cmd := bytesToInt32(buff[8:HeadLength])

	// packet length limitation
	if int(size) > MaxPacketSize {
		return nil, ErrPacketSizeExceed
	}
	return &MessageHeader{Cmd: cmd, Index: index, Size: size}, nil
}

func (c *Protocol) Decode(data []byte) ([]*Packet, error) {
	c.buf.Write(data)
	// check length
	if c.buf.Len() < HeadLength {
		return nil, errors.New("length too small")
	}
	var packets []*Packet

	for c.buf.Len() > HeadLength {
		// 保存读取索引
		c.buf.MarkReadIndex()
		header, err := c.readHeader()
		if err != nil {
			return packets, err
		}
		// 消息体长度
		bodySize := int(header.Size) - HeadLength
		if bodySize <= c.buf.Len() {
			body, err := c.buf.Next(bodySize)
			if err != nil {
				return packets, err
			}
			p := &Packet{Header: *header, Data: body}
			packets = append(packets, p)
		} else {
			c.buf.ResetReadIndex()
			break
		}
	}

	return packets, nil
}

func (c *Protocol) Encode(cmd int32, index int32, data []byte) ([]byte, error) {
	bodyLen := len(data)
	buf := make([]byte, HeadLength+bodyLen)
	msgLength := HeadLength + bodyLen
	copy(buf[0:4], int32ToBytes(int32(msgLength)))
	copy(buf[4:8], int32ToBytes(index))
	copy(buf[8:HeadLength], int32ToBytes(int32(cmd)))
	copy(buf[HeadLength:], data)

	return buf, nil
}

// Decode packet data length byte to int(Big end)
func bytesToInt32(b []byte) int32 {
	result := int32(0)
	for _, v := range b {
		result = result<<8 + int32(v)
	}
	return result
}

func int32ToBytes(n int32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte((n >> 24) & 0xFF)
	buf[1] = byte((n >> 16) & 0xFF)
	buf[2] = byte((n >> 8) & 0xFF)
	buf[3] = byte(n & 0xFF)
	return buf
}
