package util

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

// 默认四个日志级别
var (
	TRACE    = NewLogger(os.Stderr, false, "[TRACE] ")
	DEBUG    = NewLogger(os.Stderr, false, "[DEBUG] ")
	WARN     = NewLogger(os.Stderr, true, "[WARN] ")
	ERROR    = NewLogger(os.Stderr, true, "[ERROR] ")
	INFO     = NewLogger(os.Stdout, true, "[INFO] ")
	CRITICAL = NewLogger(os.Stdout, true, "[CRITICAL] ")
)

// 调试时日志 输出行号[调用深度还真头疼啊]
func DebugLog(args ...interface{}) {
	debugFlog(os.Stdout, args...)
}

// 调试时日志 输出到流
func DebugFlog(w io.Writer, args ...interface{}) {
	debugFlog(w, args...)
}

func debugFlog(w io.Writer, args ...interface{}) {
	fup, file, line, _ := runtime.Caller(2)
	fu := runtime.FuncForPC(fup)
	fmt.Fprintln(w, args...)
	fmt.Fprintln(w, "        at", fu.Name()+"("+path.Base(file)+":"+ToString(line)+")")
}

// 日志对象
type Logger struct {
	output io.Writer
	enable bool
	perfix string
}

// 创建新日志对象
func NewLogger(w io.Writer, enable bool, perfix string) *Logger {
	return &Logger{w, enable, perfix}
}

// 输出日志
func (l *Logger) Log(s ...interface{}) {
	if l.enable {
		fmt.Fprintf(l.output, "%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(s...))
	}
}

// 输出日志和调试信息
func (l *Logger) Err(s ...interface{}) {
	if l.enable {
		debugFlog(l.output, fmt.Sprintf("%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(s...)))
	}
}
func (l *Logger) Logf(format string, s ...interface{}) {
	if l.enable {
		fmt.Fprintf(l.output, "%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, s...))
	}
}
func (l *Logger) Errf(format string, s ...interface{}) {
	if l.enable {
		debugFlog(l.output, fmt.Sprintf("%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, s...)))
	}
}

// 返回启用状态
func (l *Logger) Enable() bool {
	return l.enable
}

// 设置启用状态
func (l *Logger) SetEnable(v bool) {
	l.enable = v
}

// 返回输出
func (l *Logger) Output() io.Writer {
	return l.output
}

// 设置输出
func (l *Logger) SetOutput(v io.Writer) {
	l.output = v
}
