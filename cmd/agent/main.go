package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
	"yalerting/cmd/app"
)

func main() {
	var cfg app.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	var address, reportInterval, pollInterval string
	flag.StringVar(&address, "a", cfg.Address, "address to send metrics")
	flag.StringVar(&reportInterval, "r", cfg.ReportInterval, "report interval")
	flag.StringVar(&pollInterval, "p", cfg.PollInterval, "poll interval")
	flag.Parse()

	_, present := os.LookupEnv("ADDRESS")
	if !present && len(address) > 0 {
		cfg.Address = address
	}
	_, present = os.LookupEnv("REPORT_INTERVAL")
	if !present && len(reportInterval) > 0 {
		cfg.ReportInterval = reportInterval
	}
	_, present = os.LookupEnv("POLL_INTERVAL")
	if !present && len(pollInterval) > 0 {
		cfg.PollInterval = pollInterval
	}

	app.MonitorMetrics(cfg)
}
