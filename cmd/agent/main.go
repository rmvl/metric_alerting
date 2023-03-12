package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"yalerting/cmd/app"
)

func main() {
	var cfg app.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	app.MonitorMetrics(cfg)
}
