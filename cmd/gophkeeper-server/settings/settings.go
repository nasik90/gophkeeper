package settings

import (
	"encoding/json"
	"flag"
	"os"
)

// Options - структура для хранения настроек сервиса.
type Options struct {
	ServerAddress string `json:"server_address"`
	BaseURL       string `json:"base_url"`
	DatabaseDSN   string `json:"database_dsn"`
	Config        string
}

// ParseFlags - парсит флаги командной строки или переменные окружения.
// Результат сохраняет в структуру Options.
func ParseFlags(o *Options) {
	//flag.StringVar(&o.Config, "c", "config.json", "config path")
	flag.StringVar(&o.Config, "c", "", "config path")
	if config := os.Getenv("CONFIG"); config != "" {
		o.Config = config
	}

	//var config Options
	// var err error
	// if o.Config != "" {
	// 	config, err = readConfig(o.Config)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	fillDefaultOptions(o)
	// overrideOptionsFromConfig(o, &config)
	// overrideOptionsFromCmd(o)
	// overrideOptionsFromEnv(o)

}

func readConfig(fname string) (Options, error) {
	var config Options

	f, err := os.Open(fname)
	if err != nil {
		return config, err
	}
	out := make([]byte, 1024)
	var n int
	if n, err = f.Read(out); err != nil {
		return config, err
	}
	data := out[:n]
	err = json.Unmarshal(data, &config)
	return config, err

}

func fillDefaultOptions(o *Options) {
	o.ServerAddress = ":8080"
	o.BaseURL = "http://localhost:8080"
	//o.LogLevel = "debug"
	o.DatabaseDSN = "host=localhost user=postgres password=xxxx dbname=gophkeeper sslmode=disable"
	//o.DatabaseDSN = ""
}
