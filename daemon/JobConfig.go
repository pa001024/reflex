package daemon

import (
	"encoding/json"
	"github.com/pa001024/MoeWorker/filter"
	"github.com/pa001024/MoeWorker/source"
	"github.com/pa001024/MoeWorker/target"
	"github.com/pa001024/MoeWorker/util"
	"io"
)

type JobConfig struct {
	Debug bool `json:"debug"`
	// 注: 用interface{}是为了读取所有数据以进行实例化封装
	Source map[string]interface{} `json:"source"` // 源
	Filter map[string]interface{} `json:"filter"` // 过滤器
	Target map[string]interface{} `json:"target"` // 目标

	// 解析target后使用
	targetBacklinks map[string][]target.ITarget `json:"-"`
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
		if this.Source[i] == nil {
			util.Log("source." + i + ".type \"" + v.(map[string]interface{})["type"].(string) + "\" not exists, skip")
			delete(this.Source, i)
			continue
		}
	}
	for i, v := range this.Filter {
		b, _ := json.Marshal(v)
		this.Filter[i] = filter.New(i, b)
		if this.Filter[i] == nil {
			util.Log("filter." + i + ".type \"" + v.(map[string]interface{})["type"].(string) + "\" not exists, skip")
			delete(this.Filter, i)
			continue
		}
	}
	for i, v := range this.Target {
		b, _ := json.Marshal(v)
		this.Target[i] = target.New(i, b)
		if this.Target[i] == nil {
			util.Log("target." + i + ".type \"" + v.(map[string]interface{})["type"].(string) + "\" not exists, skip")
			delete(this.Target, i)
			continue
		}
	}

	// 解析targetBacklinks
	this.targetBacklinks = make(map[string][]target.ITarget)
	for _, tv := range this.Target {
		for _, mv := range tv.(target.ITarget).GetMethod() {
			if mv.Action == "repost" {
				for _, sv := range mv.Source { // 类型为 "repost" 的 Method 中的Source是需要转发的Target的ID
					if this.targetBacklinks[sv] == nil {
						this.targetBacklinks[sv] = make([]target.ITarget, 0)
					}
					this.targetBacklinks[sv] = append(this.targetBacklinks[sv], tv.(target.ITarget))
				}
			}
		}
	}
	return
}

func (this *JobConfig) Save(w io.Writer) (err error) {
	json.NewEncoder(w).Encode(this)
	return
}
