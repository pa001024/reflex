package source

import (
	"encoding/json"
	"github.com/pa001024/reflex/util"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	rep_mw_pic = regexp.MustCompile(
		`` +
			`\[\[(?:[Ii][Mm][Aa][Gg][Ee]|[Ff][Ii][Ll][Ee]):\s*(\S.+?\.(?:jpg|png|gif))(?:\|.+)?\]\]|` +
			`\| *?(?:[I][Mm][Aa][Gg][Ee]_[Nn][Aa][Mm][Ee]|[Ii][Mm][Aa][Gg][Ee]) *?=\s*(\S.+\.(?:jpg|png|gif))|` +
			`(https?://[A-z][A-z0-9\-\./]+/[A-z0-9\-\./%]+\.(?:jpg|png|gif))`) // 提取图片
	rep_mw_redirect = regexp.MustCompile(`#(?:重定向|[Rr]edirect) \[\[(.+)\]\]`) // 处理重定向
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
	go func() {
		for {
			tc := time.After(time.Duration(this.Interval) * time.Second)
			chw <- this.Get()
			<-tc
		}
	}()
	this.C = chw
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
			util.WARN.Err("Time Parse Fail", err)
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
	t := rep_mw_pic.FindAllStringSubmatch(src, -1) // TODO: 有的WIKI有两个以上的图片地址 需要补完算法
	dst = make([]string, 0, len(t))
endf:
	for _, v := range t {
		var raw string
		for i, s := range v[1:] {
			if s != "" {
				if i == 2 {
					dst = append(dst, s)
					continue endf
				}
				raw = s
				break
			}
		}

		if raw == "" {
			continue
		}
		raw = strings.Replace(raw, " ", "_", -1)
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
		"titles": {strings.Replace(name, " ", "_", -1)},
	}).Encode())
	if err != nil {
		util.ERROR.Err("Network Fetch Fail", err)
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
