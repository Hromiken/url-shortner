package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppInfo  AppInfo        `yaml:"app" mapstructure:"app"`
	Server   ServerConfig   `yaml:"http_server" mapstructure:"http_server"`
	Logger   LoggerConfig   `yaml:"logger_config" mapstructure:"logger_config"`
	Postgres PostgresConfig `yaml:"postgres_config" mapstructure:"postgres_config"`
	Redis    RedisConfig    `yaml:"redis_config" mapstructure:"redis_config"`
}

type AppInfo struct {
	Name    string `yaml:"name" mapstructure:"name"`
	Version string `yaml:"version" mapstructure:"version"`
}
type ServerConfig struct {
	Addr        string        `yaml:"addr" mapstructure:"addr"`
	Timeout     time.Duration `yaml:"timeout" mapstructure:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

type LoggerConfig struct {
	LogLevel string `yaml:"log_level" mapstructure:"log_level"`
}

type PostgresConfig struct {
	DSN             string        `yaml:"dsn" mapstructure:"dsn"`
	SlavesDSN       []string      `yaml:"slaves_dsn" mapstructure:"slaves_dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_сonns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	Port            int           `yaml:"port" mapstructure:"port"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

func NewAppConfig() (*AppConfig, error) {
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("CONFIG ERROR:", err) // добавь это
		return nil, err
	}
	fmt.Println("Used config:", viper.ConfigFileUsed())

	var appCfg AppConfig
	err = viper.Unmarshal(&appCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &appCfg, nil
}
