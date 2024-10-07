package protocol

import (
	"bytes"
	"errors"
)

// Protocol constants.
const (
	HeadLength    = 12
	MaxPacketSize = 64 * 1024
)

type MessageHeader struct {
	Cmd   int //消息类型
	Index int //客户端消息索引，由客户端维护自增长
	Size  int //消息长度
}

// ErrPacketSizeExceed is the error used for encode/decode.
var ErrPacketSizeExceed = errors.New("protocol: packet size exceed")

type Protocol struct {
	buf *bytes.Buffer
}

// NewDecoder returns a new decoder that used for decode network bytes slice.
func NewDecoder() *Protocol {
	return &Protocol{
		buf: bytes.NewBuffer(nil),
	}
}

func (c *Protocol) readHeader() (*MessageHeader, error) {
	buff := c.buf.Next(HeadLength)
	id := bytesToInt(buff[0:4])
	index := bytesToInt(buff[4:8])
	size := bytesToInt(buff[8:HeadLength])

	// packet length limitation
	if size > MaxPacketSize {
		return nil, ErrPacketSizeExceed
	}
	return &MessageHeader{Cmd: id, Index: index, Size: size}, nil
}

func (c *Protocol) Decode(data []byte) ([]*Packet, error) {
	c.buf.Write(data)
	// check length
	if c.buf.Len() < HeadLength {
		return nil, errors.New("length too small")
	}
	var packets []*Packet

	for c.buf.Len() > 0 {
		header, err := c.readHeader()
		if err != nil {
			return packets, err
		}

		if header.Size <= c.buf.Len() {
			body := c.buf.Next(header.Size)
			p := &Packet{Header: *header, Data: body}
			packets = append(packets, p)
		} else {
			break
		}
	}

	return packets, nil
}

func (c *Protocol) Encode(cmd int, index int, data []byte) ([]byte, error) {
	// p := &Packet{Cmd: cmd, Length: len(data)}
	bodyLen := len(data)
	buf := make([]byte, HeadLength+bodyLen)

	copy(buf[0:4], intToBytes(cmd))
	copy(buf[4:8], intToBytes(index))
	copy(buf[8:HeadLength], intToBytes(bodyLen))
	copy(buf[HeadLength:], data)

	return buf, nil
}

// Decode packet data length byte to int(Big end)
func bytesToInt(b []byte) int {
	result := 0
	for _, v := range b {
		result = result<<8 + int(v)
	}
	return result
}

func intToBytes(n int) []byte {
	buf := make([]byte, 4)
	buf[0] = byte((n >> 24) & 0xFF)
	buf[1] = byte((n >> 16) & 0xFF)
	buf[2] = byte((n >> 8) & 0xFF)
	buf[3] = byte(n & 0xFF)
	return buf
}
