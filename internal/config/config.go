// Package config - application configs structures
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config - application config
type Config struct {
	Database *DatabaseConfig `mapstructure:"database"`
	Server   *ServerConfig   `mapstructure:"server"`
	WebAPI   *WebAPI         `mapstructure:"web_api"`
}

// DatabaseConfig database connection preferences
type DatabaseConfig struct {
	ConnString string `mapstructure:"conn_string"`
}

// WebAPI contains config for web api
type WebAPI struct {
	InternalAPI *APIConf `mapstructure:"internal"`
}

// APIConf - config of single server in web api
type APIConf struct {
	Addr string `mapstructure:"addr"`
	Cors *Cors  `mapstructure:"cors"`
}

// Cors contains cors-policy settings
type Cors struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	Debug          bool     `mapstructure:"debug"`
	AllowCreds     bool     `mapstructure:"allow_creds"`
}

// New - initialize app configuration
func New() (*Config, error) {
	var config = &Config{}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); nil != err {
		return nil, fmt.Errorf("unable to read config from file")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.Unmarshal(config)
	if nil != err {
		return nil, fmt.Errorf("unable to decode into struc: '%w'", err)
	}

	return config, nil
}
