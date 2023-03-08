package app

type ServerConfig struct {
	Address       string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"300"`
	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool   `env:"RESTORE" envDefault:"true"`
}
