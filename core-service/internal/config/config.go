package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string         `yaml:"env" env-default:"local"`
	Database       DatabaseConfig `yaml:"database" env-required:"true"`
	HTTP           HTTPConfig     `yaml:"http"`
	S3             S3Config       `yaml:"s3"`
	MigrationsPath string
}

type S3Config struct {
	EndpointUrl     string `yaml:"endpoint_url"`
	Database        string `yaml:"database"`
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" env-required:"true"`
}

type DatabaseConfig struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"user"`
	Password string `env:"PG_PASSWORD" env-required:"true"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check whether file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist" + configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty" + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to fetch necessary env variables" + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or enviroment variable.
// Priority: flag > env > default/
// Default value is empty string
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to the config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
