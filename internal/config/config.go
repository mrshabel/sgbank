package config

type Config struct {
	ServerAddr string
	Env        ENV
}

type ENV string

const (
	DEV  ENV = "local"
	PROD ENV = "production"
)

func New() *Config {
	// TODO: move config variables into env file
	return &Config{
		ServerAddr: "localhost:8000",
		Env:        DEV,
	}
}
