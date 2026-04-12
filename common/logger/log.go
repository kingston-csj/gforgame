package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type plainTextHandler struct {
	out   io.Writer
	level slog.Level
	mu    sync.Mutex
}

func (h *plainTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *plainTextHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}
	level := strings.ToLower(r.Level.String())
	line := fmt.Sprintf("%s [%s] %s\n", ts.Format("2006-01-02 15:04:05"), level, r.Message)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.out, line)
	return err
}

func (h *plainTextHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *plainTextHandler) WithGroup(_ string) slog.Handler {
	return h
}

var (
	logs       map[string]*slog.Logger
	logsMu     sync.RWMutex
	consoleLog *slog.Logger
	errorLog   *slog.Logger
)

func init() {
	logs = make(map[string]*slog.Logger)
	// 创建一个新的Logger
	consoleLog = createConsoleLog()

	errorLog = createErrorLog()
}

func createBusinessLog(name string) *slog.Logger {
	name = strings.ToLower(name)
	writer, _ := rotatelogs.New(
		"logs/"+name+"/"+name+".%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
	)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler).With("biz", name)
}

func createConsoleLog() *slog.Logger {
	// 设置Logger的输出
	writer, _ := rotatelogs.New(
		"logs/app/"+"app.%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
	)
	// 创建多输出目标：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)
	handler := &plainTextHandler{out: multiWriter, level: slog.LevelDebug}
	return slog.New(handler)
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
func Log(name string, args ...interface{}) {
	if len(args)%2 != 0 {
		panic("log arguments must be odd number")
	}
	key := strings.ToLower(name)
	logsMu.RLock()
	logger := logs[key]
	logsMu.RUnlock()
	if logger == nil {
		logsMu.Lock()
		// double-check，避免并发重复创建
		logger = logs[key]
		if logger == nil {
			logger = createBusinessLog(key)
			logs[key] = logger
		}
		logsMu.Unlock()
	}

	fields := make([]any, 0, len(args)+2)
	fields = append(fields, "time", time.Now().UnixNano()/1000000)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %v", args[i]))
		}
		value := args[i+1]
		fields = append(fields, key, value)
	}
	logger.Info("biz_log", fields...)
}

// Info 记录一条日志
func Info(v string) {
	consoleLog.Info(v)
}

func Debugf(format string, args ...interface{}) {
	consoleLog.Debug(fmt.Sprintf(format, args...))
}

func createErrorLog() *slog.Logger {
	// 确保日志目录存在
	logDir := "logs/err"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	// 创建 rotatelogs 实例
	writer, err := rotatelogs.New(
		filepath.Join(logDir, "error.%Y%m%d"),
		rotatelogs.WithMaxAge(24*time.Hour),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create rotatelogs: %v", err))
	}

	// 创建多输出目标：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler).With("biz", "error")
}

// Error 记录带自定义异常内容、原始错误和调用栈的错误日志
// 参数：
//   - customMsg: 自定义异常描述（如"活动调度任务取消失败"），可为空
//   - err: 原始错误对象，若为nil则不记录日志
func Error(customMsg string, err error) {
	if err != nil {
		// 获取调用栈信息
		stack := make([]byte, 1024)
		n := runtime.Stack(stack, false)
		stackStr := string(stack[:n])

		// 优化调用栈输出（保留原有逻辑）
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

		// 拼接自定义异常内容（为空则不显示，避免多余字符）
		var msgPrefix string
		if customMsg != "" {
			msgPrefix = fmt.Sprintf("%s - ", customMsg)
		}

		// 记录日志（整合自定义内容、原始错误、调用栈）
		errorLog.Error("ERROR EXCEPTION",
			"custom", msgPrefix,
			"error", err.Error(),
			"stack", strings.Join(optimizedStack, "\n"),
		)
	}
}

func ErrorNoStack(customMsg any) {
	if customMsg == nil {
		return
	}
	var message string
	switch v := customMsg.(type) {
	case string:
		message = strings.TrimSpace(v)
	case error:
		message = strings.TrimSpace(v.Error())
	default:
		message = strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	if message == "" {
		return
	}
	// 重要错误单行日志：写入 error 文件，不带调用栈
	errorLog.Error("ERROR", "message", message)
}
