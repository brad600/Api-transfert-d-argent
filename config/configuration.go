package APIServer

// Config contient les données nécessaires au fonctionnement de la base de données et du serveur api
type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	DbPath       string `toml:"db_path"`
	QueryTimeout uint32 `toml:"query_timeout"`
}

// NewConfig instancie le nouvel obbjet de configuration
func NewConfig() *Config {
	return &Config{
		BindAddr:     ":8010",
		LogLevel:     "debug",
		DbPath:       "/tmp/mysql.db",
		QueryTimeout: 10,
	}
}
