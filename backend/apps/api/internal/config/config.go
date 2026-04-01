package config

import (
	"fmt"
	"packages/configloader"
	"packages/configloader/parsers/yaml"
	"packages/configloader/providers/env"
	"packages/configloader/providers/file"

	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Env      string
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Host       string
	Port       int
	CtxTimeout int
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Singleton instance of AppConfig
var config *AppConfig

// Load loads the configuration from environment variables and a file,
// and returns a singleton instance of AppConfig.
// Prioritizes environment variables over file values.
func Load(configPath string) *AppConfig {
	if config != nil {
		fmt.Println("Config already loaded, returning existing instance.")
		return config
	}

	cfLoader := configloader.New()

	// Load environment variables with the specified prefix.
	cfLoader.Load(
		env.Provider(env.Opt{
			Prefix: "EASYSTUDY_",
		}),
		nil,
	)

	// Load configuration from the specified file path.
	cfLoader.Load(
		file.Provider(configPath),
		yaml.Parser(),
	)

	config = buildConfig(cfLoader)
	return config
}

// Get returns the loaded AppConfig instance.
// It panics if Load has not been called yet.
func Get() *AppConfig {
	if config == nil {
		panic("Config not loaded. Please call Load() first.")
	}
	return config
}

// buildConfig constructs an AppConfig instance from the provided ConfigLoader.
func buildConfig(cfLoader *configloader.ConfigLoader) *AppConfig {
	return &AppConfig{
		Env: pickString(cfLoader, "local", "EASYSTUDY_ENV", "env"),
		Server: ServerConfig{
			Host:       pickString(cfLoader, "localhost", "EASYSTUDY_SERVER_HOST", "server.host"),
			Port:       pickInt(cfLoader, 8080, "EASYSTUDY_SERVER_PORT", "server.port"),
			CtxTimeout: pickInt(cfLoader, 30, "EASYSTUDY_SERVER_CTX_TIMEOUT", "server.ctx_timeout"),
		},
		Postgres: PostgresConfig{
			Host:     pickString(cfLoader, "localhost", "EASYSTUDY_POSTGRES_HOST", "postgres.host"),
			Port:     pickInt(cfLoader, 5432, "EASYSTUDY_POSTGRES_PORT", "postgres.port"),
			Database: pickString(cfLoader, "EasyStudy", "EASYSTUDY_POSTGRES_DATABASE", "postgres.database"),
			Username: pickString(cfLoader, "postgres", "EASYSTUDY_POSTGRES_USERNAME", "postgres.username"),
			Password: pickString(cfLoader, "postgres", "EASYSTUDY_POSTGRES_PASSWORD", "postgres.password"),
		},
		Redis: RedisConfig{
			Host:     pickString(cfLoader, "localhost", "EASYSTUDY_REDIS_HOST", "redis.host"),
			Port:     pickInt(cfLoader, 6379, "EASYSTUDY_REDIS_PORT", "redis.port"),
			Password: pickString(cfLoader, "", "EASYSTUDY_REDIS_PASSWORD", "redis.password"),
			DB:       pickInt(cfLoader, 0, "EASYSTUDY_REDIS_DB", "redis.db"),
		},
	}
}

// pickString checks the provided keys in order and returns the first non-empty string value found.
// If no value is found, it returns the provided fallback.
func pickString(cfLoader *configloader.ConfigLoader, fallback string, keys ...string) string {
	for _, key := range keys {
		if v := cfLoader.String(key); v != "" {
			return v
		}
	}
	return fallback
}

// pickInt checks the provided keys in order and returns the first non-zero integer value found.
// If no value is found, it returns the provided fallback.
func pickInt(cfLoader *configloader.ConfigLoader, fallback int, keys ...string) int {
	for _, key := range keys {
		if v := cfLoader.Int(key); v != 0 {
			return v
		}
	}
	return fallback
}
