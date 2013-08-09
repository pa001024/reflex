package util

// 动态数组封装
type ArrayList []interface{}

// 追加元素
func (this ArrayList) Append(v interface{}) {
	this = append(this, v)
}

// 追加数组
func (this ArrayList) AppendArray(v []interface{}) {
	this = append(this, v...)
}

// 返回长度
func (this ArrayList) Len() int {
	return len(this)
}

// 返回数组元素位置
func (this ArrayList) IndexOf(val interface{}) int {
	for i, v := range this {
		if v == val {
			return i
		}
	}
	return -1
}
