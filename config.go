package main

import ()

type JobConfig struct {
	Name   string      `json:"name"`
	Source []JobSource `json:"source"`
	Target []JobTarget `json:"target"`
}

type JobSource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type JobTarget struct {
	Name string `json:"name"`
}
