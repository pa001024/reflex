package source

import (
	"encoding/xml"
	"github.com/pa001024/MoeCron/util"
	"net/http"
	"time"
)

// "format":"atom"
type FeedAtom struct {
	XMLName     xml.Name         `xml:"feed"`
	Id          string           `xml:"id"`
	Title       string           `xml:"title"`
	Updated     string           `xml:"updated"`
	Description string           `xml:"subtitle"`
	Generator   string           `xml:"generator"`
	Entry       []*FeedAtomEntry `xml:"entry"`
}
type FeedAtomEntry struct {
	Id      string `xml:"id"`
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Updated string `xml:"updated"`
	Summary string `xml:"summary"`
	Author  string `xml:"author>name"`
}

type SourceAtom struct { // Atom 实现接口ISource
	ISource
	Source

	FeedUrl string `json:"feed_url"` // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=atom&namespace=0
}

func (this *SourceAtom) GetChan() <-chan []*FeedInfo {
	return this.super_GetChan()
}

func (this *SourceAtom) Get() (rst []*FeedInfo) {
	f := this.FetchFeed()
	if f == nil {
		return
	}
	last := this.LastUpdate
	fetched := 0
	rst = make([]*FeedInfo, 0)
	for _, v := range f.Entry {
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
		fv := this.GetByFeedAtomEntry(v)
		rst = append(rst, fv)
		fetched++
	}
	return
}
func (this *SourceAtom) GetByFeedAtomEntry(v *FeedAtomEntry) (rst *FeedInfo) {
	rst = &FeedInfo{
		Id:       v.Id,
		SourceId: this.Name,
		Title:    v.Title,
		Author:   v.Author,
		Link:     v.Link,
		Content:  v.Summary,
	}
	return
}

func (this *SourceAtom) FetchFeed() (rst *FeedAtom) {
	res, err := http.Get(this.FeedUrl)
	if err != nil {
		util.Log("FetchFeed Fail")
		return
	}
	defer res.Body.Close()
	rst = &FeedAtom{}
	xml.NewDecoder(res.Body).Decode(rst)
	return
}
