package config

import (
	"fmt"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type AppConfig struct {
	Env      string         `koanf:"env"`
	Server   ServerConfig   `koanf:"server"`
	Postgres PostgresConfig `koanf:"postgres"`
	Redis    RedisConfig    `koanf:"redis"`
}

type ServerConfig struct {
	Host       string `koanf:"host"`
	Port       int    `koanf:"port"`
	CtxTimeout int    `koanf:"ctx_timeout"`
}

type PostgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Database string `koanf:"database"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

// Singleton instance of AppConfig
var config *AppConfig

var k = koanf.New(".")

// Load loads the configuration from environment variables and a file,
// and returns a singleton instance of AppConfig.
// Prioritizes environment variables over file values.
func Load(configPath string) *AppConfig {
	if config != nil {
		fmt.Println("Config already loaded, returning existing instance.")
		return config
	}

	// Load configuration from the specified file path.
	k.Load(
		file.Provider(configPath),
		yaml.Parser(),
	)

	// Load environment variables and merge into the loaded config.
	// "EASYSTUDY" is the prefix to filter the env vars by.
	// TransformFunc is used to transform the env var names to match the config keys.
	k.Load(env.Provider(".", env.Opt{
		Prefix: "EASYSTUDY_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(
				strings.TrimPrefix(k, "EASYSTUDY_")), "_", ".")

			// If there is a space in the value, split the value into a slice by the space.
			if strings.Contains(v, ",") {
				return k, strings.Split(v, ",")
			}
			return k, v
		},
	}), nil)

	k.Unmarshal("", &config)
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
