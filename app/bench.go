package app

import (
	"fmt"

	"bdbu/models"
)

const (
	tplRecords = "Test records   : %d\n"
	tplCreate  = "Creating table : %v\n"
	tplInsert  = "Inserting      : %v\n"
	tplSelect  = "Selecting      : %v\n"
	tplDelete  = "Deleting       : %v\n"
	tplDrop    = "Droping        : %v\n"

	sqlGetInfo = `engine, row_format,
    round(((data_length + index_length) / 1024 / 1024), 2) As size`
)

// InfoSchema structure for reading information schema tables
type InfoSchema struct {
	Engine      string  `gorm:"column:ENGINE"`
	RowFormat   string  `gorm:"column:ROW_FORMAT"`
	DataLength  float64 `gorm:"column:DATA_LENGTH"`
	IndexLength float64 `gorm:"column:INDEX_LENGTH"`
}

// TableName return information schema tables name
func (InfoSchema) TableName() string {
	return "information_schema.TABLES"
}

// TableBench used for create table to benchmark tests
type TableBench struct {
	ID uint `gorm:"primary_key"`
	// Str   string    `gorm:"type:char(10);default:''"`
}

// TableName return name of table
func (TableBench) TableName() string {
	return "temp_bench_table"
}

func Benchmark(records int, engine, rowFormat string) {
	if !Config.IsBenchmark() {
		return
	}

	drv := models.Open(Config.Benchmark, "src")
	defer drv.Close()

	db := drv.GetDB()

	fmt.Printf(tplRecords, records)

	d := NewDuration()
	err := db.DropTableIfExists(&TableBench{}).
		Set("gorm:table_options",
			fmt.Sprintf("ENGINE=%s ROW_FORMAT=%s", engine, rowFormat)).
		AutoMigrate(&TableBench{}).Error

	if err != nil {
		panic(err)
	}
	d.Completed(tplCreate)

	d = NewDuration()
	tx := db.Begin()
	for i := 0; i < records; i++ {
		tx.Create(&TableBench{})
	}
	tx.Commit()
	d.Completed(tplInsert)

	list := []TableBench{}

	d = NewDuration()
	db.Find(&list)
	d.Completed(tplSelect)

	d = NewDuration()
	db.Delete(&list)
	d.Completed(tplDelete)

	info := &InfoSchema{}

	db.Where("table_schema = ?", Config.Benchmark.Name).
		Where("table_name = ?", list[0].TableName()).First(&info)

	d = NewDuration()
	db.DropTableIfExists(&TableBench{})
	d.Completed(tplDrop)

	fmt.Printf("Table size (MB): %v, Engine: %v, Row format: %v\n",
		((info.DataLength + info.IndexLength) / 1024 / 1024),
		info.Engine, info.RowFormat)
}
