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
		panic(err)
	}
	return y
}

// 无错返回ReadAll
func MustReadAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	Throw(err.Error())
	return b
}

// 出错崩溃
func Try(e error) {
	if e != nil {
		Throw(e.Error())
	}
}

func Throw(s string) {
	fup, file, line, _ := runtime.Caller(2)
	fu := runtime.FuncForPC(fup)
	panic(fmt.Sprintf("\x0b\xad\xca\xfe%s\n    at %s(%s:%v)", s, fu.Name(), path.Base(file), line))
}

// 供给panic后恢复
func Catch() {
	if e := recover(); e != nil {
		es := fmt.Sprint(e)
		if es[:4] == "\x0b\xad\xca\xfe" {
			ERROR.Log(es[4:])
		} else {
			panic(e)
		}
	}
}
