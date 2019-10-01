package models

import (
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/jinzhu/gorm"
  "github.com/cheggaaa/pb/v3"
  "database/sql"
  "fmt"
)

type MySql struct {
  DB      *gorm.DB
  Tables  []Table
}

var DBy   *gorm.DB
var Mytx  *gorm.DB

func init() {
  my := &MySql{}
  Register("mysql", my)
}

func (my *MySql) Open(url string) {
  db, err := gorm.Open("mysql", url)
  if err != nil {
    panic(err)
  }

  fmt.Println("Connected to MySql database")

  db.LogMode(false)

  DBs = db
  Mytx = db.Begin()

  my.DB = db
}

// Close connection of MySql database
func (my *MySql) Close() {
  Mytx.Commit()
  my.DB.Close()
  fmt.Println("MySql disconnected")
}

// GetTables assign table list to MsSql Tables slice
func (my *MySql) GetTables() error {
  return DBs.Preload("Propert").Preload("Keys").Order("table").Find(&my.Tables).Error
}

func (my *MySql) CreateTables(tables []Table) error {
  tmpl := `{{ "Creating tables:" }} [{{string . "table_name"}}] {{ bar .}} {{counters .}} {{etime .}} {{percent .}}`

  // start bar based on our template
  bar := pb.ProgressBarTemplate(tmpl).Start(len(tables))
  defer bar.Finish()

  for _, table := range tables {
    bar.Increment()

    name := table.GetName()

    bar.Set("table_name", fmt.Sprintf("%-10v", name))

    DBy.DropTableIfExists(name)

    if err := my.CreateTable(table); err != nil {
      return err
    }

    if err := table.AddIndexes(); err != nil {
      // fmt.Println(name, err)
      // return err
    }
  }

  return nil
}

func (my *MySql) TableList() []Table {
  return my.Tables
}

func (my *MySql) GetTableRows(table Table) {
}

func (my *MySql) MigrateTable(count int, rows *sql.Rows) error {
  return nil
}

// CreateTable create table in MySql database
func (my *MySql) CreateTable(table Table) error {
  return my.Exec(table.CreateSql())
}

func (my *MySql) Exec(sql string) error {
  return my.DB.Exec(sql).Error
}
