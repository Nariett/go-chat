package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Protocol   string
	DBPort     string
	LPort      string
	DBHost     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		Protocol:   os.Getenv("PROTOCOL"),
		DBPort:     os.Getenv("DB_PORT"),
		LPort:      os.Getenv("LPORT"),
		DBHost:     os.Getenv("DB_HOST"),
	}
}

func (c *Config) BuildConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func (c *Config) GetProtocolAndPort() (protocol, port string) {
	return c.Protocol, fmt.Sprintf(":%s", c.LPort)
}
