package configuration

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type HttpConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MySQLConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	Database          string `mapstructure:"database"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	ConnMaxLifetimeMs int64  `mapstructure:"conn_max_lifetime_ms"`
	MaxOpenConns      int    `mapstructure:"max_open_conns"`
	MaxIdleConns      int    `mapstructure:"max_idle_conns"`
}

type Config struct {
	Http  HttpConfig  `mapstructure:"http"`
	MySQL MySQLConfig `mapstructure:"mysql"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	if err := godotenv.Load(path + "/.env"); err != nil {
		return cfg, err
	}

	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	cfg.MySQL.Password = os.Getenv("MYSQL_PASSWORD")

	return cfg, nil
}
