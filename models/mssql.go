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

// Close connection of MsSql database
func (ms *MsSql) Close() {
  ms.DB.Close()
  fmt.Println("MsSql disconnected")
}

// GetTables assign table list to MsSql Tables slice
func (ms *MsSql) GetTables() error {
  return DBs.Preload("Propert").Preload("Keys").Order("table").Find(&ms.Tables).Error
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
