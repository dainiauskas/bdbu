package models

import (
  "butent/api/config"
)

var (
  drivers = make(map[string]Driver)
)

type Driver interface {
  Open(url string)
  Close()
  GetTables() (error)
  CreateTables(tables []Table) (error)
  TableList() ([]Table)
  GetTableRows(table Table)
}

func Register(name string, driver Driver) {
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Connect(config *config.Database) Driver {
  drv := drivers[config.Dialect]

  drv.Open(config.FormatDSN())

  return drv
}
