package models

import (
  "fmt"

  "github.com/jinzhu/gorm"
  "github.com/cheggaaa/pb/v3"

  "butent/api/config"
)

var (
  drivers = make(map[string]Driver)
)

type Driver interface {
  Open(url string)
  Close()
  GetTables() ([]Table, error)
  GetDB() (*gorm.DB)
  CreateTableTpl() (string)
  AddIndexTpl() (string)
  FieldCreateSql(*Propert, string, string) (string)
  AddSysFields([]string) ([]string)
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

// open connect to one database by drivers and return Driver interface
func open(config *config.Database) Driver {
  drv := drivers[config.Dialect]
  drv.Open(config.FormatDSN())
  return drv
}

// Connect open connection to databases by configuration, return DB structure
func Connect(srcConfig, dstConfig *config.Database) *DB {
  return &DB{
    Src: open(srcConfig),
    Dst: open(dstConfig),
  }
}

type DB struct {
  Src         Driver
  Dst         Driver
}

// Close databases
func (db *DB) Close() {
  db.Src.Close()
  db.Dst.Close()
}

// recreateTable drop table and create new
func (db *DB) recreateTable(table *Table, bar *pb.ProgressBar) {
  tpl := db.Dst.CreateTableTpl()
  sql := table.CreateSql(db.Dst, tpl)

  bar.Set("action", fmt.Sprintf("%-10v", "drop"))
  db.Dst.GetDB().DropTableIfExists(table.GetName())

  bar.Set("action", fmt.Sprintf("%-10v", "create"))
  if err := db.Dst.GetDB().Exec(sql).Error; err != nil {
    panic(err)
  }
}

// ceateTable creating table by Table information
func (db *DB) migrateTable(table *Table) *DB {
  name := table.GetName()
  tmpl := `[{{string . "table_name"}}] [{{string . "action"}}] {{ bar .}} {{counters .}} {{etime .}} {{percent .}}`

  // start bar based on our template
  bar := pb.ProgressBarTemplate(tmpl).Start(1)
  defer bar.Finish()

  bar.Set("table_name", fmt.Sprintf("%-10v", name))
  db.recreateTable(table, bar)

  bar.Set("action", fmt.Sprintf("%-10v", "indexing"))
  table.AddIndexes(db.Dst)

  bar.Set("action", fmt.Sprintf("%-10v", "   done"))
  bar.Increment()

  return db
}

// GetTableList select list from table [tables] and preload data from tables
// [propert] and [keys]
func (db *DB) GetTableList() []Table {
  var tableList []Table

  err := db.Src.GetDB().Preload("Propert").Preload("Keys").
            Order("table").Find(&tableList).Error

  if err != nil {
    panic(err)
  }

  return tableList
}

// CreateTables create tables by []Table list on destination
func (db *DB) Migrate() {
  tables := db.GetTableList()
  for _, table := range tables {
    db.migrateTable(&table)
  }
}
