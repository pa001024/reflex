/*
简介
=======
Mediawiki 的 ISource 接口实现

功能
--------
1. 获取文字内容(使用wiki API提取源代码)
2. 将图片和文字转交给filter进行文字转换裁剪/添加水印操作后提交给Target进行最终的上传/发送/转发等

*/
package source

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

type SourceMediawiki struct { // Mediawiki 实现接口ISource
	ISource

	FeedUrl string // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=atom&namespace=0
	APIUrl  string // http://www.mediawiki.org/w/api.php
}

var ()

func (this *SourceMediawiki) Get() *SingleSource {
}

func (this *SourceMediawiki) GetByName(name string) {
	res, err := http.Get(this.APIUrl + "?" + (url.Values{
		"format": {"json"},
		"action": {"query"},
		"prop":   {"revisions"},
		"rvprop": {"content"},
		"titles": {name},
	}).Encode())
	json.NewDecoder(res.Body).Decode(v)
	defer res.Body.Close()
	return

}
