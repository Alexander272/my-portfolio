package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		Http        HttpConfig
		// FileStorage FileStorageConfig
		// CacheTTL    time.Duration `mapstructure:"ttl"`
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	// FileStorageConfig struct {
	// 	Endpoint  string
	// 	Basket    string
	// 	AccessKey string
	// 	SecretKey string
	// }

	HttpConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configDir string) (*Config, error) {
	if err := parseConfigFile(configDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg *Config
	if err := unmarhal(cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("dev")

	if env == "production" {
		viper.SetConfigName("prod")
	}

	return viper.MergeInConfig()
}

func unmarhal(conf *Config) error {
	// if err := viper.UnmarshalKey("cache.ttl", &conf.CacheTTL); err != nil {
	// 	return err
	// }
	if err := viper.UnmarshalKey("mongo", &conf.Mongo); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &conf.Http); err != nil {
		return err
	}
	// if err := viper.UnmarshalKey("fileStorage", &conf.FileStorage); err != nil {
	// 	return err
	// }

	return nil
}

func setFromEnv(conf *Config) error {
	if err := envconfig.Process("mongo", conf.Mongo); err != nil {
		return err
	}
	if err := envconfig.Process("http", conf.Http); err != nil {
		return err
	}
	if err := envconfig.Process("app", conf.Environment); err != nil {
		return err
	}
	// if err := envconfig.Process("storage", conf.FileStorage); err != nil {
	// 	return err
	// }

	return nil
}
