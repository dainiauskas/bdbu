package app

import (
  "fmt"
  "time"

  "bdbu/models"
)

type duration struct {
  Start time.Time
}

func (d *duration) Completed(tpl string) {
  fmt.Printf(tpl, time.Now().Sub(d.Start))
}

func NewDuration() *duration {
  return &duration{
    Start: time.Now(),
  }
}

func Copy(tableName string) {
  d := NewDuration()
  defer d.Completed("Completed in %v\n")

  db := models.Connect(Config.Source, Config.Destination)
  defer db.Close()

  db.Migrate(tableName)
}
