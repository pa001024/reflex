package main

import (
	"github.com/pa001024/MoeWorker/daemon"
	"github.com/pa001024/MoeWorker/util"
	"os"
)

var (
	conf = &daemon.JobConfig{}
)

func main() {
	reload()
	daemon.NewDaemon(conf).Serve()
}

func reload() {
	r, err := os.Open("config.json")
	if err != nil {
		util.ERROR.Err("Cound not Load config.json")
		return
	}
	conf.Load(r)
}
