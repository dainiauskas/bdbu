package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Table struct {
	Name     string `gorm:"column:table;unique"`
	RecordID string
	RelKey   string
	Mode     string    `gorm:"column:table_mode"`
	Records  uint      `gorm:"column:rec_count"`
	Propert  []Propert `gorm:"foreignkey:table;association_foreignkey:Name"`
	Keys     []Key     `gorm:"foreignkey:table;association_foreignkey:Name"`
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
		list = append(list, "["+strings.ToLower(strings.TrimSpace(row.Field))+"]")
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

// Columns return columns list for create table
func (t *Table) Columns(drv Driver) string {
	var list []string

	for _, row := range t.Propert {
		list = append(list, "\t"+drv.FieldCreateSql(&row, t.RecordID, t.RelKey))
	}

	list = drv.AddSysFields(list)

	return strings.Join(list, ",\n")
}

// CreateSql generate SQL string for create table
func (t *Table) CreateSql(drv Driver, tpl string) string {
	return fmt.Sprintf(tpl, t.GetName(), t.Columns(drv))
}

// AddIndexes create indexes for table
func (t *Table) AddIndexes(drv Driver) {
	for _, key := range t.Keys {
		// Checking for primary key, him created with table
		if key.Field != t.RecordID {
			if err := key.AddIndex(drv, t.GetName()); err != nil {
				panic(err)
			}
		}
	}
}

// CountRecords get rows numbers in table
func (t *Table) CountRecords(drv Driver) int {
	var count int
	drv.GetDB().Table(t.GetName()).Select("*").Count(&count)
	return count
}

// GetRows get rows cursor to table
func (t *Table) GetRows(drv Driver) (*sql.Rows, []string) {
	rows, err := drv.GetDB().Table(t.GetName()).Select("*").Rows()
	if err != nil {
		panic(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	return rows, columns
}

// CanByEmpty check if table have been with data or not
func (t *Table) CanByEmpty() bool {
	if strings.Index(t.Mode, "u") >= 0 {
		return false
	}

	return t.Records == 0
}
