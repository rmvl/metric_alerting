package app

type AgentConfig struct {
	Address        string `env:"ADDRESS" envDefault:"http://localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

type ServerConfig struct {
	Address       string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"300"`
	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool   `env:"RESTORE" envDefault:"true"`
}

func init() {

}

//AssetsPath := flag.String("assets", os.Getenv("ASSETS_PATH"), "a string")
//func (a ServerConfig) String() string {
//	return a.Host + ":" + strconv.Itoa(a.Port)
//}
//func (a *NetAddress) Set(s string) error {
//	hp := strings.Split(s, ":")
//	if len(hp) != 2 {
//		return errors.New("Need address in a form host:port")
//	}
//	port, err := strconv.Atoi(hp[1])
//	if err != nil {
//		return err
//	}
//	a.Host = hp[0]
//	a.Port = port
//	return nil
//}
