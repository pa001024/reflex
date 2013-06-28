package main

import (
	"github.com/pa001024/MoeCron/util"
	"os"
)

var (
	conf *JobConfig
)

func main() {
	conf = &JobConfig{}
	r, err := os.Open("config.json")
	if err != nil {
		util.Log("Cound not Load config.json")
		return
	}
	conf.Load(r)
	Start(conf)
}
