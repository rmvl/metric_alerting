package app

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
