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
	"encoding/xml"
	"github.com/pa001024/MoeCron/util"
	"net/http"
	"net/url"
	"time"
)

type SourceMediawiki struct { // Mediawiki 实现接口ISource
	ISource
	SourceRSS

	FeedUrl    string    `json:"feed_url"` // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=rss&namespace=0
	APIUrl     string    `json:"api_url"`  // http://www.mediawiki.org/w/api.php
	LastUpdate time.Time `json:"lastupdate"`
}

func (this *SourceMediawiki) Get() (rst []*FeedInfo) {
	f := this.FetchFeed()
	for _, v := range f.Item {
		d, _ := time.Parse(time.RFC1123, v.Updated)
		if d.Before(time.Now().Add(time.Duration(this.Interval) * -time.Second)) {
			break
		}
	}
	return
}

func (this *SourceMediawiki) FetchFeed() (rst *FeedRSS) {
	res, err := http.Get(this.FeedUrl)
	if err != nil {
		util.Log("FetchFeed fail")
	}
	defer res.Body.Close()
	rst = &FeedRSS{}
	xml.NewDecoder(res.Body).Decode(rst)
	return
}

func (this *SourceMediawiki) GetByName(name string) (rst *FeedInfo) {
	res, err := http.Get(this.APIUrl + "?" + (url.Values{
		"format": {"json"},
		"action": {"query"},
		"prop":   {"revisions"},
		"rvprop": {"content"},
		"titles": {name},
	}).Encode())
	if err != nil {
		util.Log("GetByName fail")
	}
	defer res.Body.Close()
	rst = &FeedInfo{}
	json.NewDecoder(res.Body).Decode(rst)
	return
}
