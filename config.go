package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JobConfig struct {
	Source []*JobSource `json:"source"` // 源
	Filter []*JobFilter `json:"filter"` // 过滤器
	Target []*JobTarget `json:"target"` // 目标
}

type JobSource struct {
	Type      string `json:"type"`      // 类型 source工厂ID 如[mediawiki,rss,atom]等
	Format    string `json:"format"`    // encoding工厂ID 如[xml,json]
	Name      string `json:"name"`      // 名字
	Limit     string `json:"limit"`     // 发送频次 格式[次数/时间] 比如1/1为每分钟最多发送一次 超过的丢弃 (iron.io上无效)
	EnablePic string `json:"enablepic"` // 启用图片 如果target支持图片会一并发出
	PicTag    string `json:"pictag"`    // 图片标签(CSS选择器) 一般为 img 亦或指定ID或class 使用 img.class#xxx 等
	PicSize   []int  `json:"picsize"`   // 图片大小 [minX,maxX,minY,maxY] 不填表示不限制
}

type JobFilter struct {
	Type string `json:"type"` // 类型 filter工厂ID 如[moegirlwiki,简繁转换]等
	Name string `json:"name"` // 名字 不能跟别的target或source名字相同
}

type JobTarget struct {
	Type   string             `json:"type"`   // 类型 target工厂ID 如[sina,qq,renren]等
	Name   string             `json:"name"`   // 名字 不能跟别的target或source名字相同
	Method []*JobTargetMethod `json:"method"` // 处理方法
}

type JobTargetMethod struct {
	Action string   `json:"action"` // 动作 可选[update,upload]
	Source []string `json:"source"` // 目标 填写相应source或target名字
	Filter []string `json:"filter"` // 过滤器 从左到右依次管道
}

func (this *JobConfig) Load(r io.Reader) (err error) {
	json.NewDecoder(r).Decode(this)
	return
}
