package models

import (
  "strings"
)

type Propert struct {
  Table     string
  Field     string
  Type      string
  Size      int
  Decimals  int
  Unique    string
  Required  string
  ColNull   int
}

func (Propert) TableName() string {
  return "propert"
}

func (p *Propert) GetField() string {
  return strings.ToLower(strings.TrimSpace(p.Field))
}
