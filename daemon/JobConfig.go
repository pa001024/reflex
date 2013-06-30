package daemon

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/filter"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/target"
	"io"
)

type JobConfig struct {
	Debug  bool             `json:"debug"`
	Source []*source.Source `json:"source"` // 源
	Filter []*filter.Filter `json:"filter"` // 过滤器
	Target []*target.Target `json:"target"` // 目标
}

func (this *JobConfig) Load(r io.Reader) (err error) {
	json.NewDecoder(r).Decode(this)
	return
}
