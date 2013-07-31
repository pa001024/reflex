package util

// 动态数组封装
type ArrayList []interface{}

// 动态数组封装
func (this *ArrayList) Add(v interface{}) {
	*this = append(*this, v)
}

// 动态数组封装
func (this *ArrayList) Len() int {
	return len(*this)
}
