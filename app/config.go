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
  MsSql       *config.Database
  MySql       *config.Database
  Source      *config.Database
  Destination *config.Database
}

func (c *MainConfig) MsConfig() *config.Database {
  return c.MsSql
}

func (c *MainConfig) MyConfig() *config.Database {
  return c.MySql
}

func GetConfig() error {
  if err := viper.Unmarshal(&Config); err != nil {
    return err
  }

  return nil
}

var Config  *MainConfig
