package logger

import (
	"fmt"
	"testing"
)

func ExampleLog() {
	Info("this is a normal log")
	Error(fmt.Errorf("this is an error log"))
	Log(Player, "Id", "123", "name", "gforgame")
}

func TestLog(t *testing.T) {
	Info("this is a normal log")
	Error(fmt.Errorf("this is an error log"))
	Log(Player, "Id", "123", "name", "gforgame")
}
