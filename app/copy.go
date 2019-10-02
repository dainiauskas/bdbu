package app

import (
  "fmt"
  "time"

  "bdbu/models"
)

type duration struct {
  Start time.Time
}

func (d *duration) Completed() {
  fmt.Printf("Completed in %v\n", time.Now().Sub(d.Start))
}

func NewDuration() *duration {
  return &duration{
    Start: time.Now(),
  }
}

func Copy() {
  d := NewDuration()
  defer d.Completed()

  db := models.Connect(Config.Source, Config.Destination)
  defer db.Close()

  db.Migrate()
}
