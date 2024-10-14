package configs

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env 	   string `yaml:"env" env-default:"local"`
	HttpServer 		  `yaml:"httpServer"`
	Postgres	      `yaml:"postgres"`
	Kafka			  `yaml:"kafka"`
}

type Postgres struct {
	Host 	 string `yaml:"host"`
  	Port     string `yaml:"port"`
  	Username string `yaml:"username"`
  	Password string 
  	Name  	 string	`yaml:"dbname"`
  	SSL      string `yaml:"ssl_mode"`
}

type HttpServer struct {
	Port         string        `yaml:"port" env-default:":80"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type Kafka struct {
	Brokers     string `yaml:"brokers"`
	Topic 		string `yaml:"topic"`
}


func Init() (*Config, error){
	var config Config

	err := cleanenv.ReadConfig("./configs/config.yml", &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file: %w", err)
	}

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("could not load env variables: %w", err)
	}

	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return nil,  errors.New("database password is empty") 
	}
	config.Postgres.Password = pass

	return &config, nil
}