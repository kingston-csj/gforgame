package logger_test

import (
	"fmt"
	"io/github/gforgame/logger"
	"testing"
)

func ExampleLog() {
	logger.Info("this is a normal log")
	logger.Error(fmt.Errorf("this is an error log"))
	logger.Log(logger.Player, "Id", "123", "name", "gforgame")
}

func TestLog(t *testing.T) {
	logger.Info("this is a normal log")
	logger.Error(fmt.Errorf("this is an error log"))
	logger.Log(logger.Player, "Id", "123", "name", "gforgame")
}
