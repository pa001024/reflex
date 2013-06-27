package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JobConfig struct {
	Source []*JobSource `json:"source"`
	Target []*JobTarget `json:"target"`
}

type JobSource struct {
	Type      string `json:"type"`
	Format    string `json:"format"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	EnablePic string `json:"enablepic"`
	PicTag    string `json:"pictag"`
	PicSize   []int  `json:"picsize"`
}

type JobTarget struct {
	Type   string             `json:"type"`
	Name   string             `json:"name"`
	Method []*JobTargetMethod `json:"method"`
	// 微博部分
	AppKey      string    `json:client_id`     // AppKey
	AppSecret   string    `json:client_secret` // AppSecret
	CallbackUrl string    `json:redirect_uri`  // 验证URL
	Token       string    `json:access_token`  // OAuth2.0 验证码
	ExpiresIn   time.Time `json:expires_in`    // 失效时间
}

type JobTargetMethod struct {
	Action string   `json:"action"`
	Source []string `json:"source"`
}

func (this *JobConfig) Load(r io.Reader) (err error) {
	json.NewDecoder(r).Decode(this)
	return
}
