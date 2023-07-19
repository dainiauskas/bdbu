package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MySql struct {
	DB *gorm.DB
}

func init() {
	Register("mysql-src", &MySql{})
	Register("mysql-dst", &MySql{})
}

// SetParams set parameters before migration
func (my *MySql) SetParams() {
	my.DB.Exec("SET @@sql_mode='';")
}

// Open connect to database and and assign to DB
func (my *MySql) Open(url string) {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MySql database")

	db.LogMode(false)

	my.DB = db
}

// GetDB return gorm DB
func (my *MySql) GetDB() *gorm.DB {
	return my.DB
}

// Close connection of MySql database
func (my *MySql) Close() {
	my.DB.Close()
	fmt.Println("MySql disconnected")
}

// CreateTableTpl return create table string
func (my *MySql) CreateTableTpl() string {
	return "CREATE TABLE IF NOT EXISTS `%s` (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=cp1257 COLLATE=cp1257_lithuanian_ci;"
}

// AddIndexTpl return alter table string
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
		return fmt.Sprintf("LONGTEXT")
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

func (my *MySql) CreateIndexes(t *Table) error {
	for _, key := range t.Keys {
		if key.Field != t.RecordID {
			q := fmt.Sprintf("ALTER TABLE `%s` ADD KEY IF NOT EXISTS `%s` (`%s`);",
				t.GetName(), key.Field, key.Field)

			if err := my.DB.Exec(q).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (my *MySql) Quote(s string) string {
	return "`" + s + "`"
}

func (my *MySql) InsertSql(table string, columns []string, args []string) string {
	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table,
		strings.Join(columns, ","), strings.Join(args, ","))
}

// TableNotExists check if table exist in database or not
func (my *MySql) TableNotExists(table string) bool {
	q := fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", table)
	if err := my.DB.Exec(q).Error; err != nil {
		return true
	}

	return false
}

func (my *MySql) Truncate(table string) {
	my.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", my.Quote(table)))
}
