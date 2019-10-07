package app

import (
  "butent/api/config"
  "github.com/spf13/viper"
)

type Configuration interface {
  Get() (*config.Database)
}

type MainConfig struct {
  config.App
  Source      *config.Database
  Destination *config.Database
  Benchmark   *config.Database
}

func GetConfig() error {
  if err := viper.Unmarshal(&Config); err != nil {
    return err
  }

  return nil
}

var Config  *MainConfig
