package models

import (
  _ "github.com/jinzhu/gorm/dialects/mssql"
  "github.com/jinzhu/gorm"
  "github.com/cheggaaa/pb/v3"
  "database/sql"
  "fmt"
  "strings"
)

type MsSql struct {
  DB      *gorm.DB
  Tables  []Table
}

var DBs *gorm.DB

func init() {
  ms := &MsSql{}
  Register("mssql", ms)
}

func (ms *MsSql) Open(url string) {
  db, err := gorm.Open("mssql", url)
  if err != nil {
    panic(err)
  }

  fmt.Println("Connected to MsSql database")

  db.LogMode(false)

  DBs = db

  ms.DB = db
}

func (ms *MsSql) GetDB() *gorm.DB {
  return ms.DB
}

// Close connection of MsSql database
func (ms *MsSql) Close() {
  ms.DB.Close()
  fmt.Println("MsSql disconnected")
}

func (ms *MsSql) CreateTableTpl() string {
  return "CREATE TABLE `%s` (\n%s\n);"
}

func (ms *MsSql) AddIndexTpl() string {
  return "ALTER TABLE `%s` ADD KEY IF NOT EXISTS `%s` (`%s`)"
}

func (ms *MsSql) FieldCreateSql(p *Propert, primaryKey, relKey string) string {
  return "[" + p.GetField() + "]" + " " + ms.dbfToSqlType(p, primaryKey, relKey)
}

func (ms *MsSql) dbfToSqlType(p *Propert, primaryKey, relKey string) string {
  var null string

  if p.Field == primaryKey {
    null = "NOT NULL PRIMARY KEY"
  } else if p.Field == relKey {
    null = "NOT NULL UNIQUE"
  } else if p.ColNull == 0 {
    null = "NOT NULL"
  } else {
    null = "DEFAULT NULL"
  }

  switch p.Type {
  case "I":
    return fmt.Sprintf("INT %s", p.Size, null)
  case "C":
    return fmt.Sprintf("CHAR(%d) %s", p.Size, null)
  case "N":
    return fmt.Sprintf("DECIMAL(%d,%d) %s", p.Size, p.Decimals, null)
  case "M":
    return fmt.Sprintf("TEXT")
  case "D":
    return fmt.Sprintf("DATE %s", null)
  case "T":
    return fmt.Sprintf("DATETIME %s", null)
  }

  return ""
}

// AddSysFields adding system fields to columns list for creating table
func (ms *MsSql) AddSysFields(list []string) []string {
  list = append(list, "[inp_time] datetime DEFAULT NULL")
  list = append(list, "[inp_user] int(11) DEFAULT NULL")
  list = append(list, "[mod_time] datetime DEFAULT NULL")
  list = append(list, "[mod_user] int(11) DEFAULT NULL")
  list = append(list, "[lock_time] datetime DEFAULT NULL")
  list = append(list, "[lock_user] int(11) DEFAULT NULL")
  list = append(list, "[print_time] datetime DEFAULT NULL")
  list = append(list, "[print_user] int(11) DEFAULT NULL")

  return list
}

// GetTables assign table list to MsSql Tables slice
func (ms *MsSql) GetTables() ([]Table, error) {
  var tableList []Table

  err := ms.DB.Preload("Propert").Preload("Keys").Order("table").Find(&tableList).Error

  return tableList, err
}

func (ms *MsSql) CreateTables(tables []Table) error {
  return nil
}

func (ms *MsSql) TableList() []Table {
  return ms.Tables
}

// TableRows return rows from table
func (ms *MsSql) TableRows(name string) (int, *sql.Rows, error) {
  var count int

  ms.DB.Table(name).Select("*").Count(&count)
  if count == 0 {
    return 0, nil, nil
  }

  rows, err := ms.DB.Table(name).Select("*").Rows()

  return count, rows, err
}

func (ms *MsSql) GetTableRows(table Table) {
  var count int

  name := table.GetName()

  ms.DB.Table(name).Select("*").Count(&count)
  if count == 0 {
    return
  }

  tmpl := `{{ " Inserting data:" }} [{{string . "table_name"}}] {{ bar .}} {{counters .}} {{etime .}} {{percent .}}`

  bar := pb.ProgressBarTemplate(tmpl).Start(count)
  bar.Set("table_name", fmt.Sprintf("%-10v", name))
  defer bar.Finish()

  rows, err := ms.DB.Table(name).Select("*").Rows()
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  // Get the column names from the query
  var columns, values []string
  columns, err = rows.Columns()
  if err != nil {
    panic(err)
  }

  for i, column := range columns {
    columns[i] = "`" + column + "`"
    values = append(values, "?")
  }

  query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", name,
    strings.Join(columns, ","), strings.Join(values, ","))

  tx := DBy.Begin()

  for rows.Next() {
    bar.Increment()

    r := make([]interface{}, len(columns))
    for i := range r {
      r[i] = &r[i]
    }

    if err := rows.Scan(r...); err != nil {
      panic(err)
    }

    if err := tx.Exec(query, r...).Error; err != nil {
      Mytx.Rollback()
      fmt.Println(name, err)
      panic(err)
    }
  }

  tx.Commit()

  bar.Finish()
}
