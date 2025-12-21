// Package infra provides infrastructure components including configuration, database, logging, and routing.
package infra

type Config struct {
	App struct {
		Addr string
	}
	DB struct {
		Driver string
		DSN    string
	}
	Log struct {
		Level string
	}
}

func LoadConfig() *Config {
	cfg := &Config{}
	cfg.App.Addr = ":18080"
	cfg.DB.Driver = "sqlite3"
	cfg.DB.DSN = "./data.db?_loc=auto"
	cfg.Log.Level = "info"
	return cfg
}
