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

	APIUrl string `json:"api_url"` // http://www.mediawiki.org/w/api.php
}

func (this *SourceMediawiki) Get() (rst []*FeedInfo) {
	f := this.FetchFeed()
	if f == nil {
		return
	}
	last := time.Now().Add(-time.Duration(this.Interval) * time.Second)
	rst = make([]*FeedInfo, 0)
	for _, v := range f.Item {
		d, err := time.Parse(time.RFC1123, v.Updated)
		if err != nil {
			util.Log("Time Parse Fail", err)
			continue
		}
		if d.Before(this.LastUpdate) {
			break
		}
		if d.After(last) {
			last = d
		}
		f := this.GetByFeedRSSItem(v)
		rst = append(rst, f)
	}
	this.LastUpdate = last
	return
}

func (this *SourceMediawiki) FetchFeed() (rst *FeedRSS) {
	res, err := http.Get(this.FeedUrl)
	if err != nil {
		util.Log("FetchFeed Fail")
	}
	defer res.Body.Close()
	rst = &FeedRSS{}
	xml.NewDecoder(res.Body).Decode(rst)
	return
}

func (this *SourceMediawiki) GetByFeedRSSItem(v *FeedRSSItem) (rst *FeedInfo) {
	rst = &FeedInfo{
		Id:       v.Id,
		SourceId: this.Name,
		Title:    v.Title,
		Author:   v.Author,
		Content:  this.GetByName(v.Title),
	}
	return
}

func (this *SourceMediawiki) GetByName(name string) (rst string) {
	res, err := http.Get(this.APIUrl + "?" + (url.Values{
		"format": {"json"},
		"action": {"query"},
		"prop":   {"revisions"},
		"rvprop": {"content"},
		"titles": {name},
	}).Encode())
	if err != nil {
		util.Log("Network Fetch Fail", err)
	}
	defer res.Body.Close()
	var v map[string]map[string]map[string]map[string][]map[string]interface{}
	json.NewDecoder(res.Body).Decode(v)
	for _, v1 := range v["query"]["pages"] {
		rv := v1["revisions"]
		if rv != nil && len(rv) > 0 {
			rst = rv[0]["*"].(string)
			break
		}
	}
	return
}
