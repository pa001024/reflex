package source

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/util"
)

type ISource interface {
	Get() []*FeedInfo
}

type FeedInfo struct { // 目标
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	PicUrl  []string `json:"picurl"`
}

type Source struct { // 配置持久模板
	Type      string `json:"type"`      // 类型 source工厂ID 如[mediawiki,rss,atom]等
	Format    string `json:"format"`    // encoding工厂ID 如[xml,json]
	Name      string `json:"name"`      // 名字
	Interval  int    `json:"interval"`  // 发送频率 单位为秒
	Limit     int    `json:"limit"`     // 单次发送上限 超过的丢弃
	EnablePic string `json:"enablepic"` // 启用图片 如果target支持图片会一并发出
	PicTag    string `json:"pictag"`    // 图片标签(CSS选择器) 一般为 img 亦或指定ID或class 使用 img.class#xxx 等
	PicSize   []int  `json:"picsize"`   // 图片大小 [minX,maxX,minY,maxY] 不填表示不限制
}

func New(b []byte) (rst ISource) {
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
		break
		// case "rss":
	}
	return
}
