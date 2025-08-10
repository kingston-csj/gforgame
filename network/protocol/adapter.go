package protocol

import (
	"encoding/json"
	"fmt"
)

// ProtocolType 协议类型
type ProtocolType int

// 网络消息的编解码分为两步：
// 1. 私有协议数据类型的编解码，例如包头，包体的协议编解码
// 2. 具体消息体的编解码，例如ReqLogin,ResLogin这种消息可以使用json编解码,也可使用二进制编解码
// 这里的类型，特指私有协议数据类型，并不是指具体消息的实际表示编解码
const (
	ProtocolTypeBinary ProtocolType = iota // 二进制协议
	ProtocolTypeJSON                       // JSON文本协议,主要用于websocket的文本协议
)

// ProtocolAdapter 协议适配器接口
type ProtocolAdapter interface {
	Decode(data []byte) ([]*Packet, error)
	Encode(cmd int, index int, data []byte) ([]byte, error)
}

// BinaryProtocolAdapter 二进制协议适配器
type BinaryProtocolAdapter struct {
	*Protocol
}

func NewBinaryProtocolAdapter() *BinaryProtocolAdapter {
	return &BinaryProtocolAdapter{
		Protocol: NewDecoder(),
	}
}

// Decode 实现ProtocolAdapter接口
func (b *BinaryProtocolAdapter) Decode(data []byte) ([]*Packet, error) {
	return b.Protocol.Decode(data)
}

// Encode 实现ProtocolAdapter接口
func (b *BinaryProtocolAdapter) Encode(cmd int, index int, data []byte) ([]byte, error) {
	return b.Protocol.Encode(cmd, index, data)
}

// JSONProtocolAdapter JSON协议适配器
type JSONProtocolAdapter struct {
	*Protocol
}

func NewJSONProtocolAdapter() *JSONProtocolAdapter {
	return &JSONProtocolAdapter{
		Protocol: NewDecoder(),
	}
}

// WebSocketJsonFrame JSON格式的数据包
type WebSocketJsonFrame struct {
	Type  string      `json:"$type,omitempty"` // 消息类型标识
	Cmd   int         `json:"cmd"`             // 消息类型
	Index int         `json:"index"`           // 客户端消息索引
	Msg   interface{} `json:"msg,omitempty"`   // 消息数据
	Data  interface{} `json:"data,omitempty"`  // 兼容data字段
}

// Decode 解码JSON格式的数据
func (j *JSONProtocolAdapter) Decode(data []byte) ([]*Packet, error) {
	var jsonPacket WebSocketJsonFrame
	if err := json.Unmarshal(data, &jsonPacket); err != nil {
		return nil, fmt.Errorf("unmarshal json failed: %v", err)
	}

	// 将消息数据转换为字节数组
	var dataBytes []byte
	var msgData interface{}

	// 优先使用Msg字段，如果没有则使用Data字段
	if jsonPacket.Msg != nil {
		msgData = jsonPacket.Msg
	} else if jsonPacket.Data != nil {
		msgData = jsonPacket.Data
	}

	if msgData != nil {
		var err error

		// 检查msgData是否是字符串（JSON字符串）
		if msgStr, ok := msgData.(string); ok {
			// 如果是字符串，直接使用字符串的字节数组
			dataBytes = []byte(msgStr)
		} else {
			// 如果是对象，序列化为JSON字节数组
			dataBytes, err = json.Marshal(msgData)
			if err != nil {
				return nil, fmt.Errorf("marshal msg failed: %v", err)
			}
		}
	}

	// 创建Packet
	packet := &Packet{
		Header: MessageHeader{
			Cmd:   jsonPacket.Cmd,
			Index: jsonPacket.Index,
			Size:  len(dataBytes),
		},
		Data: dataBytes,
	}

	return []*Packet{packet}, nil
}

// Encode 编码为JSON格式
func (j *JSONProtocolAdapter) Encode(cmd int, index int, data []byte) ([]byte, error) {
	// 尝试将data解析为JSON对象
	var dataObj interface{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &dataObj); err != nil {
			// 如果不是JSON，则作为字符串处理
			dataObj = string(data)
		}
	}

	jsonPacket := WebSocketJsonFrame{
		Cmd:   cmd,
		Index: index,
		Msg:   dataObj,
	}

	return json.Marshal(jsonPacket)
}

// ProtocolFactory 协议工厂
type ProtocolFactory struct{}

// NewProtocolAdapter 根据协议类型创建适配器
func (f *ProtocolFactory) NewProtocolAdapter(protocolType ProtocolType) ProtocolAdapter {
	switch protocolType {
	case ProtocolTypeBinary:
		return NewBinaryProtocolAdapter()
	case ProtocolTypeJSON:
		return NewJSONProtocolAdapter()
	default:
		return NewBinaryProtocolAdapter() // 默认使用二进制协议
	}
}

// DetectProtocolType 检测协议类型
func DetectProtocolType(data []byte) ProtocolType {
	// 检查是否为JSON格式
	if len(data) > 0 && (data[0] == '{' || data[0] == '[') {
		return ProtocolTypeJSON
	}
	return ProtocolTypeBinary
}
