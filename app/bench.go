package app

import (
  "bdbu/models"
)

type TableBench struct {
  Id    uint `gorm:"primary_key"`
}

func (TableBench) TableName() string {
  return "bench_table"
}

func Benchmark() {
  drv := models.Open(Config.Benchmark, "src")
  defer drv.Close()

  db := drv.GetDB()

  db.DropTableIfExists(&TableBench{}).AutoMigrate(&TableBench{})

  d := NewDuration()
  tx := db.Begin()
  for i := 0; i < 100000; i++ {
    tx.Create(&TableBench{})
  }
  tx.Commit()
  d.Completed()

  db.DropTableIfExists(&TableBench{})
}
