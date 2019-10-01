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

func (k *Key) Name() string {
  return strings.TrimSpace(strings.ToLower(k.Field))
}

func (k *Key) AddIndex(table string) error {
  name := k.Name()
  query := fmt.Sprintf("ALTER TABLE `%s` ADD KEY IF NOT EXISTS `%s` (`%s`)", table, name, name)

  return DBy.Exec(query).Error
}
