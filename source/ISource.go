package source

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/util"
	"time"
)

type ISource interface {
	Get() []*FeedInfo
}

type FeedInfo struct { // 目标
	Id       string   `json:"id"`
	SourceId string   `json:"source"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Content  string   `json:"content"`
	PicUrl   []string `json:"picurl"`
}

type Source struct { // 配置持久模板
	Type     string `json:"type"`          // 类型 source工厂ID 如[mediawiki,rss,atom]等
	Name     string `json:"name"`          // 名字
	Interval int    `json:"interval"`      // 发送频率 单位为秒
	Limit    int    `json:"limit"`         // 单次发送上限 超过的丢弃
	Pic      []int  `json:"pic,omitempty"` // 图片大小 [minX,maxX,minY,maxY] 不填表示不发送图片
}

func New(name string, b []byte) (rst ISource) {
	obj := &Source{}
	err := json.Unmarshal(b, obj)
	if err != nil {
		util.Log("JSON Parse Error", err)
		return
	}
	switch obj.Type {
	default:
	case "mediawiki":
		rst = &SourceMediawiki{}
		json.Unmarshal(b, rst)
		rst.(*SourceMediawiki).Name = name
		rst.(*SourceMediawiki).LastUpdate = time.Now()
		break
		// case "rss":
	}
	return
}
