package main

import (
	"github.com/pa001024/MoeCron/daemon"
	"github.com/pa001024/MoeCron/util"
	"os"
)

var (
	conf *daemon.JobConfig
)

func main() {
	conf = &daemon.JobConfig{}
	r, err := os.Open("config.json")
	if err != nil {
		util.Log("Cound not Load config.json")
		return
	}
	conf.Load(r)
	daemon.NewDaemon(conf).Serve()
}
