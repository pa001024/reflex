package source

import (
	"encoding/json"
	"github.com/pa001024/reflex/util"
	"time"
)

type ISource interface {
	Get() []*FeedInfo
	GetChan() <-chan []*FeedInfo
}

type FeedInfo struct { // 目标
	Id       string   `json:"id"`
	SourceId string   `json:"source"`
	RepostId string   `json:"rid"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Content  string   `json:"content"`
	PicUrl   []string `json:"picurl"`
	Link     string   `json:"link"`
}

type Source struct { // 配置持久模板
	Type     string `json:"type"`     // 类型 source工厂ID 如[mediawiki,rss,atom]等
	Name     string `json:"name"`     // 名字
	Interval int    `json:"interval"` // 发送频率 单位为秒
	Limit    int    `json:"limit"`    // 单次发送上限 超过的丢弃
	// Pic      []int  `json:"pic,omitempty"` // 图片大小 [minX,maxX,minY,maxY] 不填表示不发送图片
	C                   chan []*FeedInfo `json:"-"`
	LastUpdate          time.Time        `json:"-"`
	DoUpdateBeforeStart bool             `json:"do_update_before_start"` // 不更新程序启动之前的条目
}

func (this *Source) GetId() (rst string) { return this.Name }

func New(name string, b []byte) (rst ISource) {
	obj := &Source{}
	err := json.Unmarshal(b, obj)
	if err != nil {
		util.ERROR.Err("JSON Parse Error", err)
		return
	}
	obj.Name = name
	if !obj.DoUpdateBeforeStart {
		obj.LastUpdate = time.Now()
	}
	switch obj.Type { // TODO: 使用反射替代这么啰嗦的初始化
	case "atom", "atomfeed":
		dst := &SourceAtom{}
		json.Unmarshal(b, dst)
		rst = dst
	case "rss", "rssfeed":
		dst := &SourceRSS{}
		json.Unmarshal(b, dst)
		rst = dst
	case "mediawiki", "wikifeed", "wiki":
		dst := &SourceMediawiki{}
		json.Unmarshal(b, dst)
		rst = dst
	default:
		util.WARN.Logf("source.%s \"%s\" not exists.", obj.Type, name)
		return
	}
	util.INFO.Logf("source..%s \"%s\" Loaded.", obj.Type, name)
	return
}
