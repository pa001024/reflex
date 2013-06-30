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
	Debug      bool             `json:"debug"`
	Source     []source.ISource `json:"-"` // 源
	Filter     []filter.IFilter `json:"-"` // 过滤器
	Target     []target.ITarget `json:"-"` // 目标
	DataSource []interface{}    `json:"source"`
	DataFilter []interface{}    `json:"filter"`
	DataTarget []interface{}    `json:"target"`
}

func (this *JobConfig) Load(r io.Reader) (err error) {
	err = json.NewDecoder(r).Decode(this)
	if err != nil {
		util.Log("Parse Config Fail", err)
		return
	}
	this.Source = make([]source.ISource, len(this.DataSource))
	for i, v := range this.DataSource {
		b, _ := json.Marshal(v)
		this.Source[i] = source.New(b)
	}
	this.DataSource = nil

	this.Filter = make([]filter.IFilter, len(this.DataFilter))
	for i, v := range this.DataFilter {
		b, _ := json.Marshal(v)
		this.Filter[i] = filter.New(b)
	}
	this.DataFilter = nil
	this.Target = make([]target.ITarget, len(this.DataTarget))
	for i, v := range this.DataTarget {
		b, _ := json.Marshal(v)
		this.Target[i] = target.New(b)
	}
	this.DataTarget = nil
	return
}

func (this *JobConfig) Save(w io.Writer) (err error) {
	this.DataSource = make([]interface{}, len(this.Source))
	for i, v := range this.Source {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, this.DataSource[i])
	}
	this.DataFilter = make([]interface{}, len(this.Filter))
	for i, v := range this.Filter {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, this.DataFilter[i])
	}
	this.DataTarget = make([]interface{}, len(this.Target))
	for i, v := range this.Target {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, this.DataTarget[i])
	}
	json.NewEncoder(w).Encode(this)
	return
}
