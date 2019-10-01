package app

import (
  "fmt"
  "time"

  "bdbu/models"
)

func Copy() {
  starting := time.Now()

  src := models.Connect(Config.Source)
  defer src.Close()

  dst := models.Connect(Config.Destination)
  defer dst.Close()


  if err := src.GetTables(); err != nil {
    panic(err)
  }
  fmt.Printf("%+v\n", src.TableList())

  if err := dst.CreateTables(src.TableList()); err != nil {
    panic(err)
  }

  for _, table := range src.TableList() {
    src.GetTableRows(table)
  }

  ending := time.Now()

  fmt.Printf("Completed in %v\n", ending.Sub(starting))
}
