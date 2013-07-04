/*
简介
=======
Mediawiki 的 ISource 接口实现

功能
--------

1. 获取文字内容(使用wiki API提取源代码)
2. 将图片和文字转交给filter进行文字转换裁剪/添加水印操作后提交给Target进行最终的上传/发送/转发等

注意
--------

1. 这里的方法都是同步的 如担心阻塞请使用异步方式调用

TODO
--------

这里写在source可能有代码冗余 请转移到filter

*/
package source

import (
	"encoding/json"
	"encoding/xml"
	"github.com/pa001024/MoeCron/util"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	rep_mw_round0   = regexp.MustCompile(`\[\[File:\s*(\S.+?)\|.+\]\]|https?://[A-z][A-z0-9\-\./]+/[A-z0-9\-\./%]+\.(?:jpg|png|gif)|\| image ?=\s*(\S.+.(?:jpg|png|gif))`) // 提取图片
	rep_mw_redirect = regexp.MustCompile(`#(?:重定向|redirect) \[\[(.+)\]\]`)                                                                                                 // 处理重定向
)

type SourceMediawiki struct { // Mediawiki 实现接口ISource
	ISource
	SourceRSS

	APIUrl  string `json:"api_url"` // http://www.mediawiki.org/w/api.php
	PicBase string `json:"pic_url"` // http://upload.wikimedia.org/wikipedia/commons/
}

func (this *SourceMediawiki) GetChan() <-chan []*FeedInfo {
	if this.C != nil {
		return this.C
	}
	chw := make(chan []*FeedInfo)
	t := time.NewTimer(time.Duration(this.Interval) * time.Second)
	go func() {
		<-t.C
		chw <- this.Get()
	}()
	return chw
}

func (this *SourceMediawiki) Get() (rst []*FeedInfo) {
	f := this.FetchFeed()
	if f == nil {
		return
	}
	last := this.LastUpdate
	fetched := 0
	rst = make([]*FeedInfo, 0)
	for _, v := range f.Item {
		if fetched >= this.Limit {
			break
		}
		d, err := time.Parse(time.RFC1123, v.Updated)
		if err != nil {
			util.Log("Time Parse Fail", err)
			continue
		}
		if d.Sub(last) <= 0 { // It means if feed.Updated >= this.LastUpdate
			break
		}
		if d.After(this.LastUpdate) {
			this.LastUpdate = d
		}
		fv := this.GetByFeedRSSItem(v)
		rst = append(rst, fv)
		fetched++
	}
	return
}

func (this *SourceMediawiki) FetchFeed() (rst *FeedRSS) {
	res, err := http.Get(this.FeedUrl)
	if err != nil {
		util.Log("FetchFeed Fail", err)
		return
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
		Link:     v.Link,
		Content:  this.GetByName(v.Title),
	}
	this.ClearRedirect(rst)                     // 处理重定向
	rst.PicUrl = this.FilterPicUrl(rst.Content) // 处理图片
	return
}
func (this *SourceMediawiki) ClearRedirect(v *FeedInfo) {
	d := rep_mw_redirect.FindAllStringSubmatch(v.Content, 1)
	if len(d) > 0 && len(d[0]) > 1 {
		v.Content = this.GetByName(d[0][1])
	}
}

func (this *SourceMediawiki) FilterPicUrl(src string) (dst []string) {
	t := rep_mw_round0.FindAllStringSubmatch(src, 0)
	dst = make([]string, 0, len(t))
	for _, v := range t {
		var raw string
		for _, s := range v {
			if s != "" {
				raw = s
			}
		}
		if raw == "" {
			continue
		}
		raw = strings.Replace(raw, " ", "_", 0)
		h := util.Md5String(raw)
		raw = this.PicBase + h[0:1] + "/" + h[0:2] + "/" + raw
		dst = append(dst, raw)
	}
	return
}

func (this *SourceMediawiki) GetByName(name string) (rst string) {
	res, err := http.Get(this.APIUrl + "?" + (url.Values{
		"format": {"json"},
		"action": {"query"},
		"prop":   {"revisions"},
		"rvprop": {"content"},
		"titles": {strings.Replace(name, " ", "_", 0)},
	}).Encode())
	if err != nil {
		util.Log("Network Fetch Fail", err)
		return
	}
	defer res.Body.Close()

	var v map[string]map[string]map[string]map[string]interface{}
	json.NewDecoder(res.Body).Decode(&v)
	for _, v1 := range v["query"]["pages"] {
		rv := v1["revisions"]
		if rv != nil && len(rv.([]interface{})) > 0 && rv.([]interface{})[0] != nil {
			rst = rv.([]interface{})[0].(map[string]interface{})["*"].(string)
			break
		}
	}
	return
}
