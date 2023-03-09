package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"yalerting/cmd/app"
)

func main() {
	var cfg app.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "address to send metrics")
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "report interval")
	flag.IntVar(&cfg.PollInterval, "p", cfg.PollInterval, "poll interval")
	flag.Parse()

	app.MonitorMetrics(cfg)
}
