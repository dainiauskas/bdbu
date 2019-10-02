package models

import (
  "fmt"
  "strings"
  // "reflect"

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
  SetParams()
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

  count := table.CountRecords(db.Src)

  // start bar based on our template
  bar := pb.ProgressBarTemplate(tmpl).Start(count)
  defer bar.Finish()
  defer bar.Set("action", fmt.Sprintf("%-10v", "   done"))

  bar.Set("table_name", fmt.Sprintf("%-10v", name))
  db.recreateTable(table, bar)

  bar.Set("action", fmt.Sprintf("%-10v", "indexing"))
  table.AddIndexes(db.Dst)

  if count == 0 {
    return db
  }

  bar.Set("action", fmt.Sprintf("%-10v", "inserting"))
  
  rows, err := db.Src.GetDB().Table(name).Select("*").Rows()
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

  tx := db.Dst.GetDB().Begin()
  for rows.Next() {
    bar.Increment()

    r := make([]interface{}, len(columns))
    for i := range r {
      r[i] = &r[i]
    }

    if err := rows.Scan(r...); err != nil {
      panic(err)
    }

    for i, v := range r {
      switch v.(type) {
      case string:
        r[i] = strings.TrimSpace(v.(string))
      }
    }

    if err := tx.Exec(query, r...).Error; err != nil {
      tx.Rollback()
      fmt.Println(name, err)
      fmt.Printf("%#v\n", columns)
      fmt.Printf("%#v\n", r)
      panic(err)
    }
  }

  bar.Set("action", fmt.Sprintf("%-10v", "   done"))

  tx.Commit()

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

  db.Dst.SetParams()

  for _, table := range tables {
    db.migrateTable(&table)
  }
}
