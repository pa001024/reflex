package util

import ()

type ArrayList []interface{}

func (this *ArrayList) Add(v interface{}) {
	*this = append(*this, v)
}

func (this *ArrayList) Len() int {
	return len(*this)
}
