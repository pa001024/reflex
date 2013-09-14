package util

import (
	"fmt"
	"path"
	"runtime"
)

func LineNumber() string {
	return LineNumberAt(1)
}

func LineNumberAt(deep int) string {
	fup, file, line, _ := runtime.Caller(deep + 1)
	fu := runtime.FuncForPC(fup)
	return fmt.Sprintf("%s(%s:%v)", fu.Name(), path.Base(file), line)
}
