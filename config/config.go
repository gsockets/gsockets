package config

import (
	"github.com/gsockets/gsockets"
	"github.com/spf13/viper"
)

func Load(configPath string) (Config, error) {
	vp := viper.New()
	vp.AddConfigPath(configPath)
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")

	if err := vp.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := vp.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

type Config struct {
	Server
	AppManager `mapstructure:"app_manager"`
}

type AppManager struct {
	Driver string
	Array  []gsockets.App
}

type Server struct {
	Port int
}
