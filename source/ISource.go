package source

import ()

type SingleSource struct { // 目标
	Id      string `json:"id"`
	Content string `json:"content"`
	PicUrl  string `json:"picurl"`
}

type ISource interface {
	Get() SingleSource
}
