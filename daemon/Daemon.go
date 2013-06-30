package daemon

import (
	"github.com/pa001024/MoeCron/source"
)

var (
	D Daemon
)

type Daemon struct {
	Config *JobConfig
	C      chan []*source.Source
}

func NewDaemon(conf *JobConfig) (rst *Daemon) {
	rst = &Daemon{
		Config: conf,
		C:      make(chan []*source.Source),
	}
	return
}

func (this *Daemon) Serve() {
	for {
		// ...
	}
}
