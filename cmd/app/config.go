package app

import (
	"fmt"
	"strconv"
	"strings"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval string `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ServerConfig struct {
	Address       string `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval string `env:"STORE_INTERVAL" envDefault:"11s"`
	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.2.json"`
	Restore       bool   `env:"RESTORE" envDefault:"true"`
}

func (cfg *ServerConfig) GetStoreInterval() int {
	if strings.HasSuffix(cfg.StoreInterval, "s") {
		storeInterval, err := strconv.Atoi(strings.TrimSuffix(cfg.StoreInterval, "s"))
		if err != nil {
			fmt.Println(err)
		}
		return storeInterval
	}
	if strings.HasSuffix(cfg.StoreInterval, "m") {
		storeInterval, err := strconv.Atoi(strings.TrimSuffix(cfg.StoreInterval, "m"))
		if err != nil {
			fmt.Println(err)
		}
		return storeInterval
	}
	storeInterval, err := strconv.Atoi(cfg.StoreInterval)
	if err != nil {
		fmt.Println(err)
	}
	return storeInterval
}

func (cfg *AgentConfig) GetReportInterval() int {
	reportInterval, _ := strconv.Atoi(strings.TrimSuffix(cfg.ReportInterval, "s"))
	return reportInterval
}

func (cfg *AgentConfig) GetPollInterval() int {
	pollInterval, _ := strconv.Atoi(strings.TrimSuffix(cfg.PollInterval, "s"))
	return pollInterval
}
