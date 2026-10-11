package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
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

// 滚动周期时间格式
const (
	RollingDaily  = "%Y%m%d"   // 按天滚动
	RollingHourly = "%Y%m%d%H" // 按小时滚动
)

const DefaultBusinessLogName = "business"

var (
	logs       map[string]*slog.Logger
	logsMu     sync.RWMutex
	consoleLog *slog.Logger
	errorLog   *slog.Logger
)

type stackTracer interface {
	StackTrace() string
}

type panicWithStack struct {
	value any
	stack []byte
}

func (p *panicWithStack) Error() string {
	return fmt.Sprintf("%v", p.value)
}

func (p *panicWithStack) StackTrace() string {
	return string(p.stack)
}

// PanicToError converts a recovered panic value into an error with the original panic stack.
func PanicToError(p any) error {
	return &panicWithStack{
		value: p,
		stack: debug.Stack(),
	}
}

func init() {
	logs = make(map[string]*slog.Logger)
	// 创建一个新的Logger
	consoleLog = createConsoleLog()

	errorLog = createErrorLog()
}

func createBusinessLog(dirName string, name string, rollingPattern string) *slog.Logger {
	name = strings.ToLower(name)
	dirName = strings.ToLower(dirName)
	writer, _ := rotatelogs.New(
		"logs/"+dirName+"/"+name+"."+rollingPattern,
		rotatelogs.WithMaxAge(30*time.Duration(24)*time.Hour),
	)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 业务日志字段映射：
			// - 默认 time -> date(格式化字符串)
			// - 自定义 time(毫秒时间戳) 在 Log() 中写入
			switch a.Key {
			case slog.TimeKey:
				// 只处理 slog 内置的 time（Value 为 time.Time 类型）
				// 改键为 date，存格式化字符串
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String("date", t.Format("2006-01-02 15:04:05.000"))
				}
				// 用户自定义的同名 attr（非 time.Time 类型）保持原键名不动
				return a
			case slog.LevelKey, slog.MessageKey:
				return slog.Attr{}
			default:
				return a
			}
		},
	})
	return slog.New(handler)
}

func createConsoleLog() *slog.Logger {
	// 设置Logger的输出
	writer, _ := rotatelogs.New(
		"logs/app/"+"app.%Y%m%d",
		rotatelogs.WithMaxAge(7*time.Duration(24)*time.Hour),
	)
	// 创建多输出目标：同时输出到文件和控制台
	multiWriter := io.MultiWriter(writer, os.Stdout)
	handler := &plainTextHandler{out: multiWriter, level: slog.LevelDebug}
	return slog.New(handler)
}

// logWithPeriod 内部通用写日志函数
// dirName:  实际写入的二级目录名。系统日志=name（每个日志独立目录），运营日志="business"（统一目录）
func logWithPeriod(dirName string, name string, rollingPattern string, args ...any) {
	if len(args)%2 != 0 {
		panic("log arguments must be even number (key-value pairs)")
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
			logger = createBusinessLog(dirName, key, rollingPattern)
			logs[key] = logger
		}
		logsMu.Unlock()
	}

	fields := make([]any, 0, len(args)+2)
	// timestamp: 毫秒级 unix 时间戳（数字型），与 slog 内置 time 区分开
	fields = append(fields, "time", time.Now().UnixNano()/1000000)
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %v", args[i]))
		}
		value := args[i+1]
		fields = append(fields, k, value)
	}
	logger.Info("biz_log", fields...)
}

// LogSystem 系统日志（按天滚动，每个日志类型独立二级目录）
// 路径：logs/{name}/{name}.%Y%m%d
// 例：logs/admin/admin.20260920  logs/debug/debug.20260920
//
// 【禁止业务代码直接调用】请使用 internal/logger.LogSystem(category, args...)
// 由封装层统一注入 sid/module/logType/logDesc 等标准字段。
func LogSystem(loggerName string, args ...any) {
	logWithPeriod(loggerName, loggerName, RollingDaily, args...)
}

// LogBusiness 运营/业务日志（按小时滚动，统一文件 business.%Y%m%d%H）
// 路径：logs/business/business.%Y%m%d%H
//
// 【禁止业务代码直接调用】请使用 internal/logger.LogBusinessDefault(args...)
// 由封装层统一注入 sid 等标准字段。
func LogBusiness(args ...any) {
	logWithPeriod("business", DefaultBusinessLogName, RollingHourly, args...)
}

// LogBusinessTo 指定业务模块名 + 按小时滚动，统一写入 business.%Y%m%d%H
// 路径：logs/business/business.%Y%m%d%H
// 注："module" 字段已由业务层（internal/logger/logutil.go）写入，此处不再重复注入，避免 JSON 中出现重复 key。
//
// 【禁止业务代码直接调用】请使用 internal/logger.LogBusinessWithActor 或 LogBusinessNoActor
// 由封装层统一注入 module/logType/logDesc 等标准字段。
func LogBusinessTo(loggerName string, args ...any) {
	logWithPeriod("business", DefaultBusinessLogName, RollingHourly, args...)
}

// Info 记录一条日志
func Info(v string) {
	consoleLog.Info(v)
}

func Debugf(format string, args ...any) {
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
		rotatelogs.WithMaxAge(7*time.Duration(24)*time.Hour),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create rotatelogs: %v", err))
	}

	// 错误日志仅写入 err 文件；控制台和 app 日志由 consoleLog 统一输出，避免重复打印
	handler := &plainTextHandler{out: writer, level: slog.LevelInfo}
	return slog.New(handler)
}

// Error 记录带自定义异常内容、原始错误和调用栈的错误日志
// 参数：
//   - customMsg: 自定义异常描述（如"活动调度任务取消失败"），可为空
//   - err: 原始错误对象，若为nil则不记录日志
func Error(customMsg string, err error) {
	if err != nil {
		stackStr := ""
		if st, ok := err.(stackTracer); ok {
			stackStr = st.StackTrace()
		} else {
			// 回退：记录当前调用栈，增大容量减少截断概率
			stack := make([]byte, 64*1024)
			n := runtime.Stack(stack, false)
			stackStr = string(stack[:n])
		}

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

		// 拼接多行日志，堆栈按真实换行输出，便于排查
		message := "ERROR EXCEPTION"
		if customMsg != "" {
			message += "\ncustom: " + customMsg
		}
		message += "\nerror: " + err.Error()
		message += "\nstack:\n" + strings.Join(optimizedStack, "\n")
		errorLog.Error(message)
		consoleLog.Error(message)
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
	message = "ERROR: " + message
	errorLog.Error(message)
	consoleLog.Error(message)
}
