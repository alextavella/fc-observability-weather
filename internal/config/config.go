package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PORT       int    `mapstructure:"PORT"`
	APP_NAME   string `mapstructure:"APP_NAME"`
	APP_A_HOST string `mapstructure:"APP_A_HOST"`
	APP_B_HOST string `mapstructure:"APP_B_HOST"`
	OTEL_HOST  string `mapstructure:"OTEL_HOST"`
}

func NewConfig(path, name string) (*Config, error) {
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	fmt.Println("Configuração lida com sucesso:")
	fmt.Printf("PORT: %d\n", cfg.PORT)
	fmt.Printf("APP_NAME: %s\n", cfg.APP_NAME)
	fmt.Printf("APP_A_HOST: %s\n", cfg.APP_A_HOST)
	fmt.Printf("APP_B_HOST: %s\n", cfg.APP_B_HOST)
	fmt.Printf("OTEL_HOST: %s\n", cfg.OTEL_HOST)

	return &cfg, nil
}
