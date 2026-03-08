package settings

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

// Options - структура для хранения настроек сервиса.
type Options struct {
	ServerAddress string `json:"server_address"`
	BaseURL       string `json:"base_url"`
	DatabaseDSN   string `json:"database_dsn"`
	Config        string
	LogLevel      string
}

// ParseFlags - парсит флаги командной строки или переменные окружения.
// Результат сохраняет в структуру Options.
func ParseFlags(o *Options) {
	//flag.StringVar(&o.Config, "c", "config.json", "config path")
	flag.StringVar(&o.Config, "c", "", "config path")
	if config := os.Getenv("CONFIG"); config != "" {
		o.Config = config
	}

	var config Options
	var err error
	if o.Config != "" {
		config, err = readConfig(o.Config)
		if err != nil {
			panic(err)
		}
	}

	fillDefaultOptions(o)
	overrideOptionsFromConfig(o, &config)
	overrideOptionsFromCmd(o)
	overrideOptionsFromEnv(o)

}

func overrideOptionsFromConfig(o *Options, c *Options) {
	if c.ServerAddress != "" {
		o.ServerAddress = c.ServerAddress
	}
	if c.BaseURL != "" {
		o.BaseURL = c.BaseURL
	}
	if c.LogLevel != "" {
		o.LogLevel = c.LogLevel
	}
	if c.DatabaseDSN != "" {
		o.DatabaseDSN = c.DatabaseDSN
	}
}

func readConfig(fname string) (Options, error) {
	var config Options

	f, err := os.Open(fname)
	if err != nil {
		return config, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err

}

func fillDefaultOptions(o *Options) {
	o.ServerAddress = ":8080"
	o.BaseURL = "http://localhost:8080"
	o.LogLevel = "debug"
	o.DatabaseDSN = "host=localhost user=postgres password=xxxx dbname=gophkeeper sslmode=disable"
}

func overrideOptionsFromCmd(o *Options) {
	flag.StringVar(&o.ServerAddress, "a", o.ServerAddress, "address and port to run server")
	flag.StringVar(&o.BaseURL, "b", o.BaseURL, "base address for short URL")
	flag.StringVar(&o.LogLevel, "l", o.LogLevel, "log level")
	flag.StringVar(&o.DatabaseDSN, "d", o.DatabaseDSN, "database connection string")
	flag.Parse()
}

func overrideOptionsFromEnv(o *Options) {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		o.ServerAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		o.BaseURL = baseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		o.LogLevel = envLogLevel
	}
	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		o.DatabaseDSN = databaseDSN
	}
}
