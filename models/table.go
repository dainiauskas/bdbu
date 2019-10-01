package models

import (
  "strings"
  "fmt"
)

type Table struct {
  Name      string    `gorm:"column:table;unique"`
  RecordId  string
  RelKey    string
  Propert   []Propert `gorm:"foreignkey:table;association_foreignkey:Name"`
  Keys      []Key     `gorm:"foreignkey:table;association_foreignkey:Name"`
}

// TableName return name of Table structure
func (Table) TableName() string {
  return "tables"
}

// GetName return name of Table structure in lowercase
func (t *Table) GetName() string {
  return strings.ToLower(strings.TrimSpace(t.Name))
}

func (t *Table) ColumnsList() string {
  var list []string

  for _, row := range t.Propert {
    list = append(list, "[" + strings.ToLower(strings.TrimSpace(row.Field)) + "]")
  }

  list = append(list, "[inp_time]")
  list = append(list, "[inp_user]")
  list = append(list, "[mod_time]")
  list = append(list, "[mod_user]")
  list = append(list, "[lock_time]")
  list = append(list, "[lock_user]")
  list = append(list, "[print_time]")
  list = append(list, "[print_user]")

  return strings.Join(list, ", ")
}

// Columns return columns in string for create table
func (t *Table) Columns() string {
  var list []string

  for _, row := range t.Propert {
    list = append(list, "\t" + row.MyCreateField(t.RecordId, t.RelKey))
  }

  list = append(list, "`inp_time` datetime DEFAULT NULL")
  list = append(list, "`inp_user` int(11) DEFAULT NULL")
  list = append(list, "`mod_time` datetime DEFAULT NULL")
  list = append(list, "`mod_user` int(11) DEFAULT NULL")
  list = append(list, "`lock_time` datetime DEFAULT NULL")
  list = append(list, "`lock_user` int(11) DEFAULT NULL")
  list = append(list, "`print_time` datetime DEFAULT NULL")
  list = append(list, "`print_user` int(11) DEFAULT NULL")

  return strings.Join(list, ",\n")
}

// CreateSql generate SQL string for create table
func (t *Table) CreateSql() string {
  return fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=cp1257 COLLATE=cp1257_lithuanian_ci;",
    t.GetName(), t.Columns())
}

func (t *Table) AddIndexes() error {
  for _, key := range t.Keys {
    // Checking for primary key, him created with table
    if key.Field != t.RecordId {
      if err := key.AddIndex(t.GetName()); err != nil {
        return err
      }
    }

  }

  return nil
}
