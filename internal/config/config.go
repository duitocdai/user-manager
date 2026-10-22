package config

import (
	"fmt"
	"os"

	"github.com/duitocdai/internal/utils"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
type Config struct {
	ServerAddress     string
	DB                DatabaseConfig
	MailProviderType  string
	MailProvideConfig map[string]any
}

func NewConfig() *Config {
	mailProviderConfig := make(map[string]any)
	mailProviderType := utils.GetEnv("MAIL_PROVIDER_TYPE", "mailtrap")
	if mailProviderType == "mailtrap" {
		mailTrapConfig := map[string]any{
			"mail_sender":      utils.GetEnv("MAILTRAP_MAIL_SENDER", "admin@gmail.com"),
			"name_sender":      utils.GetEnv("MAILTRAP_NAME_SENDER", "Support Team"),
			"mail_trap_url":    utils.GetEnv("MAILTRAP_URL", "https://sandbox.api.mailtrap.io/api/send/4966932"),
			"mailtrap_api_key": utils.GetEnv("MAILTRAP_API_KEY", "43ff3d018adc745f8bca52ec46a67428"),
		}
		mailProviderConfig["mailtrap"] = mailTrapConfig
	}
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),

		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "myapp"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
		MailProviderType:  mailProviderType,
		MailProvideConfig: mailProviderConfig,
	}
}
func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.SSLMode)
}
