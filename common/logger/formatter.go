package logger

import (
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type consoleLogFormatter struct{}

// Format 实现 logrus.Formatter 接口的 Format 方法。
func (f *consoleLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 按照自定义格式写入日志信息
	_, _ = fmt.Fprintf(b, "%s [%s] %s\n", entry.Time.Format(time.DateTime), entry.Level, entry.Message)
	return b.Bytes(), nil
}

type businessLogFormatter struct{}

// Format 实现 logrus.Formatter 接口的 Format 方法。
func (f *businessLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	_, _ = fmt.Fprintf(b, "%s", entry.Message)
	return b.Bytes(), nil
}
