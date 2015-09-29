package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
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
func Try(err error) {
	if err != nil {
		Throw(err.Error())
	}
}

// 输出调用者
func lastCaller(msg string, deep int) string {
	return fmt.Sprintf("%s\n    at %s", msg, LineNumberAt(deep+1))
}

// 内部抛出错误
func throw(msg string) {
	panic(lastCaller(msg, 1))
}

// 抛出错误
func Throw(msg string) {
	panic(lastCaller(msg, 2))
}

// 供给panic后恢复 (部分内部错误无法恢复)
// exmaple:
// func xx() (err error){
//   defer util.Catch(&err)
//   ...
//   util.Try(err)
// }
//
func Catch(err ...*error) {
	if e := recover(); e != nil {
		es := fmt.Sprint(e)
		if err != nil && len(err) > 0 {
			*err[0], _ = e.(error)
		}
		ERROR.Log(es)
	}
}
