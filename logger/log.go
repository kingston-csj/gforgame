package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Type int

// 日志类型枚举，每一个类型对应独立的文件
const (
	Admin Type = iota
	Application
	Player
	Mail
	Item
	Activity
)

var logName = map[Type]string{
	Admin:       "admin",
	Application: "application",
	Player:      "player",
	Mail:        "mail",
	Item:        "item",
	Activity:    "activity",
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
	// 创建多输出目标：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)
	logger.Out = multiWriter
	// 设置Logger的日志级别
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

func Debugf(format string, args ...interface{}) {
	consoleLog.Debugf(format, args...)
}

func createErrorLog() *logrus.Logger {
	logger := logrus.New()

	// 确保日志目录存在
	logDir := "logs/err"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatalf("Failed to create log directory: %v", err)
	}

	// 创建 rotatelogs 实例
	writer, err := rotatelogs.New(
		filepath.Join(logDir, "error.%Y%m%d"),
		rotatelogs.WithMaxAge(24*time.Hour),
	)
	if err != nil {
		logrus.Fatalf("Failed to create rotatelogs: %v", err)
	}

	// 创建多输出目标：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)

	// 自定义日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true, // 禁用默认时间戳
		DisableQuote:     true, // 禁用字段值的引号
	})

	// 设置日志输出
	logger.Out = multiWriter

	// 设置日志级别
	logger.Level = logrus.ErrorLevel

	return logger
}

func Error(err error) {
	if err != nil {
		// 获取调用栈信息
		stack := make([]byte, 1024)
		n := runtime.Stack(stack, false)
		stackStr := string(stack[:n])

		// 优化调用栈输出
		stackLines := strings.Split(stackStr, "\n")
		var optimizedStack []string
		for i, line := range stackLines {
			if i%2 == 0 {
				// 函数名行
				optimizedStack = append(optimizedStack, line)
			} else {
				// 文件路径行，只显示文件名
				parts := strings.Split(line, "/")
				if len(parts) > 0 {
					optimizedStack = append(optimizedStack, "\t"+parts[len(parts)-1])
				}
			}
		}

		// 记录日志
		errorLog.Out.Write([]byte(fmt.Sprintf(
			"%s ERROR EXCEPTION - %s \n%s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			fmt.Sprintf("\"%s\"", err),
			strings.Join(optimizedStack, "\n"),
		)))
	}
}
