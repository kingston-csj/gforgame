package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"time"
)

type Type int

// 日志类型枚举，每一个类型对应独立的文件
const (
	Admin Type = iota
	Application
	Player
)

var logName = map[Type]string{
	Admin:       "admin",
	Application: "application",
	Player:      "player",
}

var (
	logs       map[string]*logrus.Logger
	consoleLog *logrus.Logger
	errorLog   *logrus.Logger
)

func init() {
	logs = make(map[string]*logrus.Logger)

	for _, logType := range logName {
		logger := createBusinessLog(logType)
		logs[logType] = logger
	}

	// 创建一个新的Logger
	consoleLog = createConsoleLog()

	errorLog = createErrorLog()
}

func createBusinessLog(name string) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &businessLogFormatter{}
	writer, _ := rotatelogs.New(
		"logs/"+name+"/"+name+".%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
	)
	logger.Out = writer
	// 设置Logger的日志级别
	logger.Level = logrus.InfoLevel
	return logger
}

func createConsoleLog() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &consoleLogFormatter{}
	// 设置Logger的输出
	writer, _ := rotatelogs.New(
		"logs/app/"+"app.%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
	)
	logger.Out = writer
	// 设置Logger的日志级别
	logger.Level = logrus.InfoLevel
	return logger
}

func createErrorLog() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	writer, _ := rotatelogs.New(
		"logs/err/"+"error.%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
	)
	logger.Out = writer
	logger.Level = logrus.InfoLevel
	return logger
}

// Log records a message with the specified log type and arguments.
// The Log function writes a log message to the log file
// Parameters:
//
//	name LogType: The type of log message being recorded.
//	args ...interface{}: A variadic list of arguments to include in the log message.
//
// Example:
//
//	Log(Admin, "User logged in", user)
//	Log(Player, "Failed to save data", err)
func Log(name Type, args ...interface{}) {
	if len(args)%2 != 0 {
		panic("log arguments must be odd number")
	}
	logger := logs[logName[name]]

	sb := &strings.Builder{}
	sb.WriteString("time|")
	sb.WriteString(fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	sb.WriteString("|")
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %v", args[i]))
		}
		value := args[i+1]
		sb.WriteString(key)
		sb.WriteString("|")
		sb.WriteString(fmt.Sprintf("%v", value))
		sb.WriteString("|")
	}
	sb.WriteString("\n")
	logger.Info(sb.String())
}

// Info 记录一条日志
func Info(v string) {
	consoleLog.Info(v)
}

func Error(err error) {
	if err != nil {
		stack := make([]byte, 1024)
		n := runtime.Stack(stack, false)
		errorLog.WithFields(logrus.Fields{
			"error": err,
		}).Error(string(stack[:n]))
	}
}
