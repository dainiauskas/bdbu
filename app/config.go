package app

import (
  "butent/api/config"
  "github.com/spf13/viper"
  "github.com/gookit/color"
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

func (mc *MainConfig) IsBenchmark() bool {
  if (Config.Benchmark == nil) {
    color.Red.Println("Please configure Benchmark section")
    return false
  }

  return true
}

func (mc *MainConfig) IsConfigured() bool {
  if (Config.Source == nil || Config.Destination == nil) {
    color.Red.Println("Please configure source and destination sections")
    return false
  }

  return true
}
