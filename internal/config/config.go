package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceConfiguration   ServiceConfiguration   `mapstructure:"ServiceConfiguration"`
	PostgresConfiguration  PostgresConfiguration  `mapstructure:"PostgresConfiguration"`
	RedisConfiguration     RedisConfiguration     `mapstructure:"RedisConfiguration"`
	AdminConfiguration     AdminConfiguration     `mapstructure:"AdminConfiguration"`
	SecurityConfiguration  SecurityConfiguration  `mapstructure:"SecurityConfiguration"`
	RateLimitConfiguration RateLimitConfiguration `mapstructure:"RateLimitConfiguration"`
	BotConfiguration       BotConfiguration       `mapstructure:"BotConfiguration"`
}

type ServiceConfiguration struct {
	Port           string   `mapstructure:"Port"`
	Debug          bool     `mapstructure:"Debug"`
	ExportURL      string   `mapstructure:"ExportURL"`
	LogLevel       string   `mapstructure:"LogLevel"`
	TrustedProxies []string `mapstructure:"TrustedProxies"`
}

type PostgresConfiguration struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
	DBName   string `mapstructure:"DBName"`
	SSLMode  bool   `mapstructure:"SSLMode"`
	TimeZone string `mapstructure:"TimeZone"`
}

func (p PostgresConfiguration) DSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		p.Host, p.User, p.Password, p.DBName, p.Port)
	if p.TimeZone != "" {
		dsn += " TimeZone=" + p.TimeZone
	}
	if !p.SSLMode {
		dsn += " sslmode=disable"
	}
	return dsn
}

type RedisConfiguration struct {
	Addr     string `mapstructure:"Addr"`
	Password string `mapstructure:"Password"`
	Db       int    `mapstructure:"Db"`
}

type AdminConfiguration struct {
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
}

type SecurityConfiguration struct {
	TokenSecret       string `mapstructure:"TokenSecret"`
	TSWindowSeconds   int    `mapstructure:"TSWindowSeconds"`
	NonceTTLSeconds   int    `mapstructure:"NonceTTLSeconds"`
	DedupSeconds      int    `mapstructure:"DedupSeconds"`
	RSAPrivateKeyPath string `mapstructure:"RSAPrivateKeyPemPath"`
	RSAPublicKeyPath  string `mapstructure:"RSAPublicKeyPemPath"`
	KID               string `mapstructure:"KID"`
}

type RateLimitConfiguration struct {
	PerIPPerMinute        int `mapstructure:"PerIPPerMinute"`
	PerIPUAPerMinute      int `mapstructure:"PerIPUAPerMinute"`
	PerTrackerIPPerMinute int `mapstructure:"PerTrackerIPPerMinute"`
}

type BotConfiguration struct {
	MarkThreshold  int    `mapstructure:"MarkThreshold"`
	BlockThreshold int    `mapstructure:"BlockThreshold"`
	BlockMode      string `mapstructure:"BlockMode"`
}

func InitConfiguration(configName string, configPaths []string, config *Config) error {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("toml")

	for _, path := range configPaths {
		trimmed := strings.TrimSpace(path)
		if trimmed != "" {
			v.AddConfigPath(trimmed)
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err := v.Unmarshal(config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	return nil
}
