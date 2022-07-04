package config

import (
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	APIHost         string        `env:"APIHOST" env-default:"0.0.0.0:3000"`
	ReadTimeout     time.Duration `env:"READTIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"WRITETIMEOUT" env-default:"10s"`
	IdleTimeout     time.Duration `env:"IDLETIMEOUT" env-default:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWNTIMEOUT" env-default:"20s"`
	DB              struct {
		DBUser         string `env:"DBUSER" env-default:"postgres"`
		DBPassword     string `env:"DBPASSWORD" env-default:"postgres"`
		DBHost         string `env:"DBHOST" env-default:"localhost"`
		DBPort         string `env:"DBPORT" env-default:"5432"`
		DBName         string `env:"DBNAME" env-default:"postgres"`
		DBMaxIdleConns int    `env:"DBMAXIDLECONNS" env-default:"0"`
		DBMaxOpenConns int    `env:"DBMAXOPENCONNS" env-default:"0"`
		DBDisableTLS   bool   `env:"DBDISABLETLS" env-default:"true"`
	}
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func (cfg *Config) GetDBConnString() string {
	builder := strings.Builder{}
	builder.WriteString("postgres://")
	builder.WriteString(cfg.DB.DBUser)
	builder.WriteString(":")
	builder.WriteString(cfg.DB.DBPassword)
	builder.WriteString("@")
	builder.WriteString(cfg.DB.DBHost)
	builder.WriteString(":")
	builder.WriteString(cfg.DB.DBPort)
	builder.WriteString("/")
	builder.WriteString(cfg.DB.DBName)
	return builder.String()
}
