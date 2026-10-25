package session

import (
	"net"
	"strconv"
	"testing"
	"time"

	jsoncodec "github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network/protocol"
)

type factoryMessage struct{ Value string }

func TestSessionUsesOwnFactoryForSendAndDecode(t *testing.T) {
	for _, cmd := range []int32{101, 202} {
		t.Run(strconv.Itoa(int(cmd)), func(t *testing.T) {
			factory := protocol.NewMessageFactory()
			factory.RegisterMessage(cmd, &factoryMessage{})
			a, b := net.Pipe()
			defer a.Close()
			defer b.Close()
			s := NewSession(a, jsoncodec.NewSerializer(), factory)
			defer s.Close()
			if s.GetMessageFactory() != factory {
				t.Fatal("session lost its factory")
			}
			if err := s.Send(&factoryMessage{Value: "outgoing"}, 7); err != nil {
				t.Fatal(err)
			}
			queued := <-s.dataToSend
			packets, err := protocol.NewBinaryProtocolAdapter().Decode(queued.frame)
			if err != nil || len(packets) != 1 {
				t.Fatalf("decode sent frame: %v", err)
			}
			if packets[0].Header.Cmd != cmd || packets[0].Header.Index != 7 {
				t.Fatalf("wrong header: %+v", packets[0].Header)
			}
			s.decodeAndDeliverPackets(packets)
			select {
			case frame := <-s.DataReceived:
				msg, ok := frame.Msg.(*factoryMessage)
				if !ok || msg.Value != "outgoing" {
					t.Fatalf("wrong decoded message: %#v", frame.Msg)
				}
			case <-time.After(time.Second):
				t.Fatal("message was not decoded using the session factory")
			}
		})
	}
}
