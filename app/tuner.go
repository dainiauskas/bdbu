package app

import (
  "bdbu/tuner"
  "bdbu/models"
)

func Tuner() {
  tuner.Start(models.Open(Config.Benchmark, "dst"))
}
