package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type MsSql struct {
	DB *gorm.DB
}

var DBs *gorm.DB

func init() {
	Register("mssql-src", &MsSql{})
	Register("mssql-dst", &MsSql{})
}

// Set parameters before migration
func (ms *MsSql) SetParams() {}

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
		return fmt.Sprintf("INT %d", p.Size)
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

// func (ms *MsSql) TableList() []Table {
//   return ms.Tables
// }

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

func (ms *MsSql) Quote(s string) string {
	return "[" + s + "]"
}

func (ms *MsSql) InsertSql(table string, columns []string, args []string) string {
	return fmt.Sprintf("INSERT INTO [%s] (%s) VALUES (%s)", table,
		strings.Join(columns, ","), strings.Join(args, ","))
}

// TableNotExists check if table exist in database or not
func (ms *MsSql) TableNotExists(table string) bool {
	q := fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", table)
	if err := ms.DB.Exec(q).Error; err != nil {
		return true
	}

	return false
}
