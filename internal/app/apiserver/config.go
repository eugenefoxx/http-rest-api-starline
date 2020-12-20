package apiserver

//import "github.com/eugenefoxx/http-rest-api/internal/store"

// Config ...
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
	//	HTML        string `toml:"html"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "debug",
	}
}
