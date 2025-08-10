package ws

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	conn   *websocket.Conn
	typ    int // message type
	reader io.Reader
	// 协议类型：二进制或文本
	protocolType int
}

func newWSConn(conn *websocket.Conn) (*wsConn, error) {
	c := &wsConn{conn: conn}

	// 不在这里读取消息，让Session来处理
	// 只设置默认的协议类型
	c.protocolType = websocket.BinaryMessage // 默认二进制协议

	return c, nil
}

func (c *wsConn) Read(b []byte) (int, error) {
	// WebSocket连接不应该使用这个方法，应该使用ReadMessage
	// 这里提供一个兼容的实现，但会返回错误
	return 0, fmt.Errorf("WebSocket connection should use ReadMessage() instead of Read()")
}

func (c *wsConn) Write(b []byte) (int, error) {
	// 根据协议类型选择消息类型
	messageType := websocket.BinaryMessage
	if c.protocolType == websocket.TextMessage {
		messageType = websocket.TextMessage
	}

	err := c.conn.WriteMessage(messageType, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// ReadMessage 实现WebSocketConn接口
func (c *wsConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

// WriteMessage 实现WebSocketConn接口
func (c *wsConn) WriteMessage(messageType int, data []byte) error {
	return c.conn.WriteMessage(messageType, data)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *wsConn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *wsConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *wsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.conn.SetReadDeadline(t); err != nil {
		return err
	}

	return c.conn.SetWriteDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
