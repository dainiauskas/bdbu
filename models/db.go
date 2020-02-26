package models

import (
	"fmt"

	"github.com/cheggaaa/pb/v3"
	"github.com/jinzhu/gorm"

	"bitbucket.org/butenta/pkg-config"
)

var (
	drivers = make(map[string]Driver)
)

// Driver interface to use databases
type Driver interface {
	Open(url string)
	Close()
	GetTables() ([]Table, error)
	GetDB() *gorm.DB
	CreateTableTpl() string
	AddIndexTpl() string
	FieldCreateSql(*Propert, string, string) string
	AddSysFields([]string) []string
	SetParams()
	Quote(string) string
	InsertSql(string, []string, []string) string
	TableNotExists(string) bool
}

// Register drivers
func Register(name string, driver Driver) {
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Open connect to one database by drivers and return Driver interface
func Open(config *config.Database, where string) Driver {
	drv := drivers[config.Dialect+"-"+where]
	drv.Open(config.FormatDSN())

	return drv
}

// Connect open connection to databases by configuration, return DB structure
func Connect(srcConfig, dstConfig *config.Database) *DB {
	return &DB{
		Src: Open(srcConfig, "src"),
		Dst: Open(dstConfig, "dst"),
	}
}

type DB struct {
	Src Driver
	Dst Driver
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

	rows, columns := table.GetRows(db.Src)
	defer rows.Close()

	// Get the column names from the query
	var values []string
	for i, column := range columns {
		columns[i] = db.Dst.Quote(column)
		values = append(values, "?")
	}

	query := db.Dst.InsertSql(name, columns, values)

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

		if err := tx.Exec(query, r...).Error; err != nil {
			tx.Rollback()
			fmt.Println(name, err)
			panic(err)
		}
	}

	tx.Commit()

	return db
}

// GetTableList select list from table [tables] and preload data from tables
// [propert] and [keys]
func (db *DB) GetTableList(tableName string) []Table {
	var tableList []Table

	qr := db.Src.GetDB().Preload("Propert").Preload("Keys")

	if tableName != "" {
		qr = qr.Where(map[string]interface{}{"table": tableName})
	}

	err := qr.Order("table").Find(&tableList).Error

	if err != nil {
		panic(err)
	}

	return tableList
}

// Migrate function to migrate tables from source to destination
func (db *DB) Migrate(tableName string, dropTables bool) {
	db.Src.SetParams()
	db.Dst.SetParams()

	tables := db.GetTableList(tableName)

	for _, table := range tables {
		ok := true
		name := table.GetName()

		if !dropTables {
			ok = db.Src.TableNotExists(name)
		}

		if ok {
			db.migrateTable(&table)
		} else {
			fmt.Printf("Skipping table: %s\n", name)
		}
	}
}
