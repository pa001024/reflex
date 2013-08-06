package util

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"time"
)

var (
	DEBUG = false
)

// 日志
func Log(args ...interface{}) {
	if DEBUG {
		debugLog(args...)
	} else {
		procLog(args...)
	}
}

// 运行时日志
func procLog(args ...interface{}) {
	log.Println(args...)
}

// 调试时日志
func debugLog(args ...interface{}) {
	fup, file, line, _ := runtime.Caller(2)
	fu := runtime.FuncForPC(fup)
	log.Println(args...)
	fmt.Println("        at", fu.Name()+"("+path.Base(file)+":"+ToString(line)+")")
}

// 日志对象
type Logger struct {
	enable bool
	perfix string
}

func NewLogger(enable bool, perfix string) *Logger {
	return &Logger{enable, perfix}
}

func (l *Logger) Log(s ...interface{}) {
	if l.enable {
		fmt.Printf("%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(s...))
	}
}
func (l *Logger) Err(s ...interface{}) {
	if l.enable {
		debugLog(fmt.Sprintf("%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(s...)))
	}
}
func (l *Logger) Logf(format string, s ...interface{}) {
	if l.enable {
		l.Log(fmt.Sprintf(format, s...))
	}
}
func (l *Logger) Errf(format string, s ...interface{}) {
	if l.enable {
		l.Err(fmt.Sprintf(format, s...))
	}
}
func (l *Logger) Enable() {
	l.enable = true
}
func (l *Logger) Disable() {
	l.enable = false
}
