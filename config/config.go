package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server    ServerConfig
	Postgres  PostgresConfig
	Generator GeneratorConfig
}

type ServerConfig struct {
	AppVersion    string `json:"appVersion"`
	Host          string `json:"host" validate:"required"`
	Port          string `json:"port" validate:"required"`
	Timeout       time.Duration
	WithGenerator bool `json:"withGenerator"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"-"`
	DBName   string `json:"DBName"`
	SSLMode  string `json:"sslMode"`
	PgDriver string `json:"pgDriver"`
}

type GeneratorConfig struct {
	CntRow      int64  `json:"cntRow"`
	CntFeature  int64  `json:"cntFeature"`
	CntTag      int64  `json:"cntTag"`
	MaxTagInRow int64  `json:"maxTagInRow"`
	NameFile    string `json:"nameFile"`
}

func LoadConfig() (*viper.Viper, error) {

	viperInstance := viper.New()

	viperInstance.AddConfigPath("./config")
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yml")

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {

	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
		return nil, err
	}
	return &c, nil
}
