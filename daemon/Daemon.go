package daemon

import (
	"github.com/pa001024/MoeCron/filter"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/target"
	"github.com/pa001024/MoeCron/util"
	"time"
)

type Daemon struct {
	Config *JobConfig
	C      chan []*source.FeedInfo
}

func NewDaemon(conf *JobConfig) (rst *Daemon) {
	if conf == nil {
		panic("NO CONFIG")
	}
	rst = &Daemon{
		Config: conf,
		C:      make(chan []*source.FeedInfo),
	}
	return
}

/*
Daemon程序处理流程:

获取到[]FeedInfo
  ↓          ↓
方法A      方法B → ...
  ↓
过滤器 →  目标B(转发A)
  ↓
目标A(发表) -[回调]-↑

*/

func (this *Daemon) Serve() {
	for _, v := range this.Config.Source { // 异步(并行)获取source
		go func() {
			for {
				this.C <- <-v.(source.ISource).GetChan()
			}
		}()
	}
	for {
		src := <-this.C
		for _, tv := range this.Config.Target {
			for _, mv := range tv.(target.ITarget).GetMethod() {
				go this.DoAction(tv.(target.ITarget), mv, src)
			}
		}
	}
}

// 处理TargetMethod.action请求类型
func (this *Daemon) DoAction(tv target.ITarget, act *target.TargetMethod, src []*source.FeedInfo) bool {
	// 过一遍过滤器
	nc := src
	for _, fv := range act.Filter {
		f := this.Config.Filter[fv]
		if f == nil {
			util.Log("Warning: filter \"" + fv + "\" not exists.")
			continue
		}
		nc = f.(filter.IFilter).Process(nc)
	}
	if act.Action == "update" {
		for _, sv := range act.Source {
			if this.Config.Source[sv] == nil {
				util.Log("Warning: update.source \"" + sv + "\" not exists. ")
				continue
			}
			for _, rv := range nc {
				b, err := tv.Send(rv)
				if b == "" {
					util.Log("Master send fail, stop repost. ", err)
					continue
				}
				rp := *rv
				rp.RepostId = b
				rp.Content = "转发微博"
				for _, rtv := range this.Config.targetBacklinks[tv.GetId()] {
					<-time.After(20 * time.Second) // 防止刷爆 TODO: 移到配置文件
					rtv.Send(&rp)
				}
			}
		}
	}
	return false
}
