package main

import (
	"flag"
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

	flag.StringVar(&cfg.Address, "a", cfg.Address, "address to send metrics")
	flag.StringVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "report interval")
	flag.StringVar(&cfg.PollInterval, "p", cfg.PollInterval, "poll interval")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "key")
	flag.Parse()

	fmt.Println(cfg)
	cfg.Key = "lox"
	app.MonitorMetrics(cfg)
}
