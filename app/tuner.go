package app

import (
  "bdbu/tuner"
  "bdbu/models"
)

func Tuner() {
  if (Config.IsBenchmark()) {
    tuner.Start(models.Open(Config.Benchmark, "dst"))
  }
}
