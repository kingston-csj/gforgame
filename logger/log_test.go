package logger

import "fmt"

func ExampleLog() {
	Info("this is a normal log")
	Error(fmt.Errorf("this is an error log"))
	Log(Player, "Id", "123", "name", "gforgame")
}
