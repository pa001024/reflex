package daemon

import (
	"encoding/json"
	"github.com/pa001024/MoeCron/filter"
	"github.com/pa001024/MoeCron/source"
	"github.com/pa001024/MoeCron/target"
	"github.com/pa001024/MoeCron/util"
	"io"
)

type JobConfig struct {
	Debug  bool                   `json:"debug"`
	Source map[string]interface{} `json:"source"` // 源
	Filter map[string]interface{} `json:"filter"` // 过滤器
	Target map[string]interface{} `json:"target"` // 目标
}

func (this *JobConfig) Load(r io.Reader) (err error) {
	err = json.NewDecoder(r).Decode(this)
	if err != nil {
		util.Log("Parse Config Fail", err)
		return
	}
	for i, v := range this.Source {
		b, _ := json.Marshal(v)
		this.Source[i] = source.New(i, b)
	}
	for i, v := range this.Filter {
		b, _ := json.Marshal(v)
		this.Filter[i] = filter.New(i, b)
	}
	for i, v := range this.Target {
		b, _ := json.Marshal(v)
		this.Target[i] = target.New(i, b)
	}
	return
}

func (this *JobConfig) Save(w io.Writer) (err error) {
	json.NewEncoder(w).Encode(this)
	return
}
