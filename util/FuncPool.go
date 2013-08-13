package util

import (
	"fmt"
	"reflect"
)

// 动态模式函数调用
type FuncPool map[string]interface{}

// 调用
func (m FuncPool) Call(name string, args ...interface{}) (out []interface{}, err error) {
	if m[name] == nil {
		return nil, fmt.Errorf("Func \"%s\" not exists ", name)
	}
	v := reflect.ValueOf(m[name])
	vs := v.Type().NumIn()
	if len(args) < vs {
		return nil, fmt.Errorf("[%s] Not enough params: need %v, given %v", name, vs, len(args))
	}
	in := make([]reflect.Value, vs)
	for i, _ := range in {
		in[i] = reflect.ValueOf(args[i])
	}
	out = make([]interface{}, v.Type().NumOut())
	ret := v.Call(in)
	for i, _ := range out {
		out[i] = ret[i].Interface()
	}
	return
}

// 出错则panic
func (m FuncPool) MustCall(name string, args ...interface{}) (out []interface{}) {
	o, err := m.Call(name, args...)
	if err != nil {
		panic(err)
	}
	return o
}
