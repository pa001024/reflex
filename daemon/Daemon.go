package daemon

import (
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/util"
)

type Daemon struct {
	Config *JobConfig
	C      chan []*source.FeedInfo
}

func NewDaemon(conf *JobConfig) (rst *Daemon) {
	rst = &Daemon{
		Config: conf,
		C:      make(chan []*source.FeedInfo),
	}
	return
}

func (this *Daemon) Serve() {
	for _, v := range this.Config.Source {
		go func() {
			for {
				this.C <- v.(source.ISource).Get()
			}
		}()
	}
	for {
		for _, v := range <-this.C {
			util.Log(v.SourceId)
		}

	}
}
