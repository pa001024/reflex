package util

import (
	"sync"
)

type Key interface{}
type Value interface{}

// Thread-safe map 线程安全的Map
type SafeMap struct {
	sync.RWMutex
	m       map[Key]Value
	reverse map[Value][]Key
}

func NewSafeMap() (m *SafeMap) {
	return &SafeMap{sync.RWMutex{}, make(map[Key]Value), nil}
}

// 全部返回 供range使用 在多线程下使用请手动进行 RLock(); defer RUnlock()
func (this *SafeMap) DirtyItems() map[Key]Value { return this.m }

// 脏写 请勿在多线程下使用
func (this *SafeMap) DirtyPut(key Key, v Value) { this.m[key] = v }

// 脏读 请勿在多线程下使用
func (this *SafeMap) DirtyGet(key Key) Value { return this.m[key] }

// 脏读 请勿在多线程下使用
func (this *SafeMap) DirtyDelete(key Key) { delete(this.m, key) }

// 返回全部数据的副本[很慢] 一般情况下请使用 DirtyItems()
func (this *SafeMap) Items() (nv map[Key]Value) {
	nv = make(map[Key]Value)
	this.RLock()
	for k, v := range this.m {
		nv[k] = v
	}
	this.RUnlock()
	return
}

// 放入数据
func (this *SafeMap) Put(key Key, v Value) {
	this.Lock()
	this.m[key] = v
	this.Unlock()
}

// 取出
func (this *SafeMap) Get(key Key) (v Value) {
	this.RLock()
	v, ok := this.m[v]
	this.RUnlock()
	if !ok {
		return nil
	}
	return
}

// 删除
func (this *SafeMap) Delete(key Key) {
	this.Lock()
	delete(this.m, key)
	this.Unlock()
}

// 检查是否存在某键
func (this *SafeMap) DirtyContainsKey(key Key) (v bool) {
	_, v = this.m[key]
	return
}

// 检查是否存在某几个键
func (this *SafeMap) DirtyContainsKeys(keys []Key) (v bool) {
	for _, key := range keys {
		_, v = this.m[key]
		if !v {
			return
		}
	}
	return
}

// 检查是否存在某键
func (this *SafeMap) ContainsKey(key Key) (v bool) {
	this.RLock()
	_, v = this.m[key]
	this.RUnlock()
	return
}

// 检查是否存在某几个键
func (this *SafeMap) ContainsKeys(keys []Key) (v bool) {
	this.RLock()
	defer this.RUnlock() // 必须
	for _, key := range keys {
		_, v = this.m[key]
		if !v {
			return
		}
	}
	return
}

// [脏的]对Value建立索引 原理是建立一个反向 map[Value]Key 建立索引后ContainValue()效率将变为O(1)
func (this *SafeMap) DirtyMakeReverseIndex() {
	this.reverse = make(map[Value][]Key)
	for k, v := range this.m {
		if _, ok := this.reverse[v]; !ok {
			this.reverse[v] = make([]Key, 0, 1)
		}
		this.reverse[v] = append(this.reverse[v], k)
	}
}

// [线程安全的]对Value建立索引 原理是建立一个反向 map[Value]Key 建立索引后ContainValue()效率将变为O(1) 但索引不会自动变化(因为影响效率) 需要手动重建
func (this *SafeMap) MakeReverseIndex() {
	this.reverse = make(map[Value][]Key)
	this.RLock()
	for k, v := range this.m {
		if _, ok := this.reverse[v]; !ok {
			this.reverse[v] = make([]Key, 0, 1)
		}
		this.reverse[v] = append(this.reverse[v], k)
	}
	this.RUnlock()
}

// 检测是否建立Value反向索引
func (this *SafeMap) HaveReverseIndex() bool { return this.reverse != nil }

// 无索引下是O(n)遍历 慎用
func (this *SafeMap) ContainsValue(val Value) (v bool) {
	if this.reverse != nil {
		_, v = this.reverse[val]
	} else {
		for _, v1 := range this.m {
			if v1 == val {
				return true
			}
		}
	}
	return
}

// 取值返回int类型
func (this *SafeMap) Int(key Key) (v int) {
	v, _ = this.Get(key).(int)
	return
}

// 取值返回int8类型
func (this *SafeMap) Int8(key Key) (v int8) {
	v, _ = this.Get(key).(int8)
	return
}

// 取值返回int16类型
func (this *SafeMap) Int16(key Key) (v int16) {
	v, _ = this.Get(key).(int16)
	return
}

// 取值返回int32类型
func (this *SafeMap) Int32(key Key) (v int32) {
	v, _ = this.Get(key).(int32)
	return
}

// 取值返回int64类型
func (this *SafeMap) Int64(key Key) (v int64) {
	v, _ = this.Get(key).(int64)
	return
}

// 取值返回uint类型
func (this *SafeMap) UInt(key Key) (v uint) {
	v, _ = this.Get(key).(uint)
	return
}

// 取值返回int8类型
func (this *SafeMap) UInt8(key Key) (v uint8) {
	v, _ = this.Get(key).(uint8)
	return
}

// 取值返回int16类型
func (this *SafeMap) UInt16(key Key) (v uint16) {
	v, _ = this.Get(key).(uint16)
	return
}

// 取值返回uint32类型
func (this *SafeMap) UInt32(key Key) (v uint32) {
	v, _ = this.Get(key).(uint32)
	return
}

// 取值返回uint64类型
func (this *SafeMap) UInt64(key Key) (v uint64) {
	v, _ = this.Get(key).(uint64)
	return
}

// 取值返回float32类型
func (this *SafeMap) Float32(key Key) (v float32) {
	v, _ = this.Get(key).(float32)
	return
}

// 取值返回float64类型
func (this *SafeMap) Float64(key Key) (v float64) {
	v, _ = this.Get(key).(float64)
	return
}

// 取值返回byte类型
func (this *SafeMap) Byte(key Key) (v byte) {
	v, _ = this.Get(key).(byte)
	return
}

// 取值返回byte类型
func (this *SafeMap) Bool(key Key) (v bool) {
	v, _ = this.Get(key).(bool)
	return
}

// 取值返回string类型
func (this *SafeMap) String(key Key) (v string) {
	v, _ = this.Get(key).(string)
	return
}
