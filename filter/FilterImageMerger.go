/*
图片合并模块
=============

简介
-------------

将多张图片合并成一张 使用freetype来渲染字体

*/

package filter

import (
	"code.google.com/p/freetype-go/freetype"
	"github.com/pa001024/MoeCron/source"
)

type FilterImageMerger struct { // 图片合并
	IFilter
	Filter

	MaxPic int    `json:"max_pic"`
	Format string `json:"format"` // png
}

func (this *FilterImageMerger) Process(src []*source.FeedInfo) (dst []*source.FeedInfo) {
	freetype.ParseFont([]byte{})
	return
}