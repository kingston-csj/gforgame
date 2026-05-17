package client

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	// 服务器地址，支持 ws://host:port/path 或 host:port
	RemoteAddress string
	// websocket 路径，未配置时默认使用 /ws
	WsPath string
	// 消息编解码
	MsgCodec codec.MessageCodec
	// 消息分发器
	IoDispatcher network.IoDispatch
	// 可选请求头
	Header http.Header
	// 可选自定义拨号器
	Dialer *websocket.Dialer
	// 消息类型，默认二进制
	MessageType int
}

func (c *WebSocketClient) OpenSession() (*network.Session, error) {
	targetURL, err := c.buildURL()
	if err != nil {
		return nil, err
	}

	dialer := c.Dialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}

	conn, _, err := dialer.Dial(targetURL, c.Header)
	if err != nil {
		return nil, err
	}

	wsConn := &websocketClientConn{
		conn:        conn,
		messageType: c.resolveMessageType(),
	}
	session := network.NewSession(wsConn, c.MsgCodec)
	go session.Write()
	go session.Read()
	return session, nil
}

func (c *WebSocketClient) buildURL() (string, error) {
	address := strings.TrimSpace(c.RemoteAddress)
	if address == "" {
		return "", fmt.Errorf("remote address is empty")
	}

	if strings.HasPrefix(address, "ws://") || strings.HasPrefix(address, "wss://") {
		return c.normalizeURL(address)
	}

	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		u, err := url.Parse(address)
		if err != nil {
			return "", err
		}
		if u.Scheme == "http" {
			u.Scheme = "ws"
		} else {
			u.Scheme = "wss"
		}
		return c.normalizeParsedURL(u)
	}

	defaultScheme := "ws"
	u := &url.URL{
		Scheme: defaultScheme,
		Host:   address,
	}
	return c.normalizeParsedURL(u)
}

func (c *WebSocketClient) normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return c.normalizeParsedURL(u)
}

func (c *WebSocketClient) normalizeParsedURL(u *url.URL) (string, error) {
	if strings.TrimSpace(u.Host) == "" {
		return "", fmt.Errorf("websocket host is empty")
	}

	wsPath := strings.TrimSpace(c.WsPath)
	if wsPath == "" {
		wsPath = "/ws"
	}
	if !strings.HasPrefix(wsPath, "/") {
		wsPath = "/" + wsPath
	}

	if strings.TrimSpace(u.Path) == "" || u.Path == "/" {
		u.Path = wsPath
	} else {
		u.Path = path.Clean(u.Path)
		if !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}
	}

	return u.String(), nil
}

func (c *WebSocketClient) resolveMessageType() int {
	if c.MessageType == websocket.TextMessage {
		return websocket.TextMessage
	}
	return websocket.BinaryMessage
}

type websocketClientConn struct {
	conn        *websocket.Conn
	messageType int
}

func (c *websocketClientConn) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("websocket connection should use ReadMessage() instead of Read()")
}

func (c *websocketClientConn) Write(b []byte) (int, error) {
	if err := c.conn.WriteMessage(c.messageType, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *websocketClientConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *websocketClientConn) WriteMessage(messageType int, data []byte) error {
	return c.conn.WriteMessage(messageType, data)
}

func (c *websocketClientConn) Close() error {
	return c.conn.Close()
}

func (c *websocketClientConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *websocketClientConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *websocketClientConn) SetDeadline(t time.Time) error {
	if err := c.conn.SetReadDeadline(t); err != nil {
		return err
	}
	return c.conn.SetWriteDeadline(t)
}

func (c *websocketClientConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *websocketClientConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
