package main

import (
	"encoding/json"
	"fmt"
	"io/github/gforgame/network/protocol"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 测试用的消息结构
type TestMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func main() {
	// 启动WebSocket服务器
	go startWebSocketServer()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 测试JSON协议
	testJSONProtocol()

	// 测试二进制协议
	testBinaryProtocol()
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("WebSocket客户端已连接")

	for {
		// 读取消息类型和数据
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}

		fmt.Printf("收到消息类型: %d, 数据长度: %d\n", messageType, len(messageData))

		// 根据消息类型处理
		if messageType == websocket.TextMessage {
			// JSON协议
			handleJSONMessage(conn, messageData)
		} else {
			// 二进制协议
			handleBinaryMessage(conn, messageData)
		}
	}
}

func handleJSONMessage(conn *websocket.Conn, data []byte) {
	fmt.Printf("处理JSON消息: %s\n", string(data))

	// 解析JSON消息
	var jsonPacket protocol.WebSocketJsonFrame
	if err := json.Unmarshal(data, &jsonPacket); err != nil {
		fmt.Printf("JSON解析失败: %v\n", err)
		return
	}

	fmt.Printf("JSON消息 - Cmd: %d, Index: %d, Data: %v\n", jsonPacket.Cmd, jsonPacket.Index, jsonPacket.Msg)

	// 发送响应
	response := protocol.WebSocketJsonFrame{
		Cmd:   jsonPacket.Cmd,
		Index: jsonPacket.Index,
		Msg: string(time.Now().Unix()),
	}

	responseData, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseData)
}

func handleBinaryMessage(conn *websocket.Conn, data []byte) {
	fmt.Printf("处理二进制消息，长度: %d\n", len(data))

	// 使用二进制协议解码
	protocol := protocol.NewDecoder()
	packets, err := protocol.Decode(data)
	if err != nil {
		fmt.Printf("二进制协议解码失败: %v\n", err)
		return
	}

	for _, packet := range packets {
		fmt.Printf("二进制消息 - Cmd: %d, Index: %d, Size: %d\n",
			packet.Header.Cmd, packet.Header.Index, packet.Header.Size)
	}

	// 发送响应
	responseData, _ := protocol.Encode(1001, 1, []byte("binary response"))
	conn.WriteMessage(websocket.BinaryMessage, responseData)
}

func testJSONProtocol() {
	fmt.Println("\n=== 测试JSON协议 ===")

	// 连接WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Printf("连接失败: %v", err)
		return
	}
	defer conn.Close()

	// 发送JSON消息
	jsonPacket := protocol.WebSocketJsonFrame{
		Cmd:   1001,
		Index: 1,
		Msg: string(time.Now().Unix()),
	}

	jsonData, _ := json.Marshal(jsonPacket)
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Printf("发送JSON消息失败: %v", err)
		return
	}

	// 读取响应
	_, responseData, err := conn.ReadMessage()
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return
	}

	fmt.Printf("收到JSON响应: %s\n", string(responseData))
}

func testBinaryProtocol() {
	fmt.Println("\n=== 测试二进制协议 ===")

	// 连接WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Printf("连接失败: %v", err)
		return
	}
	defer conn.Close()

	// 发送二进制消息
	protocol := protocol.NewDecoder()
	testData := []byte("Hello Binary!")
	binaryData, _ := protocol.Encode(1002, 2, testData)

	err = conn.WriteMessage(websocket.BinaryMessage, binaryData)
	if err != nil {
		log.Printf("发送二进制消息失败: %v", err)
		return
	}

	// 读取响应
	_, responseData, err := conn.ReadMessage()
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return
	}

	fmt.Printf("收到二进制响应，长度: %d\n", len(responseData))
}
