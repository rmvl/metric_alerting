package app

type AgentConfig struct {
	Address        string `env:"ADDRESS" envDefault:"http://localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

type ServerConfig struct {
	Address       string `env:"ADDRESS" envDefault:"http://localhost:8080"`
	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"300"`
	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool   `env:"RESTORE" envDefault:"true"`
}
