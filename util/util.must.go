package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"runtime"
)

// 无错返回ParseUrl
func MustParseUrl(rawurl string) *url.URL {
	y, err := url.Parse(rawurl)
	if err != nil {
		Throw(err.Error())
	}
	return y
}

// 无错返回ReadAll
func MustReadAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		Throw(err.Error())
	}
	return b
}

// 出错崩溃
func Try(e error) {
	if e != nil {
		Throw(e.Error())
	}
}

// 输出调用者
func lastCaller(msg string, deep int) string {
	fup, file, line, _ := runtime.Caller(deep + 1)
	fu := runtime.FuncForPC(fup)
	return fmt.Sprintf("%s\n    at %s(%s:%v)", msg, fu.Name(), path.Base(file), line)
}

func LastCaller(msg string) string {
	return lastCaller(msg, 1)
}

// 抛出错误
func throw(msg string) {
	panic(lastCaller(msg, 1))
}

// 抛出错误
func Throw(msg string) {
	panic(lastCaller(msg, 2))
}

// 供给panic后恢复
func Catch(err ...*error) {
	if e := recover(); e != nil {
		es := fmt.Sprint(e)
		*err[0], _ = e.(error)
		ERROR.Log(es)
	}
}
