package models

import (
  "fmt"
  "strings"
)

type Key struct {
  Table   string
  Field   string
  TagType *string
}

func (Key) TableName() string {
  return "keys"
}

// Name return field name in lowercase
func (k *Key) Name() string {
  return strings.TrimSpace(strings.ToLower(k.Field))
}

// AddIndex create index for table
func (k *Key) AddIndex(drv Driver, table string) error {
  return drv.GetDB().Exec(fmt.Sprintf(drv.AddIndexTpl(),
            table, k.Name(), k.Name())).Error
}
