package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
	Groups   GroupsConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Driver string // "sqlite" or "postgres"
	DSN    string
}

type APIConfig struct {
	Key string
}

type GroupsConfig struct {
	AllowedIDs []int64 // 允许监控的群ID列表
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "file:seekyourguild.db?cache=shared&mode=rwc"),
		},
		API: APIConfig{
			Key: getEnv("API_KEY", "sk-123456"),
		},
		Groups: GroupsConfig{
			AllowedIDs: parseGroupIDs(getEnv("ALLOWED_GROUPS", "")),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseGroupIDs 解析群ID列表，格式: "123456,789012,345678"
func parseGroupIDs(s string) []int64 {
	if s == "" {
		return []int64{}
	}
	
	parts := strings.Split(s, ",")
	ids := make([]int64, 0, len(parts))
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if id, err := strconv.ParseInt(part, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	
	return ids
}
