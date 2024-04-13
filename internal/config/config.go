package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env              string         `yaml:"env"`
	Server           ServerConfig   `yaml:"http_server"`
	Database         DatabaseConfig `yaml:"database"`
	Jwt              JwtConfig      `yaml:"jwt"`
	DefaultAdminPass string         `yaml:"default_admin_pass"`
}

type ServerConfig struct {
	Addr        string        `yaml:"addr"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}
type JwtConfig struct {
	Secret string `yaml:"secret"`
}
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH not found")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config")
	}
	return &cfg
}
