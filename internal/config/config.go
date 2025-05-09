// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	TelegramToken string
	AdminIDs      []int64
	MongoUri      string
	DBName        string
}

// New creates a new Config instance from environment variables
func New() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	adminIDs := []int64{}
	if adminIDsStr := os.Getenv("ADMIN_IDS"); adminIDsStr != "" {
		for _, idStr := range strings.Split(adminIDsStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid admin ID format: %v", err)
			}
			adminIDs = append(adminIDs, id)
		}
	}

	mongoUri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DATABASE_NAME")

	return &Config{
		TelegramToken: token,
		AdminIDs:      adminIDs,
		MongoUri:      mongoUri,
		DBName:        dbName,
	}, nil
}

// IsAdmin checks if the given user ID is an admin
func (c *Config) IsAdmin(userID int64) bool {
	for _, adminID := range c.AdminIDs {
		if adminID == userID {
			return true
		}
	}
	return false
}

var AppConfig *Config

func SetConfig(cfg *Config) {
	AppConfig = cfg
}
