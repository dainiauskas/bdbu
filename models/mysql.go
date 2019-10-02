package models

import (
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/jinzhu/gorm"
  // "github.com/cheggaaa/pb/v3"
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

func (my *MySql) GetDB() *gorm.DB {
  return my.DB
}

// Close connection of MySql database
func (my *MySql) Close() {
  Mytx.Commit()
  my.DB.Close()
  fmt.Println("MySql disconnected")
}

func (my *MySql) CreateTableTpl() string {
  return "CREATE TABLE IF NOT EXISTS `%s` (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=cp1257 COLLATE=cp1257_lithuanian_ci;"
}

func (my *MySql) AddIndexTpl() string {
  return "ALTER TABLE `%s` ADD KEY IF NOT EXISTS `%s` (`%s`)"
}

func (my *MySql) FieldCreateSql(p *Propert, primaryKey, relKey string) string {
  return "`" + p.GetField() + "`" + " " + my.dbfToSqlType(p, primaryKey, relKey)
}

func (my *MySql) dbfToSqlType(p *Propert, primaryKey, relKey string) string {
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
    return fmt.Sprintf("INT(%d) %s", p.Size, null)
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

func (my *MySql) AddSysFields(list []string) []string {
  list = append(list, "`inp_time` datetime DEFAULT NULL")
  list = append(list, "`inp_user` int(11) DEFAULT NULL")
  list = append(list, "`mod_time` datetime DEFAULT NULL")
  list = append(list, "`mod_user` int(11) DEFAULT NULL")
  list = append(list, "`lock_time` datetime DEFAULT NULL")
  list = append(list, "`lock_user` int(11) DEFAULT NULL")
  list = append(list, "`print_time` datetime DEFAULT NULL")
  list = append(list, "`print_user` int(11) DEFAULT NULL")

  return list
}


// GetTables assign table list to MsSql Tables slice
func (my *MySql) GetTables() ([]Table, error) {
  var tableList []Table

  err := my.DB.Preload("Propert").Preload("Keys").Order("table").Find(&tableList).Error

  return tableList, err
}

// func (my *MySql) CreateTables(tables []Table) error {
//   tmpl := `{{ "Creating tables:" }} [{{string . "table_name"}}] {{ bar .}} {{counters .}} {{etime .}} {{percent .}}`
//
//   // start bar based on our template
//   bar := pb.ProgressBarTemplate(tmpl).Start(len(tables))
//   defer bar.Finish()
//
//   for _, table := range tables {
//     bar.Increment()
//
//     name := table.GetName()
//
//     bar.Set("table_name", fmt.Sprintf("%-10v", name))
//
//     my.DB.DropTableIfExists(name)
//
//     if err := my.CreateTable(table); err != nil {
//       return err
//     }
//
//     my.CreateIndexes(&table)
//     // if err := table.AddIndexes(); err != nil {
//     //   // fmt.Println(name, err)
//     //   // return err
//     // }
//   }
//
//   return nil
// }

func (my *MySql) CreateIndexes(t *Table) error {
  for _, key := range t.Keys {
    if key.Field != t.RecordId {
      q := fmt.Sprintf("ALTER TABLE `%s` ADD KEY IF NOT EXISTS `%s` (`%s`);",
              t.GetName(), key.Field, key.Field)

      if err := my.DB.Exec(q).Error; err != nil {
        return err
      }
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
// func (my *MySql) CreateTable(table Table) error {
//   return my.Exec(table.CreateSql())
// }

func (my *MySql) Exec(sql string) error {
  return my.DB.Exec(sql).Error
}
